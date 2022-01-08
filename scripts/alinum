#!/bin/bash
# Show the current line number in acme

ID=$1
tag=$(9p read acme/$ID/tag)
custom=$(echo $tag | awk -F'|' '{ print $2 }')

echo cleartag | 9p write acme/$ID/ctl
echo -n " Edit =" | 9p write acme/$ID/tag
echo Mx12 19 | 9p write acme/$ID/event
echo cleartag | 9p write acme/$ID/ctl
echo -n "$custom" | 9p write acme/$ID/tag
