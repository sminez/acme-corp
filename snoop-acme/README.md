snoop-acme - automate your acme experience
==========================================

I love tinkering with stuff, which shouldn't be a surprise given that I'm
writing a bunch of tooling around acme. I'm also pretty forgetful however so I
like to be able to automate as much as I can on both my work and home machines.
The ability to manipulate editor state via a virtual file system is _wonderful_
and a true joy for someone who likes to poke at stuff. It does mean that I need
to remember all of the tooling that I'm writing though which is less good...

Previously I have always used Vim (I still use it as my daily driver at work to
be honest) but with the workflow of living in a terminal, having a LOT of
aliases, shell functions and utility scripts dotted around the place and making
use of programs such as [fzf][0] and [wd][1] to help me navigate things. I'd
like to be able to pull my current tooling into acme itself and get a proper
workflow sorted out that lets me make the edits I need to make, manage external
systems and APIs and keep track of what I am doing.

Enter acme-corp and the snooper.

The snopper is a local webserver currently. I know, I know: it would be a lot
more in the spirit of things to have it mount its own 9fs file system to do all
of this but I _do_ know how to write webservers, I don't (currently) know how to
get 9fs working in a nice way that doesn't require `9p` to work so there's not a
huge difference either way. The server comes with several utility scripts that
are essentially canned requests that set state to modify some tooling I've
written:
  - Enable / disable format on save for all windows.
  - Allow for canned "you X-clicked on this text in the tag" events to be
  injected, so that they can be bound to keys using something like shkd or my
  compiled in dwm config. (Indent the current paragraph, wrap lines etc...)
  - I'd _really_ like to get a kind of vim-style `ex` mode with a pop up command
  line via a key binding. My laptops don't have decent three button mouse
  emulation which means modifiers and frantic swiping all over the place to get
  to the tag to make a simple `,x/#/ c/;;/` style comment replacement edit. When
  I have access to my _real_ mouse this isn't a problem but I want to get a "run
  this Edit command" pop up working ASAP: a lot of my edits in vim are done
  using `:%s/.../.../g` so if I can get a replacement for that which I'm happy
  with I'll be a very happy Glenda.
