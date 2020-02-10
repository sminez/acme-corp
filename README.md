ACME-CORP: Utilities and extension programs for the acme text editor
====================================================================

This repo is a collection of scripts, extension programs and helpers for use
with the [acme][0] text editor from [plan9][1]. I highly recommend installing the
wonderful [plan9 from userspace][2] port of plan9 software to unix, optionally
using [my fork][3] for a few tweaks such as additional key bindings for acme
and a suckless style `config.h` file for setting a custom color scheme and
default fonts.

Each top level directory is a stand alone go binary that should be possible to
run with a vanilla acme installation independently from one another. That said,
the intention is to run all of these as a suite of programs that compliment one
another.


### Inspired by Suckless
I love the [suckless][4] programs. On my home machine I run dwm, st and dmenu
and I have my own [patched versions][5] of all three that I tinker with from time
to time. I'm pretty satisfied with my current setup: despite running on an old,
cheap HP laptop I get pretty snappy performance and I'm yet to see any random
crashes (other than while I'm actively tinkering of course).

With that in mind, acme-corp is written to be as simple as possible while still
being readable. So, little to no magic hacks, comments and links around any of
the more esoteric pieces of code (in particular the use of `sam` expressions
to manipulate window content) and configuration by modifying the source. As much
as I love writing parsers and command line utilities, it really is far easier to
just pull simple things from environment variables and hard code the rest.


### Installation
As far as I can tell, `go get` -ing this repo and running `go install ./...` in
the root should be all you need to grab the utilities themselves. For the scripts,
you will need to add them to your `$PATH` (I tend to conditionally add them when
starting acme so that I don't clutter up my `$PATH`). Some of the scripts are
simple triggers for the `snooper` so you will need to kick that off in order for
them to do anything. I haven't provided any sort of install script for the
dependencies as they tend to fluctuate a bit as I work on things. Make sure to
read the source of anything you intend to run to see what it expects to be on
your system.


### Current utilities
For more in depth information, please see the individual README files in each
directory and obviously the source code itself.

* dirtree
  * A directory viewer for acme. The built in support for navigating the filesystem
  can quickly get out of hand if you are jumping around directories a lot. This
  allows you to have a single window that acts as a file tree, allowing you to
  move the root when needed.

 * gq
   * Mimic the Vim `gq` key sequence. Yes I know that `fmt` exists but I wanted to
   write something that did what I wanted out of the box. Essentially this is wrap
   lines to a column count (default 80) and preserve any common prefix that is found
   in order to give a language agnostic way to tidy up comment blocks.

 * punt
   * Quickly open the current window in an external program. Originally intended
   for punting things over to Vim if I needed to do some more complicated editing
   that I can currently get my head around using `sam` commands. It can also send
   window content to GUI based programs (default is to expect that the program
   will run in a terminal) via the `-g` flag.

 * snoop-acme
   * A local TCP server that tails all acme event streams and runs hooks such as
   auto formatting and linting (for known filetypes and tools) and quickly allowing
   programs to access the currently focused window. As far as I can tell, the
   latter is simply emitted on the event log rather than being a piece of state
   that you can pull from the acme virtual file system. My main reason for wanting
   to get at this is to drive things like the `acme-fuzzy-window-search.sh`
   script.


  [0]: http://acme.cat-v.org/
  [1]: https://9p.io/plan9/
  [2]: https://9fans.github.io/plan9port/
  [3]: https://github.com/sminez/plan9port/
  [4]: https://suckless.org/
  [5]: https://github.com/sminez/suckless/
