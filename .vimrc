" Don't write any backup or swap files ever.

set noswapfile
set nowritebackup
set nobackup

" colors
syntax on
set encoding=utf-8

" spaces & tabs
set tabstop=4
set shiftwidth=4
" js files are two spaces...
autocmd Filetype javascript setlocal tabstop=2 shiftwidth=2
autocmd Filetype javascriptreact setlocal tabstop=2 shiftwidth=2
autocmd Filetype yaml setlocal tabstop=2 shiftwidth=2
set expandtab
match ErrorMsg '\s\+$'
nnoremap <Leader>W :%s/\s\+$//e<CR>
" python pep8
au BufNewFile,BufRead *.py
    \set tabstop=4
    \set softtabstop=4
    \set shiftwidth=4
    \set textwidth=79
    \set expandtab
    \set autoindent
let python_highlight_all=1

" But wrap text for txt/markdown
autocmd FileType markdown set wrap linebreak textwidth=0
autocmd FileType txt set wrap linebreak textwidth=0

" But not for txt/markdown
autocmd FileType markdown set showbreak=
autocmd FileType txt set showbreak=

" ui config
set number              " show line numbers
set showcmd             " show command in bottom bar
set wildmenu            " visual autocomplete for command menu
set ruler
set backspace=indent,eol,start

" Searching
set incsearch ignorecase smartcase hlsearch " search as characters are entered

" More reasonable scroll keys
map J <c-e>
map K <c-y>

" buffer nav
nnoremap <Tab> :bnext<CR>
nnoremap <S-Tab> :bprevious<CR>
nnoremap <leader><leader> <c-^>

" Turn off arrow to be a better person
map <up> <nop>
map <down> <nop>
map <left> <nop>
map <right> <nop>

" leader shortcuts
" " jk is escape
inoremap jk <esc>
" " save quicker
nnoremap <leader>w :w<CR>
" " turn off search highlight
nnoremap <leader><space> :nohlsearch<CR>

nnoremap <leader>v :set paste<CR>
nnoremap <leader>V :set nopaste<CR>

" vim-plug setup
call plug#begin('~/.vim/plugged')

" Language support and syntax
Plug 'MaxMEllon/vim-jsx-pretty'
Plug 'chr4/nginx.vim'
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
Plug 'heavenshell/vim-jsdoc', { 'for': ['javascript', 'javascript.jsx','typescript'] }
Plug 'mustache/vim-mustache-handlebars'
Plug 'nvie/vim-flake8'
Plug 'pangloss/vim-javascript'
Plug 'smerrill/vcl-vim-plugin'
Plug 'vim-scripts/indentpython.vim'
Plug 'sheerun/vim-polyglot'

" UI and navigation
Plug 'itchyny/lightline.vim'
Plug 'bling/vim-bufferline'

" Git integration
Plug 'tpope/vim-fugitive'
Plug 'airblade/vim-gitgutter'

" Utilities
Plug 'tpope/vim-surround'
Plug 'tpope/vim-commentary'
Plug 'jiangmiao/auto-pairs'
Plug 'madox2/vim-ai'
Plug 'github/copilot.vim'

call plug#end()

filetype plugin indent on

" go settings
let g:go_version_warning = 0

" Perl syntax stuff
autocmd BufNewFile,BufRead *.tt setf tt2

" Lightline settings
let g:lightline = {
      \ 'active': {
      \   'left': [ [ 'mode', 'paste' ],
      \             [ 'fugitive' ] ]
      \ },
      \ 'component_function': {
      \   'fugitive': 'LightLineFugitive',
      \ }
      \ }
function! LightLineFugitive()
  return exists('*fugitive#head') ? fugitive#head() : ''
endfunction

function! BFGhUrl()
    " Get current branch
    let branch = substitute(system("git rev-parse --abbrev-ref HEAD"), '\n', '', '')

    " Get git repository root
    let git_root = substitute(system("git rev-parse --show-toplevel"), '\n', '', '')

    " Get remote URL and convert to GitHub URL
    let remote_url = substitute(system("git config --get remote.origin.url"), '\n', '', '')

    " Convert SSH/HTTPS git URLs to GitHub web URLs
    if remote_url =~ '^git@github\.com:'
        let github_url = substitute(remote_url, 'git@github\.com:', 'https://github.com/', '')
        let github_url = substitute(github_url, '\.git$', '', '')
    elseif remote_url =~ '^https://github\.com/'
        let github_url = substitute(remote_url, '\.git$', '', '')
    else
        echo 'Error: Not a GitHub repository'
        return
    endif

    " Build the URL
    let git_start = github_url . '/blob/' . branch
    let line_frag = '#L'. line('.')
    let rel_path = substitute(expand('%:p'), git_root, '', '')
    let full_url = git_start . rel_path . line_frag

    call system('open ' . full_url)
    echo 'opening... ' . full_url
endfunction

nnoremap <leader>o :call BFGhUrl()<cr>

" complete text on the current line or in visual selection
nnoremap <leader>a :AI<CR>
xnoremap <leader>a :AI<CR>

" edit text with a custom prompt
xnoremap <leader>s :AIEdit fix grammar and spelling<CR>
nnoremap <leader>s :AIEdit fix grammar and spelling<CR>

" trigger chat
xnoremap <leader>c :AIChat<CR>
nnoremap <leader>c :AIChat<CR>

" Override write command in AI chat buffers
autocmd FileType aichat nnoremap <buffer> <leader>w :AIChat<CR>
