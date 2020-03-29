package route

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Route :
type Route struct {
	path    string
	pattern *regexp.Regexp
}

var placeHolder = regexp.MustCompile(":([A-Za-z0-9]+)")

func newRoute(path string) Route {
	replaced := placeHolder.ReplaceAllString(path, "(?P<$1>\\d+)")
	pattern := fmt.Sprintf("^%s$", replaced)
	return Route{
		path:    path,
		pattern: regexp.MustCompile(pattern),
	}
}

// Match :
func (r Route) Match(path string) (Params, bool) {
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
	if strings.Count(r.path, "/:") != len(params) {
		return "", errors.Errorf("invalid params: %s", params)
	}

	replaces := []string{}
	for key, val := range params {
		old := ":" + key
		replaces = append(replaces, old, val)
	}
	replacer := strings.NewReplacer(replaces...)
	replaced := replacer.Replace(r.path)

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
	TasksNew = newRoute(Schema + "tasks/new")
	// TasksOne :
	TasksOne = newRoute(Schema + "tasks/:taskId")
	// TasksList :
	TasksList = newRoute(Schema + "tasks")
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

// Reads : read routes
var Reads = Routes{
	TasksNew,
	TasksOne,
	TasksList,
}

// Writes : write routes
var Writes = Routes{
	TasksNew,
}

// Match :
func (routes Routes) Match(path string) (Route, Params, error) {
	for _, r := range routes {
		params, ok := r.Match(path)
		if ok {
			return r, params, nil
		}
	}
	return Route{}, nil, errors.Errorf("not matched route: %s", path)
}
