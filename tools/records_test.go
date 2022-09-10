package tools

import (
	"reflect"
	"testing"

	"github.com/alpine-hodler/gidari/pkg/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type testcase struct {
	Int int `json:"int"`
}

func TestAssignReadResponseRecords(t *testing.T) {
	t.Run("assign slice", func(t *testing.T) {
		rsp := new(proto.ReadResponse)

		r1, err := structpb.NewStruct(map[string]interface{}{"int": 1})
		if err != nil {
			t.Fatalf("failed to create struct: %v", err)
		}

		rsp.Records = []*structpb.Struct{r1}

		tc := []*testcase{}
		err = AssignReadResponseRecords(rsp, &tc)
		if err != nil {
			t.Fatalf("failed to assign records: %v", err)
		}

		expected := []*testcase{{1}}
		if !reflect.DeepEqual(tc, expected) {
			t.Fatalf("expected %v, got %v", expected, tc)
		}
	})

	assignSliceBenchmarkResponse := new(proto.ReadResponse)
	assignSliceBenchmarkExpected := []*testcase{}
	for i := 1; i <= 1e6; i++ {
		r, err := structpb.NewStruct(map[string]interface{}{"int": i})
		if err != nil {
			t.Fatalf("failed to create struct: %v", err)
		}
		assignSliceBenchmarkExpected = append(assignSliceBenchmarkExpected, &testcase{Int: i})
		assignSliceBenchmarkResponse.Records = append(assignSliceBenchmarkResponse.Records, r)
	}
	t.Run("assign slice benchmark", func(t *testing.T) {
		tc := []*testcase{}
		err := AssignReadResponseRecords(assignSliceBenchmarkResponse, &tc)
		if err != nil {
			t.Fatalf("failed to assign records: %v", err)
		}
		if !reflect.DeepEqual(tc, assignSliceBenchmarkExpected) {
			t.Fatalf("expected %v, got %v", assignSliceBenchmarkExpected, tc)
		}
	})
}
