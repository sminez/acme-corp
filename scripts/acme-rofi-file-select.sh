#!/usr/bin/env bash
# Helper script called by acme to fuzzy search for files and open them in acme.
#
# If the selected path is a file then a new acme window is opened containing the
# file. If it is a directory then the dirtree program is run to display an
# interactive file tree.
#
# Modified from: https://github.com/carnager/rofi-scripts/tree/master/rofi-finder

if [ ! -z "$@" ]
then
  QUERY=$@
  if [[ "$@" == /* ]]
  then
    if [[ "$@" == *\?\? ]]
    then
      if [ -d "${QUERY%\/* \?\?}" ]; then
        coproc ( xplor "${QUERY%\/* \?\?}"  > /dev/null 2>&1 )
      else
        coproc ( editinacme "${QUERY%\/* \?\?}"  > /dev/null 2>&1 )
      fi

      exec 1>&-
      exit;

    else
      if [ -d "$@" ]; then
        coproc ( xplor "$@"  > /dev/null 2>&1 )
      else
        coproc ( editinacme "$@"  > /dev/null 2>&1 )
      fi

      exec 1>&-
      exit;
    fi

  elif [[ "$@" == \!\!* ]]
  then
    find / -iname *"${QUERY#!!}"* 2>&1 | grep -v 'Permission denied\|Input/output error'

  elif [[ "$@" == \?* ]]
  then
    while read -r line; do
      echo "$line" \?\?
    done <<< $(find $HOME -iname *"${QUERY#\?}"* 2>&1 | grep -v 'Permission denied\|Input/output error')

  else
    find $HOME -iname *"${QUERY#!}"* 2>&1 | grep -v 'Permission denied\|Input/output error'
  fi

else
  # Display some some help text
  echo ">> Type your search query to find files"
  echo ">> To seach again type !<SEARCH_QUERY>"
  echo ">> To seach parent directories type ?<SEARCH_QUERY>"
  echo '>> Search from / instead of $HOME using !!'
fi
