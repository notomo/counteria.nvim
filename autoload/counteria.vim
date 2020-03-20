
let s:executable = expand('<sfile>:h:h') . '/bin/counteriad'

function! counteria#main(...) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif

    " for debug
    call s:job.stop()

    call s:job.start(s:executable)
    call s:job.notify(a:000)

    return s:job
endfunction
