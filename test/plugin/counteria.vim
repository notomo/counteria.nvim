
let s:helper = TestHelper()
let s:suite = s:helper.suite('plugin.counteria')
let s:assert = s:helper.assert

function! s:suite.do()
    call s:helper.sync_execute('task', 'create')

    call s:assert.tab_count(2)
endfunction
