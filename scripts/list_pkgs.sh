#!/usr/bin/env bash

directory="$1"
if [ -z "$directory" ]; then
    directory="."
fi
go list $directory/...
