acme itself
-----------
web links:  (http://doc.cat-v.org/plan_9/4th_edition/papers/acme/)
            (http://doc.cat-v.org/plan_9/4th_edition/papers/plumb)
local docs: acme(1)
            acme(4)
            sam(1)
            regexp(7)
            9p(1)
            9pclient(3)

### Personal motivation for using acme
`acme` really is, as described by Rob Pike in what I think is his original
paper on it, "a user interface for programmers". Since finding it I've been
slow to realise that it is the missing component (along with the plumber)
between writing your own shell scripts and utility programs and having an
editing / work environment that you are in complete control over. I have,
over time, attempted to use a wide range of editors / IDEs as my daily go-to
tool: light-table, sublime-text, atom, intellij, emacs, spacemacs, vim, neovim,
idle, kate, gedit...you get the idea. Ultimately, for the majority of my time
working as a programmer I have settled on vim from within the terminal rather
than gVim or some sort of electron skin.
The main reasons for this are twofold:
  1) vim's editing language is simply wonderful for structurally editing files
     and it can be extended with your own macros and key bindings _relatively_
     simply. (Writing plugins on the other hand?...)
  2) being able to quickly jump in and out of files while navigating around my
     file system and then drop back out to the shell when needed is a workflow
     that I find infinitely more flexible than hunting around in a gui menu
     system for a hand-rolled git implementation / terminal emulator / kitchen
     sink.

So, based on that, why acme? Well, while it does feel like I need to give up
vim's file navigation and structural editing (see footnote) by moving to acme,
what I gain is the ability to modify file contents and editor state via _any
language I chose_ and, more than that, tail the editor state and inject my
own actions from helper programs.

Now what self respecting hacker (in the original sense) can say no to that?

(footnote: yes I know that acme has `sam` structural editing embedded but the
need to jump up to the tag to type a command rather than bind common commands
to key sequences/chords is a real pain.)


### Some useful acme commands
* Select with button 3:
  * `:n`
    * jump to (and select) line n of the file
  * `:0`
    * jump to the start of the file
  *`:$`
    * jump to the end of the file
  * `:,`
    * select the entire file

* Select with button 2
  * `Edit ,`
    * select the entire file
  * `Edit , d`
    * clear the window
  * `Edit , < SOME_EXTERNAL_COMMAND`
    * replace body with the output of the external command
  * `Edit , > SOME_EXTERNAL_COMMAND`
    * pipe the current window content to an external command
  * `Edit ,s/foo/bar/g`
    * equivalent to `:%s/foo/bar/g` in vim

* `$%` is the current window name (normally the abspath of the file being edited)
* `$winid` is the current window ID