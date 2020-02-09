Library code for interacting with acme from go
----------------------------------------------
source:     (https://github.com/9fans/go)
local docs: (go doc 9fans.net/go/acme)
            (go doc 9fans.net/go/draw)
            (go doc 9fans.net/go/plan9)
            (go doc 9fans.net/go/plumb)

9fans is (mostly) the work of Russ Cox and Rob Pike, wrapping the the virtual
file system api presented by acme along with interaction with the plumber and
other plan9 system stuff. There is also some work around `draw` which I've not
experimented _too_ much with yet but seems interesting for kicking off some
of your own GUI plan9 goodness if that's your thing.

This is what is powering pretty much all of `acme-corp`, so thanks!