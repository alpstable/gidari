#!/bin/bash

# EXCLUDE_LIST is an array of files paths to exclude from prepending with a license notice.
declare -a EXCLUDE_LIST=(
    "./internal/web/auth/auth1.go",
    "./proto/db.pb.go"
)

# LICENSE_TEMPLATE is the license notice to prepend to files.
LICENSE_TEMPLATE=$(cat <<EOF
// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
EOF
)

# append all *.go in this repository with the LICENSE_TEMPLATE
for file in $(find . -name "*.go" -type f); do
    	# skip files in the EXCLUDE_LIST
    	if [[ " ${EXCLUDE_LIST[@]} " =~ " ${file} " ]]; then
		continue
    	fi

    	# skip files that already have the LICENSE_TEMPLATE
    	if grep -q "Copyright 2022 The Gidari Authors." "${file}"; then
		continue
    	fi

	echo "${LICENSE_TEMPLATE}" | cat - "${file}" > /tmp/out && mv /tmp/out "${file}"
done


