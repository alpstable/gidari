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
