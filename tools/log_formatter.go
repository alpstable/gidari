package tools

import (
	"fmt"
	"strings"
	"time"
)

// LogFormatter encapsulates data that is used to format a log message.
type LogFormatter struct {
	WorkerID      int
	Duration      time.Duration
	Msg           string
	UpsertedCount int32
	MatchedCount  int32
}

const (
	// LogFormatterWorkerID the label of the worker id.
	LogFormatterWorkerID = "w"

	// LogFormatterDuration the label of the duration.
	LogFormatterDuration = "d"

	// LogFormatterMsg the label of the message.
	LogFormatterMsg = "m"

	// LogFormatterUpsertedCount the label of the upserted count.
	LogFormatterUpsertedCount = "u"

	// LogFormmaterMatchedCount the label of the matched count.
	LogFormatterMatchedCount = "c"
)

// String uses the data from the LogFormatter object to build a log message.
func (lf LogFormatter) String() string {
	var sb strings.Builder
	if lf.WorkerID > 0 {
		sb.WriteString(fmt.Sprintf("%s:%d, ", LogFormatterWorkerID, lf.WorkerID))
	}
	if lf.Duration > 0 {
		sb.WriteString(fmt.Sprintf("%s:%s, ", LogFormatterDuration, lf.Duration))
	}
	if lf.UpsertedCount > 0 {
		sb.WriteString(fmt.Sprintf("%s:%d, ", LogFormatterUpsertedCount, lf.UpsertedCount))
	}
	if lf.MatchedCount > 0 {
		sb.WriteString(fmt.Sprintf("%s:%d, ", LogFormatterMatchedCount, lf.MatchedCount))
	}
	if lf.Msg != "" {
		sb.WriteString(fmt.Sprintf("%s:%s, ", LogFormatterMsg, lf.Msg))
	}
	return fmt.Sprintf("{%s}", strings.TrimSuffix(sb.String(), ", "))
}
