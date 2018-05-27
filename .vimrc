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
" set clipboard=unnamed
set number              " show line numbers
set showcmd             " show command in bottom bar
" set cursorline          " highlight current line
set wildmenu            " visual autocomplete for command menu
set ruler

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

set nocompatible              " required
filetype off                  " required

" set the runtime path to include Vundle and initialize
set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()

" alternatively, pass a path where Vundle should install plugins
"call vundle#begin('~/some/path/here')

" let Vundle manage Vundle, required
Plugin 'gmarik/Vundle.vim'
Plugin 'pangloss/vim-javascript'
Plugin 'mustache/vim-mustache-handlebars'
Plugin 'hallison/vim-markdown'
Plugin 'tpope/vim-fugitive'
Plugin 'bling/vim-bufferline'
Plugin 'itchyny/lightline.vim'
Plugin 'heavenshell/vim-jsdoc'
Plugin 'nvie/vim-flake8'
Plugin 'vim-scripts/indentpython.vim'
" Add all your plugins here (note older versions of Vundle used Bundle instead of Plugin)

" All of your Plugins must be added before the following line
call vundle#end()            " required
filetype plugin indent on    " required

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
    let full_path = expand('%')
    let bf_path = matchstr(full_path, '/opt/buzzfeed/')
    if empty(bf_path)
        echo 'not bf repo'
        return
    endif
    let bf_repo = matchlist(full_path, '^/opt/\zs\v(buzzfeed/(mono|buzzfeed))(/.*)')[1]
    let bf_file = matchlist(full_path, '^/opt/\zs\v(buzzfeed/(mono|buzzfeed))(/.*)')[3]
    let bf_branch = system("git rev-parse --abbrev-ref HEAD | sed 's/^ *//;s/ *$//'")
    let bf_gh = 'https://github.com/'
    let bf_gh_url = 'https://github.com/' . bf_repo . '/blob/' . matchstr(bf_branch, '\v\w+') . bf_file . '#L' . line('.')
    call system('open ' . bf_gh_url)
    echo 'opening... ' . bf_gh_url
endfunction

nnoremap <leader>o :call BFGhUrl()<cr>
