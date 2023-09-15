if exists("b:current_syntax")
    finish
endif

setlocal filetype=markdown

syntax match aiLine /AI: .*/
hi def link aiLine ModeMsg

syntax match userLine /USER: .*/
hi def link userLine Question

let b:current_syntax = "chat"
