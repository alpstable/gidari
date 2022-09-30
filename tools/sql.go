// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
package tools

import (
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
)

// SQLIterativePlaceholders will return a string of placeholders that require a number next to the placeholder string
// that iteratively increases by the number of arguments passed to the query. For example, if a string has numCols=3
// and numRows=2, this function will return "(?1,?2,?3),(?4,?5,?6)".
func SQLIterativePlaceholders(numCols int, numRows int, symbol string) string {
	if numCols == 0 || numRows == 0 {
		return "()"
	}

	if symbol == "" {
		symbol = "?"
	}

	var strBldr strings.Builder

	for pos := 0; pos < numRows*numCols; pos++ {
		if pos%numCols == 0 {
			strBldr.WriteString("(")
		}

		strBldr.WriteString(symbol)
		strBldr.WriteString(strconv.Itoa(pos + 1))

		if pos%numCols == numCols-1 {
			strBldr.WriteString(")")

			if pos != numRows*numCols-1 {
				strBldr.WriteString(",")
			}
		} else {
			strBldr.WriteString(",")
		}
	}

	return strBldr.String()
}

// SQLFlattenPartition will take a slice of structures, extract data from their fields, and append it to a slice.
// This will "flatten" the data to be used in conjunctino with placeholders in a SQL query.
func SQLFlattenPartition(columns []string, partition []*structpb.Struct) []interface{} {
	var args []interface{}

	for _, record := range partition {
		hash := record.AsMap()
		for _, column := range columns {
			args = append(args, hash[column])
		}
	}

	return args
}
