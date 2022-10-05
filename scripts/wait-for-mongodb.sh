#!/bin/bash

CHECK_CMD="mongosh --eval 'db.runCommand(\"ping\").ok' --quiet --host"

while [ -z `$CHECK_CMD mongo1 2>/dev/null` ] || \
      [ -z `$CHECK_CMD mongo2 2>/dev/null` ] || \
      [ -z `$CHECK_CMD mongo3 2>/dev/null` ]
do
    echo -n "."
    sleep 1
done

echo ""
