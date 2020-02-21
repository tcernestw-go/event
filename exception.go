package event

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tcernestw-go/moment"
)

// Exception contains details for recover
type Exception struct {
	ID         string      `json:"id" bson:"id"`
	Err        error       `json:"err" bson:"err"`
	Remark     string      `json:"remark" bson:"remark"`
	Attachment interface{} `json:"attachment" bson:"attachment"`
	Time       int64       `json:"time" bson:"time"`
	Traces     []Trace     `json:"traces" bson:"traces"`
}

// Trace contains details of programming line
type Trace struct {
	Pc   uintptr
	File string
	Line int
}

// NewException returns exception with specific error details
// If error is nil, the exception
func NewException(remark string, err error, attachment interface{}) (exc Exception) {
	exc = Exception{
		ID:         uuid.NewV4().String(),
		Remark:     remark,
		Attachment: attachment,
		Time:       moment.Now().UnixMillis(),
		Traces:     GetTraces(1),
	}
	return
}

// HasErr checks if the exception contains any error
func (exc Exception) HasErr() (has bool) {
	return exc.Err != nil
}

// Error returns the err in "error"
func (exc Exception) Error() (err string) {
	format := "2006-01-02 15:04:05"
	return fmt.Sprintf(
		`
Exception(%v)
 - Time: %v
 - Remark: %v
 - Trace: %v
 - Error: %v
 - Attachment: %v
`,
		exc.ID,
		time.Unix(0, exc.Time*time.Hour.Milliseconds()).Format(format),
		exc.Remark,
		LogTraces(exc.Traces),
		exc.Err,
		exc.Attachment,
	)
}

// GetTraces get the programming lines calling this function
func GetTraces(layer int) (traces []Trace) {
	for skip := 1; ; skip++ {
		var ok bool
		var trace Trace
		trace.Pc, trace.File, trace.Line, ok = runtime.Caller(layer + skip)
		if !ok {
			break
		}
		traces = append(traces, trace)
	}
	return
}

// LogTraces log the traces
func LogTraces(traces []Trace) (log string) {
	var length = len(traces)
	for i := 0; i < length; i++ {
		trace := traces[i]
		log += "\n	" + strconv.Itoa(i+1) + ". pc = " + fmt.Sprint(trace.Pc) + ", file = " + trace.File + ", line = " + strconv.Itoa(trace.Line)
	}
	return
}

// Recover get the panic content as error, and create a exception
func Recover(exc *Exception, msg string, attachment interface{}) {
	if r := recover(); r != nil && exc != nil {
		*exc = NewException(msg, errors.New(fmt.Sprint(r)), attachment)
	}
}
