" Don't write any backup or swap files ever.
" " They're more annoying than they are safe.
set noswapfile
set nowritebackup
set nobackup

" colors
syntax on
set encoding=utf-8

" spaces & tabs
set tabstop=4
set shiftwidth=4
set expandtab
match ErrorMsg '\s\+$'
nnoremap <Leader>rtw :%s/\s\+$//e<CR>

" ui config
" set clipboard=unnamed
set number              " show line numbers
set showcmd             " show command in bottom bar
" set cursorline          " highlight current line
set wildmenu            " visual autocomplete for command menu
set lazyredraw          " redraw only when we need to.
set ruler
set showcmd

" Searching
set incsearch           " search as characters are entered
set hlsearch            " highlight matches

" More reasonable scroll keys
map J <c-e>
map K <c-y>
" Map <leader><leader> to switch to last buffer
nnoremap <leader><leader> <c-^>

" Turn off arrow to be a better person
map <up> <nop>
map <down> <nop>
map <left> <nop>
map <right> <nop>

" leader shortcuts
" " jk is escape
inoremap jk <esc>
" " save session
nnoremap <leader>s :mksession ~/session.vim<CR>
" " save quicker
nnoremap <leader>w :w<CR>
" " turn off search highlight
nnoremap <leader><space> :nohlsearch<CR>

set nocompatible              " required
filetype off                  " required

" set the runtime path to include Vundle and initialize
set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()

" alternatively, pass a path where Vundle should install plugins
"call vundle#begin('~/some/path/here')

" let Vundle manage Vundle, required
Plugin 'gmarik/Vundle.vim'
Plugin 'scrooloose/nerdtree'
Plugin 'kien/ctrlp.vim'
Plugin 'pangloss/vim-javascript'
Plugin 'mustache/vim-mustache-handlebars'
Plugin 'rking/ag.vim'
Plugin 'tpope/vim-sleuth'

" Add all your plugins here (note older versions of Vundle used Bundle instead of Plugin)

" All of your Plugins must be added before the following line
call vundle#end()            " required
filetype plugin indent on    " required

" CtrlP settings
let g:ctrlp_match_window = 'bottom,order:ttb'
let g:ctrlp_switch_buffer = 0
let g:ctrlp_working_path_mode = 0
