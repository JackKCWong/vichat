if exists("b:current_syntax")
    finish
endif


setlocal filetype=markdown

syntax match sysLine /SYSTEM: .*/
hi def link sysLine DiffAdd

syntax match aiLine /AI: .*/
hi def link aiLine DiffText

syntax match userLine /USER: .*/
hi def link userLine DiffChange

let b:current_syntax = "chat"

