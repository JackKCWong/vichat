
function! SendToChat()
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))
    let output[0] = "AI: " . output[0]

    " Append the output of the command to the current buffer
    call append(line('$'), [""] + output + ["", "USER: "])

    norm! G
    norm! A
endfunction

function! TryToChat()
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))

    exe "vnew"
    setlocal buftype=nofile nobuflisted syntax=markdown

    call append(line('$'), output)
endfunction

command! -buffer Chat call SendToChat()
command! -buffer Try call TryToChat()

nnoremap <buffer> <c-s> :Chat<cr>
nnoremap <buffer> <c-t> :Try<cr>
nnoremap <buffer> <leader><cr> :Chat<cr>
nnoremap <buffer> <leader>s :Chat<cr>
nnoremap <buffer> <leader>t :Try<cr>

inoremap <buffer> <c-s> <esc>:Chat<cr>
inoremap <buffer> <c-t> <esc>:Try<cr>
