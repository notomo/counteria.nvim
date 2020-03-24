
let s:helper = TestHelper()
let s:suite = s:helper.suite('plugin.counteria')
let s:assert = s:helper.assert

function! s:suite.tasks_new()
    call s:helper.sync_execute('tasks/new')
    call s:assert.tab_count(2)

    call s:helper.sync_write()
    call s:assert.match_path('counteria://tasks/\d+')
endfunction

function! s:suite.tasks_list()
    call s:helper.sync_execute('tasks')
    call s:assert.tab_count(2)
endfunction
