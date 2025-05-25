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
" set clipboard=unnamed
set number              " show line numbers
set showcmd             " show command in bottom bar
" set cursorline          " highlight current line
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
Plug 'junegunn/fzf', { 'do': { -> fzf#install() } }
Plug 'junegunn/fzf.vim'

" Git integration
Plug 'tpope/vim-fugitive'
Plug 'airblade/vim-gitgutter'

" Code completion and LSP
Plug 'neoclide/coc.nvim', {'branch': 'release'}

" Utilities
Plug 'tpope/vim-surround'
Plug 'tpope/vim-commentary'
Plug 'jiangmiao/auto-pairs'
Plug 'madox2/vim-ai'

call plug#end()

filetype plugin indent on

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
    let branch = substitute(system("git rev-parse --abbrev-ref HEAD"), '\n', '', '')
    let git_start = 'https://github.com/buzzfeed/mono/blob/' . branch
    let line_frag = '#L'. line('.')
    let full_url = substitute(expand('%:p'), '/opt/buzzfeed/mono', git_start, '') . line_frag
    call system('open ' . full_url)
    echo 'opening... ' . full_url
endfunction

nnoremap <leader>o :call BFGhUrl()<cr>

" FZF mappings
nnoremap <leader>f :Files<CR>
nnoremap <leader>b :Buffers<CR>
nnoremap <leader>g :Rg<CR>
nnoremap <leader>t :Tags<CR>

" CoC configuration
inoremap <silent><expr> <TAB>
      \ pumvisible() ? "\<C-n>" :
      \ <SID>check_back_space() ? "\<TAB>" :
      \ coc#refresh()
inoremap <expr><S-TAB> pumvisible() ? "\<C-p>" : "\<C-h>"

function! s:check_back_space() abort
  let col = col('.') - 1
  return !col || getline('.')[col - 1]  =~# '\s'
endfunction

" Use <cr> to confirm completion
inoremap <expr> <cr> pumvisible() ? "\<C-y>" : "\<C-g>u\<CR>"

" Use `[g` and `]g` to navigate diagnostics
nmap <silent> [g <Plug>(coc-diagnostic-prev)
nmap <silent> ]g <Plug>(coc-diagnostic-next)

" GoTo code navigation
nmap <silent> gd <Plug>(coc-definition)
nmap <silent> gy <Plug>(coc-type-definition)
nmap <silent> gi <Plug>(coc-implementation)
nmap <silent> gr <Plug>(coc-references)

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

function! WB2BF()
    silent! %s/IF NOT EXISTS //g
    silent! g/ENGINE = InnoDB/delete
    silent! g/^--.*/delete
    silent! %s/`mydb`.//g
    silent! %s/`default_schema`.//g
    silent! %s/;/|/g
endfunction
