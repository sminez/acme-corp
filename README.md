ACME-CORP: Utilities and extension programs for the acme text editor
====================================================================

This repo is a collection of scripts, extension programs and helpers for use
with the [acme][0] text editor from [plan9][1]. I highly recommend installing the
wonderful [plan9 from userspace][2] port of plan9 software to unix (optionally
using [my fork][3] for a few tweaks).


### A mouse based text editor? The horror!
I love Vim. It's been my daily text editor for years now and I am very
productive in it. But extending Vim myself, in any way more complicated than
canning a macro for later use, is something that I find annoyingly fiddly.

Why? Vimscript...

Vimscript is horrible. Now, don't get me wrong: I've also made serious efforts
to use Emacs as well and my distaste for elisp is equally strong. So, in the
grand Unix tradition of bolting together the tools that do what you need I have
settled on an interesting hybrid:

Primary editing / complex text manipulation happens in Vim.

Editing session management and coordination happens through acme. Why? Because
of the amazing power it gives to plain text.


### This sounds overly complicated...
Yeah. It is. I'm working on something on the side to make this a bit easier to
work with as a dail driver but overall I like acme for day to day 'session
management'.


  [0]: http://acme.cat-v.org/
  [1]: https://9p.io/plan9/
  [2]: https://9fans.github.io/plan9port/
  [3]: https://github.com/sminez/plan9port
