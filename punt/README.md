punt - a quick detour from the land of acme
===========================================

Can't bare to leave your favourite editor behind? Punt your current window
contents to it, make your edits and then load it back into acme. Or,
alternatively, punt the file you're working on to an external program that can
do something with it: updating some HTML? Punt it to the browser to see what it
looks like.

### How does it work?
`punt` will copy your current acme window content (_not_ the underlying file!)
to a temporary file on your local file system before passing that file path off
to the external program that you designate. When the external program closes,
the contents of that temporary file will be read back in to acme and overwrite
the window that you were in before. Initially I wrote this as a way to quickly
open up my current buffer (window, old habits die hard) back to vim if there was
something I needed to do that I couldn't wrap my head around in acme. As I'm
getting more comfortable with acme, `punt` is less of a "must have" tool and
more of "useful in unexpected situations".
