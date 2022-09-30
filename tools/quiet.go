// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
package tools

import (
	"log"
	"os"
)

// Quiet will suppress the output during the test I use the following code. I fixes output as well as logging. After
// test is done it resets the output streams.
func Quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)

	return func() {
		defer null.Close()

		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
