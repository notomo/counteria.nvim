package route

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Method :
type Method string

var (
	// MethodRead :
	MethodRead = Method("read")
	// MethodWrite : create, update
	MethodWrite = Method("write")
	// MethodDelete :
	MethodDelete = Method("delete")
)

// Methods :
type Methods []Method

// Match :
func (methods Methods) Match(method Method) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

// Route :
type Route struct {
	Path    string
	pattern *regexp.Regexp
	methods Methods
}

var placeHolder = regexp.MustCompile(":([A-Za-z0-9]+)")

func newRoute(path string, methods ...Method) Route {
	replaced := placeHolder.ReplaceAllString(path, "(?P<$1>\\d+)")
	pattern := fmt.Sprintf("^%s$", replaced)
	return Route{
		Path:    path,
		pattern: regexp.MustCompile(pattern),
		methods: methods,
	}
}

// Match :
func (r Route) Match(method Method, path string) (Params, bool) {
	if !r.methods.Match(method) {
		return nil, false
	}

	match := r.pattern.FindStringSubmatch(path)
	if len(match) == 0 {
		return nil, false
	}

	params := Params{}
	for i, name := range r.pattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			params[name] = match[i]
		}
	}
	return params, true
}

// BuildPath :
func (r Route) BuildPath(params Params) (string, error) {
	if strings.Count(r.Path, "/:") != len(params) {
		return "", errors.Errorf("invalid params: %s", params)
	}

	replaces := []string{}
	for key, val := range params {
		old := ":" + key
		replaces = append(replaces, old, val)
	}
	replacer := strings.NewReplacer(replaces...)
	replaced := replacer.Replace(r.Path)

	remain := strings.Index(replaced, "/:")
	if remain != -1 {
		return "", errors.Errorf("build path: %s", replaced)
	}

	return replaced, nil
}

// Schema :
const Schema = "counteria://"

var (
	// TasksNew :
	TasksNew = newRoute(Schema+"tasks/new", MethodRead, MethodWrite)
	// TasksOne :
	TasksOne = newRoute(Schema+"tasks/:taskId", MethodRead, MethodDelete)
	// TasksList :
	TasksList = newRoute(Schema+"tasks", MethodRead)
)

// Params :
type Params map[string]string

// TaskID :
func (params Params) TaskID() int {
	id, err := strconv.Atoi(params["taskId"])
	if err != nil {
		panic(err)
	}
	return id
}

// Routes :
type Routes []Route

// All : all routes
var All = Routes{
	TasksNew,
	TasksOne,
	TasksList,
}

// Match :
func (routes Routes) Match(method Method, path string) (Route, Params, error) {
	for _, r := range routes {
		params, ok := r.Match(method, path)
		if ok {
			return r, params, nil
		}
	}
	return Route{}, nil, errors.Errorf("not matched route: %s", path)
}
