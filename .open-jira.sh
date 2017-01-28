#!/usr/bin/env bash
echo 'heyo'
# get current open browser command
case $( uname -s ) in
  Darwin)  open=open;;
  MINGW*)  open=start;;
  CYGWIN*) open=cygstart;;
  MSYS*)   open="powershell.exe –NoProfile Start";;
  *)       open=${BROWSER:-xdg-open};;
esac

# get current branch
if [ -z "$2" ]; then
  branch=$(git symbolic-ref -q --short HEAD)
else
  branch="$2"
fi

card=$( echo $branch | grep -o 'BF-\d*')

$open "https://buzzfeed.atlassian.net/browse/$card" &> /dev/null
