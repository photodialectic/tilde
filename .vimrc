" change
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

" But wrap text for txt/markdown
autocmd FileType markdown set wrap linebreak textwidth=0
autocmd FileType txt set wrap linebreak textwidth=0

" If wrapping is enabled, mark wrapped lines
" set showbreak=\ ->\

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
" " save quicker
nnoremap <leader>w :w<CR>
" " turn off search highlight
nnoremap <leader><space> :nohlsearch<CR>
" " boot up nerdtree
map <leader>n :NERDTreeToggle<CR>

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
Plugin 'terryma/vim-multiple-cursors'
Plugin 'tpope/vim-fugitive'
Plugin 'hallison/vim-markdown'
Plugin 'tclem/vim-arduino'

" Add all your plugins here (note older versions of Vundle used Bundle instead of Plugin)

" All of your Plugins must be added before the following line
call vundle#end()            " required
filetype plugin indent on    " required

" CtrlP settings
let g:ctrlp_match_window = 'bottom,order:ttb'
let g:ctrlp_switch_buffer = 0
let g:ctrlp_working_path_mode = 0

" vim-ardunio settings
" Default: /Applications/Arduino.app/Contents/Resources/Java
" let g:vim_arduino_library_path = /path/to/arduino/installation
" Default: result of `$(ls /dev/tty.* | grep usb)`
" let g:vim_arduino_serial_port = /my/serial/port
