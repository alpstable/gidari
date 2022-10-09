// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package tools

import (
	"reflect"
	"testing"

	"github.com/alpstable/gidari/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type testCase struct {
	Int int `json:"int"`
}

func TestAssignReadResponseRecords(t *testing.T) {
	t.Parallel()

	t.Run("assign slice", func(t *testing.T) {
		t.Parallel()
		rsp := new(proto.ReadResponse)

		r1, err := structpb.NewStruct(map[string]interface{}{"int": 1})
		if err != nil {
			t.Fatalf("failed to create struct: %v", err)
		}

		rsp.Records = []*structpb.Struct{r1}

		tcase := []*testCase{}
		err = AssignReadResponseRecords(rsp, &tcase)
		if err != nil {
			t.Fatalf("failed to assign records: %v", err)
		}

		expected := []*testCase{{1}}
		if !reflect.DeepEqual(tcase, expected) {
			t.Fatalf("expected %v, got %v", expected, tcase)
		}
	})

	assignSliceBenchmarkResponse := new(proto.ReadResponse)
	assignSliceBenchmarkExpected := []*testCase{}

	for recInt := 1; recInt <= 1e6; recInt++ {
		record, err := structpb.NewStruct(map[string]interface{}{"int": recInt})
		if err != nil {
			t.Fatalf("failed to create struct: %v", err)
		}

		assignSliceBenchmarkExpected = append(assignSliceBenchmarkExpected, &testCase{Int: recInt})

		assignSliceBenchmarkResponse.Records = append(assignSliceBenchmarkResponse.Records, record)
	}
	t.Run("assign slice benchmark", func(t *testing.T) {
		t.Parallel()

		tcase := []*testCase{}
		err := AssignReadResponseRecords(assignSliceBenchmarkResponse, &tcase)
		if err != nil {
			t.Fatalf("failed to assign records: %v", err)
		}
		if !reflect.DeepEqual(tcase, assignSliceBenchmarkExpected) {
			t.Fatalf("expected %v, got %v", assignSliceBenchmarkExpected, tcase)
		}
	})
}
