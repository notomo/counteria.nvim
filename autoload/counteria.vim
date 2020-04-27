
function! counteria#main(...) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    call s:job.start()

    let sync = v:false
    call s:job.call(sync, 'do', a:000)

    return s:job
endfunction

function! counteria#request(method, sync, path, bufnr) abort
    if !exists('s:job')
        let s:job = counteria#job#new()
    endif
    call s:job.start()

    call s:job.call(a:sync, 'exec', a:method, a:path, str2nr(a:bufnr))

    return s:job
endfunction

function! counteria#last_job() abort
    if exists('s:job')
        return s:job
    endif
    return v:null
endfunction
