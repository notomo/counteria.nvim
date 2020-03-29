
function! counteria#job#new() abort
    let job = {
        \ 'id': 0,
    \ }

    function! job.start(executable) abort
        if self.id != 0 && self.is_running()
            return
        endif

        let id = jobstart([a:executable], {
            \ 'rpc': v:true,
            \ 'on_stderr': function('s:on_stderr')
        \ })
        if id <= 0
            throw 'failed to start job: ' . id
        endif

        let self.id = id
    endfunction

    function! job.call(sync, method, ...) abort
        if a:sync
            return call('rpcrequest', [self.id, a:method] + a:000)
        endif
        return call('rpcnotify', [self.id, a:method] + a:000)
    endfunction

    function! job.wait() abort
        if !self.is_running()
            return
        endif
        " HACk: for enqueue notification and wait it
        call rpcnotify(self.id, 'startWaiting')
        call rpcrequest(self.id, 'wait')
    endfunction

    function! job.stop() abort
        if !self.is_running()
            return
        endif
        call jobstop(self.id)
    endfunction

    function! job.is_running() abort
        return jobwait([self.id], 0)[0] == -1
    endfunction

    return job
endfunction

function! s:on_stderr(id, data, event) dict
    let msgs = join(a:data, "\n")
    if empty(msgs)
        return
    endif
    let msgs = substitute(msgs, "\t", '  ', 'g')
    for msg in split(msgs, "\n")
        echomsg '[counteria] ' . msg
    endfor
endfunction
