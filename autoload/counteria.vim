
let s:executable = expand('<sfile>:h:h') . '/bin/counteriad'

function! counteria#main(...) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    " for debug
    call s:job.stop()

    call s:job.start(s:executable)

    let sync = v:false
    call s:job.call(sync, 'do', a:000)

    return s:job
endfunction

function! counteria#request(method, sync, bufnr) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    call s:job.start(s:executable)

    call s:job.call(a:sync, 'request', a:method, a:bufnr)

    return s:job
endfunction

function! counteria#last_job() abort
    if exists('s:job')
        return s:job
    endif
    return v:null
endfunction
