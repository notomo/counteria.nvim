
let s:helper = TestHelper()
let s:suite = s:helper.suite('plugin.counteria')
let s:assert = s:helper.assert

function! s:suite.tasks()
    call s:helper.sync_execute('open', 'tasks/new')
    call s:helper.search('name')

    call s:helper.sync_write()
    call s:assert.match_path('counteria://tasks/\d+')

    call s:helper.sync_execute('open', 'tasks')
    call s:assert.match_path('counteria://tasks')

    call s:helper.sync_execute('do')
    call s:assert.match_path('counteria://tasks/\d+')
endfunction

function! s:suite.open_tasks_new()
    call s:helper.sync_read('counteria://tasks/new')
    call s:helper.search('name')
endfunction
