#! /usr/bin/env bash
# Use pick to provide a quick command line a-la vim's `:` command mode
#
# Supported actions:
#   `e<file name>`    -- attempt to open <file name> in a new acme window
#   `?<input>`        -- button 3 (search) in $winid tag
#   `!<input>`        -- button 2 (execute) in $winid tag
#   `<edit command>`  -- treat input as an argument to the built in `Edit` command


# We expect to be called via a user configured global keybinding so we wont have
# a $winid env var to read for targetting the correct window so ping the snooper
# to find out what the active window is. (If the snooper returns -1 for the id
# of the window then we have not been able to determine focus yet so exit)
WINID="$(echo "active / ." | nc localhost 2009)"
(( WINID > 0 )) || exit 1


# == helper functions ==
# see acme(4) for details on the control files and event structure
# NOTE: character offsets are zero indexed

function windowDirectory() {
  dirname "$(9p read "acme/$WINID/tag" | head -1 | cut -d' ' -f1)"
}

function writeClickEvent() {
  local eventType=$1 offset=$2 end=$3

  echo "M$eventType$offset $end"
  echo "M$eventType$offset $end" | 9p write "acme/$WINID/event"
}

function setTag() {
  local text=$1

  echo "cleartag" | 9p write "acme/$WINID/ctl"
  echo -n "$text" | 9p write "acme/$WINID/tag"
  echo "clean" | 9p write acme/"$WINID"/ctl
}

function spoofClickInTag() {
  local eventType=$1 rawText=$2

  text="$(echo "$rawText" | xargs)"
  textLen="${#text}"
  echo "clean" | 9p write acme/"$WINID"/ctl
  fullTag="$(9p read "acme/$WINID/tag")"
  original="$(echo "$fullTag" | cut -d'|' -f2 | xargs)"
  nChars="$(echo "$fullTag" | cut -d'|' -f1 | wc -c)"
  offset="$(( nChars + 1 ))"
  end=$(( offset + textLen + 1 ))

  setTag " $text"
  writeClickEvent "$eventType" "$offset" "$end"
  echo "show" | 9p write acme/"$WINID"/ctl
  setTag " $original"
}


# == actions ==

function openInAcme() {
  local fname
  fname="$(echo "$1" | xargs)"

  case $fname in
    \~/*) fname="$HOME/${fname:2}";;
      /*) fname=$fname;;
       *) fname="$(windowDirectory)/$fname";;
  esac

  echo "name $fname" | 9p write acme/new/ctl
  id="$(9p read acme/index | grep "$fname" | xargs | cut -d' ' -f1)"
  echo 'get' | 9p write "acme/$id/ctl"
}

function searchFromTag() { spoofClickInTag "l" "$1"; }
function executeFromTag() { spoofClickInTag "x" "$1"; }
function runEditCommand() { executeFromTag "Edit $1"; }

# == main ==

# Open a new pick window for the user to provide their command, if pick exits
# with no output then the user hit enter with nothing typed or we were closed.
# Either way, there is nothing else to do.
input="$(echo "" | pick -s -p ':')"
[[ -n "$input" ]] || exit 0

# If we got some input, try to parse and action it
case $input in
  \?*) searchFromTag "${input:1}";;
   !*) executeFromTag "${input:1}";;
   e*) openInAcme "${input:1}";;
    *) runEditCommand "$input";;
esac