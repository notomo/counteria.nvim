if exists('g:loaded_counteria')
    finish
endif
let g:loaded_counteria = 1

command! -nargs=* Counteria call counteria#main(<f-args>)
