#! /usr/bin/env bash
# Use pick to search through the current active window, jumping to the
# selected line and highlighting it.

winid="$(echo "active / ." | nc localhost 2009)"
lnum=$(pick -n -N)
echo -n "$lnum" | 9p write acme/"$winid"/addr
echo "dot=addr" | 9p write acme/"$winid"/ctl
echo "show" | 9p write acme/"$winid"/ctl