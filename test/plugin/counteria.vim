
let s:helper = TestHelper()
let s:suite = s:helper.suite('plugin.counteria')
let s:assert = s:helper.assert

function! s:suite.do_tasks_new()
    call s:helper.sync_execute('open', 'tasks/new')
    call s:helper.search('name')

    call s:helper.sync_write()
    call s:assert.match_path('counteria://tasks/\d+')

    call s:helper.sync_execute('open', 'tasks')
    call s:assert.match_path('counteria://tasks')

    normal! j
    call s:helper.sync_execute('do', 'done')

    call s:helper.sync_execute('do')
    call s:assert.match_path('counteria://tasks/\d+')
endfunction

function! s:suite.open_tasks_new()
    call s:helper.sync_read('counteria://tasks/new')
    call s:helper.search('name')
    call s:helper.replace_line('"name": "new_task",')
    call s:helper.sync_write()
    call s:assert.match_path('counteria://tasks/\d+')

    call s:helper.sync_execute('open', 'tasks')
    call s:assert.match_path('counteria://tasks')

    let line = line('$')
    normal! j
    call s:helper.sync_execute('do', 'delete')

    call s:assert.match_path('counteria://tasks')
    call s:assert.line_count(line - 1)
endfunction

function! s:suite.update_task()
    call s:helper.sync_read('counteria://tasks/new')
    call s:helper.search('name')
    call s:helper.replace_line('"name": "new_task",')
    call s:helper.sync_write()

    call s:helper.search('name')
    call s:helper.replace_line('"name": "updated_task",')
    call s:helper.sync_write()
    call s:assert.match_path('counteria://tasks/\d+')

    call s:helper.search('updated_task')
endfunction
