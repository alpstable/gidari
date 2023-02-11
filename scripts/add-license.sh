#!/bin/bash

# EXCLUDE_LIST is an array of files paths to exclude from prepending with a license notice.
declare -a EXCLUDE_LIST=(
    "./internal/web/auth/auth1.go",
    "./proto/db.pb.go"
)

YEAR=$(date +%Y)

# LICENSE_TEMPLATE is the license notice to prepend to files.
LICENSE_TEMPLATE=$(cat <<EOF
// Copyright $YEAR The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
EOF
)

# NEWLINE is added in all files except doc.go.
NEWLINE=$(cat <<EOF
//
EOF
)

# append all *.go in this repository with the LICENSE_TEMPLATE
for file in $(find . -name "*.go" -type f); do
    	# skip files in the EXCLUDE_LIST
    	if [[ " ${EXCLUDE_LIST[@]} " =~ " ${file} " ]]; then
		continue
    	fi

      # If the file starts with "// Copied from" then we don't want to prepend
      # the license.
      if grep -q "^// Copied from" "${file}"; then
        continue
      fi

      sed -i '/^package/,$!d' "${file}"

    	# prepend the LICENSE_TEMPLATE to the file
    	if [[ "${file}" == ./doc.go ]]; then
          printf "%s\n" "${LICENSE_TEMPLATE}" | cat - "${file}" > /tmp/out
        else
          printf "%s\n%s\n" "${LICENSE_TEMPLATE}" "${NEWLINE}" | cat - "${file}" > /tmp/out
      fi
    	mv /tmp/out "${file}"

done


