
function! counteria#messenger#clear() abort
    let f = {}
    function! f.default(message) abort
        echomsg a:message
    endfunction

    let s:func = { message -> f.default(message) }
endfunction

call counteria#messenger#clear()


function! counteria#messenger#set_func(func) abort
    let s:func = { message -> a:func(message) }
endfunction

function! counteria#messenger#warn(message) abort
    echohl WarningMsg
    call s:func('[counteria] ' . a:message)
    echohl None
endfunction

function! counteria#messenger#error(message) abort
    echohl ErrorMsg
    call s:func('[counteria] ' . a:message)
    echohl None
endfunction
