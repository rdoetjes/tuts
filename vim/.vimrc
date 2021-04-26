set nocompatible              " be iMproved, required
filetype on                   " required

" set the runtime path to include Vundle and initialize
set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()
" alternatively, pass a path where Vundle should install plugins
"call vundle#begin('~/some/path/here')

" let Vundle manage Vundle, required
Plugin 'VundleVim/Vundle.vim'

" The following are examples of different formats supported.
" Keep Plugin commands between vundle#begin/end.
" plugin on GitHub repo
Plugin 'tpope/vim-fugitive'
" plugin from http://vim-scripts.org/vim/scripts.html
" Plugin 'L9'
" Git plugin not hosted on GitHub
Plugin 'git://git.wincent.com/command-t.git'
" The sparkup vim script is in a subdirectory of this repo called vim.
" Pass the path to set the runtimepath properly.
Plugin 'rstacruz/sparkup', {'rtp': 'vim/'}
" Install L9 and avoid a Naming conflict if you've already installed a
" different version somewhere else.
" Plugin 'ascenator/L9', {'name': 'newL9'}
"
" Add maktaba and codefmt to the runtimepath.
" (The latter must be installed before it can be used.)
Plugin 'google/vim-maktaba'
Plugin 'google/vim-codefmt'
" Also add Glaive, which is used to configure codefmt's maktaba flags. See
" `:help :Glaive` for usage.
Plugin 'google/vim-glaive'

" All of your Plugins must be added before the following line
call vundle#end()            " required
call glaive#Install()
Glaive codefmt plugin[mappings]
Glaive codefmt google_java_executable="java -jar /path/to/google-java-format-VERSION-all-deps.jar"

filetype plugin indent on    " required

"auto format everything, just uncomment this block, you will need to brew
"install or apt install the corresponding formatter like clang-formatter for c
"and yapf for python, vim will tell you what to install! 
" augroup autoformat_settings
"   autocmd FileType bzl AutoFormatBuffer buildifier
"   autocmd FileType c,cpp,proto,javascript,arduino AutoFormatBuffer clang-format
"   autocmd FileType dart AutoFormatBuffer dartfmt
"   autocmd FileType go AutoFormatBuffer gofmt
"   autocmd FileType gn AutoFormatBuffer gn
"   autocmd FileType html,css,sass,scss,less,json AutoFormatBuffer js-beautify
"   autocmd FileType java AutoFormatBuffer google-java-format
"   autocmd FileType python AutoFormatBuffer yapf
"   " Alternative: autocmd FileType python AutoFormatBuffer autopep8
"   autocmd FileType rust AutoFormatBuffer rustfmt
"   autocmd FileType vue AutoFormatBuffer prettier
" augroup END

" To ignore plugin indent changes, instead use:
"filetype plugin on
"
" Brief help
" :PluginList       - lists configured plugins
" :PluginInstall    - installs plugins; append `!` to update or just :PluginUpdate
" :PluginSearch foo - searches for foo; append `!` to refresh local cache
" :PluginClean      - confirms removal of unused plugins; append `!` to auto-approve removal
"
" see :h vundle for more details or wiki for FAQ
" Put your non-Plugin stuff after this line
Plugin 'Valloric/YouCompleteMe'
syntax on
set completeopt-=preview
let g:ycm_autoclose_preview_window_after_insertion = 1
set conceallevel=2
set concealcursor=vin
let g:clang_snippets=1
let g:clang_conceal_snippets=1
let g:clang_snippets_engine='clang_complete'

let mapleader=","
nnoremap <leader>gd :YcmCompleter GoToDeclaration <CR>
nnoremap <leader>gr :YcmCompleter RefactorRename 
nnoremap <leader>gf :YcmCompleter FixIt <CR> 
nnoremap <leader>gdoc :YcmCompleter GetDoc <CR> 

" Complete options (disable preview scratch window, longest removed to aways show menu)
set completeopt=menu,menuone

" Limit popup menu height
set pumheight=20

" SuperTab completion fall-back 
let g:SuperTabDefaultCompletionType='<c-x><c-u><c-p>'

set number
set relativenumber
set tabstop=2
set shiftwidth=2
set softtabstop=2

Plugin 'morhetz/gruvbox'
colorscheme gruvbox
set background=dark    " Setting dark mode

Plugin 'https://github.com/vim-airline/vim-airline'
Plugin 'https://github.com/enricobacis/vim-airline-clock'

Plugin 'preservim/nerdtree'
autocmd VimEnter * NERDTree | wincmd p
augroup AutoChdir
  autocmd!
  autocmd BufEnter * if &buftype !=# 'terminal' | lchdir %:p:h | endif
augroup END
autocmd bufenter * if (winnr("$") == 1 && exists("b:NERDTree") && b:NERDTree.isTabTree()) | q | endif

Plugin 'tpope/vim-commentary'
"use <n>gcc to comment n lines example 10gcc and to uncomment do 10gcc again
"use gcap to comment out a paragraph

"Auto Command to remove trailing spaces on save
autocmd BufWritePre * %s/\s+$//e

"Short keybinding for replace %s/<current word>/<type new>/g the leader is set
"to , above so ,s would write this for you on the command line
:nnoremap <Leader>s :%s/\<<C-r><C-w>\>//g<Left><Left>

set splitbelow splitright
:term ++rows=10
