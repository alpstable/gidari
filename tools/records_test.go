package tools

import (
	"testing"

	"github.com/alpine-hodler/sherpa/pkg/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

type testcase struct {
	Int int `json:"int"`
}

func TestAssignReadResponseRecords(t *testing.T) {
	t.Run("assign slice", func(t *testing.T) {
		rsp := new(proto.ReadResponse)

		r1, err := structpb.NewStruct(map[string]interface{}{"int": 1})
		require.NoError(t, err)

		rsp.Records = []*structpb.Struct{r1}

		tc := []*testcase{}
		err = AssignReadResponseRecords(rsp, &tc)
		require.NoError(t, err)

		expected := []*testcase{{1}}
		require.Equal(t, expected, tc)
	})

	assignSliceBenchmarkResponse := new(proto.ReadResponse)
	assignSliceBenchmarkExpected := []*testcase{}
	for i := 1; i <= 1e6; i++ {
		r, err := structpb.NewStruct(map[string]interface{}{"int": i})
		require.NoError(t, err)
		assignSliceBenchmarkExpected = append(assignSliceBenchmarkExpected, &testcase{Int: i})
		assignSliceBenchmarkResponse.Records = append(assignSliceBenchmarkResponse.Records, r)
	}
	t.Run("assign slice benchmark", func(t *testing.T) {
		tc := []*testcase{}
		err := AssignReadResponseRecords(assignSliceBenchmarkResponse, &tc)
		require.NoError(t, err)
		require.Equal(t, assignSliceBenchmarkExpected, tc)
	})
}
