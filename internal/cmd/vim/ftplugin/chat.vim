function! ShowPopupMiddle(message)
  let view_width = winwidth(0)
  let view_height = winheight(0)
  
  let popup_width = len(a:message) + 4
  let popup_height = 3
  
  let row = (view_height - popup_height) / 2
  let col = (view_width - popup_width) / 2
  
  let popup_options = {
        \ 'line': row,
        \ 'col': col,
        \ 'padding': [0, 1, 0, 1],
        \ 'border': [],
        \ 'highlight': 'Normal',
        \ }
  
  return popup_create(a:message, popup_options)
endfunction


function! SendToChat()
    let pid = ShowPopupMiddle('thinking...')
    redraw
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))

    let output[0] = "AI: " . output[0]

    call popup_close(pid)

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
    let pid = ShowPopupMiddle('thinking...')
    redraw
    " Redirect the content of the current buffer to the external command's stdin
    " Redirect the content of the current buffer to the external command's stdin
    let output = systemlist("vichat chat", getline(1, '$'))

    call popup_close(pid)

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
command! -buffer -range Chunk call ChunkText(<range>)

nnoremap <buffer> <c-s> :Chat<cr>
nnoremap <buffer> <c-t> :Try<cr>
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
