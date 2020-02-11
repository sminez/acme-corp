dirtree - a simple file tree for acme
=====================================

Dirtree is a minimal, tree based file explorer for acme. (I don't like the way
that acme's default handling of folders opens a new folder for each directory
you click on.)

When you run dirtree it will open a new window displaying the current working
directory. If you pass a directory as an argument it will use that as the
starting directory.

The current directory is shown at the top of the window, wrapped in parens for
easy 1-1-3 clicking to highlight, then open the directory in the normal acme
style (in case you want to do file operations with the listing).


### Tag commands
- Hidden
  - Toggle the display of dotfiles in the current file tree

- UpDir
  - Move the root of the file tree up to the parent directory then redraw the
  window.

- Reset
  - Collapse all expanded directories and re-fetch their contents


### Mouse Actions
- Button 2
  - directory: Set that as the new root directory of the tree and redraw

- Button 3
  - directory: Toggle the expand / collapse of the directory contents
  - file: Open that file via using the plumber
  - user typed text: attempt to execute in the shell as per normal acme windows.
    - NOTE: the directory used for execution will be the current root of the tree.


### Known Bugs
- Spaces in file names _sometimes_ causes the plumber to fail...not sure why.
- The entire tree is redrawn on each expansion / collapse of a node. If you have
  expanded a lot of nodes then you will see some noticeable redraw.