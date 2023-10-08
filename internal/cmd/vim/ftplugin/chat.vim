
function! SendToChat()
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))
    let output[0] = "AI: " . output[0]

    " Append the output of the command to the current buffer
    let l = getline('$')
    if l != '' 
        let output = [""] + output
    endif

    norm! G
    call append(line('$'), output + ["", "", "USER: "])
    norm! 3j

endfunction

function! TryToChat()
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))

    exe "vnew"
    setlocal buftype=nofile nobuflisted syntax=markdown

    call append(line('$'), output)
endfunction

function! CountTokens(ran)
    " Redirect the content of the current buffer to the external command's stdin
    if a:ran == 0
        let selection = getline(1, "$")
    else
        let selection = getline("'<", "'>")
    end 
    let output = systemlist("vichat tok", selection)

    echo "estimate: " . output[0] . " tokens"
endfunction

function! SplitText(ran)
    " Redirect the content of the current buffer to the external command's stdin
    if a:ran == 0
        let selection = getline(1, "$")
    else
        let selection = getline("'<", "'>")
    end 
    let output = systemlist("vichat split", selection)

    exe "vnew"
    setlocal buftype=nofile nobuflisted

    call append(line('$'), output)
endfunction

function! StartNewChat()
    let pos = getcurpos()
    norm! gg
    let end_of_prompt = search('^USER:', 'n')
    let output = getline(1, end_of_prompt - 1)

    call setpos('.', pos)

    exe "vnew"
    setlocal buftype=nofile nobuflisted

    call append(0, output)
    call setline(line('$'), 'USER: ')
endfunction

command! -buffer Chat call SendToChat()
command! -buffer NewChat call StartNewChat()
command! -buffer Try call TryToChat()
command! -buffer -range Count call CountTokens(<range>)
command! -buffer -range Split call SplitText(<range>)

nnoremap <buffer> <c-s> :Chat<cr>
nnoremap <buffer> <c-t> :Try<cr>
nnoremap <buffer> <c-k> :Count<cr>
nnoremap <buffer> <c-n> :NewChat<cr>:set ft=chat<cr>A
vnoremap <buffer> <c-k> :Count<cr>
nnoremap <buffer> <c-a> GA
nnoremap <buffer> q :q<cr>

nnoremap <buffer> <leader><cr> :Chat<cr>
nnoremap <buffer> <leader>s :Chat<cr>
nnoremap <buffer> <leader>t :Try<cr>
nnoremap <buffer> <leader>a GA

inoremap <buffer> <c-s> <esc>:Chat<cr>
inoremap <buffer> <c-t> <esc>:Try<cr>

if has('mac') || has('macunix')
    vnoremap <buffer> y "+y 
    nnoremap <buffer> yy "+yy 
else
    vnoremap <buffer> y "*y 
    nnoremap <buffer> yy "*yy 
endif
