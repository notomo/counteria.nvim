if exists('g:loaded_counteria')
    finish
endif
let g:loaded_counteria = 1

command! -nargs=* Counteria call counteria#main(<f-args>)

augroup counteria
    autocmd!
    autocmd BufReadCmd counteria://* call counteria#request('read', v:false, '', expand('<abuf>'))
    autocmd BufWriteCmd counteria://* call counteria#request('write', v:false, '', expand('<abuf>'))
augroup END
