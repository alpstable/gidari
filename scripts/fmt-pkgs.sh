#!/bin/bash

# Loop through all of the packages from list-pkgs.sh and run the gofmt command on them
PKGS=$(scripts/list-pkgs.sh)
for pkg in $PKGS; do
    echo "gofmt -w $pkg"
    gofmt -w -l -s $pkg
done

