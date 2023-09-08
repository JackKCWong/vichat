

function! SendToChat()
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))
    let output[0] = "AI: " . output[0]

    " Append the output of the command to the current buffer
    call append(line('$'), [""] + output + ["", "USER: "])

    norm! G
endfunction

function! TryToChat()
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))

    exe "vnew"
    setlocal buftype=nofile nobuflisted
    call append(line('$'), output)

    norm! G
endfunction

command! -buffer Chat call SendToChat()
command! -buffer Try call TryToChat()
