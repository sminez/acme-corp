#!/bin/bash
# 
# Use Rofi as a fuzzy search of the current window, jumping to the
# selected line and highlighting it.

winid=$(/home/innes/Personal/acme-corp/scripts/afocused)
if [ "$winid" -eq -1 ]; then
    # We need to have seen a focus event before this will work
    exit 1
fi

body=$(9p read acme/"$winid"/body)
# -format d already gives us the line number but this is helpful for quickly
# jumping to line numbers in a vim '$lnum G' style.
numbered=$(echo "$body" | nl -ba -nln -s'| ' -w4)
lnum=$(echo "$numbered" | rofi -dmenu -i -p '/ ' -format d)
echo -n "$lnum" | 9p write acme/"$winid"/addr
echo "dot=addr" | 9p write acme/"$winid"/ctl
echo "show" | 9p write acme/"$winid"/ctl
