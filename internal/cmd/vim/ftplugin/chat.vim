
let s:popup = -1

function! ShowPopupMiddle(message)
    if s:popup == -1
        let popup_options = {
                \ 'padding': [0, 1, 0, 1],
                \ 'border': [],
                \ 'highlight': 'Normal',
                \ }
        
        return popup_create(a:message, popup_options)
    endif

    return s:popup
endfunction

function! SendToChat(stream)
    if getline('$') != ""
        call append(line('$'), [""])
    endif

    " Redirect the content of the current buffer to the external command's stdin
    let cmd = ["vichat", "chat"]
    if a:stream == 1
        let cmd += ["--stream"]
    endif

    let job = job_start(cmd, 
                                \ {
                                \    "in_io": "buffer",
                                \    "in_buf": bufnr(),
                                \    "out_io": "buffer",
                                \    "out_buf": bufnr(),
                                \    "callback": "OnOutputToken",
                                \    "exit_cb": "OnOutputEnd",
                                \ })
    let s = job_status(job)

    if s == "fail"
        call popup_notification("failed to exec vichat")
        return
    endif

    let s:popup = ShowPopupMiddle("thinking")

endfunction

function! OnOutputToken(ch, msg)
    if stridx(a:msg, "AI: ") != 0
        return
    endif

    if s:popup != -1 
        call popup_close(s:popup)
        let s:popup = -1
    endif

    norm! G
endfunction

function! OnOutputEnd(ch, status)
    norm! GA
    exe "w"
endfunction

function! TryToChat()
    let inbuf = bufnr()

    " Redirect the content of the current buffer to the external command's stdin
    exe "vnew"
    setlocal buftype=nofile nobuflisted syntax=markdown 

    let cmd = ["vichat", "chat"]
    let job = job_start(cmd, 
                                \ {
                                \    "in_io": "buffer",
                                \    "in_buf": inbuf,
                                \    "out_io": "buffer",
                                \    "out_buf": bufnr(),
                                \    "callback": "OnOutputToken",
                                \    "exit_cb": "OnOutputEnd",
                                \ })
    let s = job_status(job)

    if s == "fail"
        call popup_notification("failed to exec vichat")
        return
    endif

    let s:popup = ShowPopupMiddle("thinking")
endfunction

function! CountTokens(ran)
    " Redirect the content of the current buffer to the external command's stdin
    if a:ran == 0
        let selection = getline(1, "$")
    else
        let selection = getline("'<", "'>")
    end 
    let output = systemlist("vichat tok", selection)

    echow "estimate: " . output[0] . " tokens in " . bufname()
endfunction

function! ChunkText(ran)
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

function! CloneChat()
    let pos = getcurpos()
    norm! gg
    let end_of_prompt = search('^USER:', 'n')
    let output = getline(1, end_of_prompt - 1)

    call setpos('.', pos)

    exe "vnew"

    call append(0, output)
    call setline(line('$'), 'USER: ')
endfunction

command! -buffer Chat call SendToChat(0)
command! -buffer ChatStream call SendToChat(1)
command! -buffer NewChat call CloneChat()
command! -buffer Try call TryToChat()
command! -buffer -range Count call CountTokens(<range>)
command! -buffer -range Chunk call ChunkText(<range>)

nnoremap <buffer> <c-s> :Chat<cr>
nnoremap <buffer> <c-t> :Try<cr>:set ft=chat<cr>
nnoremap <buffer> <c-k> :Count<cr>
nnoremap <buffer> <c-c> :Chunk<cr>
nnoremap <buffer> <c-n> :NewChat<cr>:set ft=chat<cr>A
vnoremap <buffer> <c-k> :Count<cr>
nnoremap <buffer> <c-a> GA
nnoremap <buffer> q :q<cr>
nnoremap <buffer> <c-q> :q!<cr>

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
