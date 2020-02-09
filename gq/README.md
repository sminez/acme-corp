gq - line wrapping with prefixes
================================

A common edit I make in vim is to run `gq` on a visual selection, or `gqap` for
a block, in order to quickly line wrap comments or formatted markup. There is
the GNU core-utils `fmt` program but you need to specify the prefix yourself and
it doesn't always do what I want. `gq` will wrap lines to 80 characters by
default with the option to set the column count using the `-c` flag. Prefixes
are maximally determined and stop on the first alpha-numeric character unless
the `-a` flag is given.

### Use within acme
Until I get this hooked in to the rest of acme-corp, `gq` is simply intended to
be used in the tag as a filter for you to pipe a selection through.
- For example, if you add `|gq -c 100` to the tag, select a comment block with
  button 1 and then select the command with button 2, you will wrap your comment
  to 100 characters.

### Known bugs
- Leading whitespace will be stripped.
- Indentation within a comment block will be lost.

### TODO
- Safe auto wrapping for files that don't currently have a formatter available
  for use with `afmt` (shell for example).
  - This means no breaking of pre-existing indentation or stripping leading
  whitespace.
- Hook into `snoop-acme` as a triggerable key binding (will need shkd or
  similar, such as my compiled dwm bindings) in order to get "wrap the current
  paragraph" or "wrap current selection".
  - This will involve storing some state in the snooper so that it can select
  the current window / text block / maximal text block with matching prefix etc
