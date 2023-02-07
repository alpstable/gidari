// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	"context"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

type DecodeType int32

const (
	DecodeTypeUnknown DecodeType = iota
	DecodeTypeJSON
)

type ListWriter interface {
	Write(cxt context.Context, list *structpb.ListValue) error
}
