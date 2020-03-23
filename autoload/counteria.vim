
let s:executable = expand('<sfile>:h:h') . '/bin/counteriad'

function! counteria#main(...) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    " for debug
    call s:job.stop()

    call s:job.start(s:executable)

    call s:job.notify('do', a:000)

    return s:job
endfunction

function! counteria#read(bufnr) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    call s:job.start(s:executable)

    call s:job.notify('read', a:bufnr)

    return s:job
endfunction

function! counteria#write(bufnr) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    call s:job.start(s:executable)

    call s:job.notify('write', a:bufnr)

    return s:job
endfunction

function! counteria#last_job() abort
    if exists('s:job')
        return s:job
    endif
    return v:null
endfunction
