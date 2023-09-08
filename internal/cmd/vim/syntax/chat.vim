setlocal filetype=markdown

syntax match aiHighlight "^AI:\s" contains=ALL
" syntax match userHighlight "^USER:\s" contains=ALL

hi aiHighlight ctermbg=blue guibg=blue

" hi userHighlight ctermbg=green guibg=green

if exists("b:current_syntax")
    finish
endif
let b:current_syntax = "chat"
