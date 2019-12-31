package event

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/tcernestw-go/moment"
	"runtime"
	"strconv"
	"strings"
)

// Exception contains details for recover
type Exception struct {
	ID         string      `json:"id" bson:"id"`
	Err        string      `json:"err" bson:"err"`
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
	if err != nil {
		exc.Err = err.Error()
	}
	return
}

// HasErr checks if the exception contains any error
func (exc Exception) HasErr() (has bool) {
	return exc.Err != ""
}

// Error returns the err in "error"
func (exc Exception) Error() (err error) {
	return errors.New(exc.Err)
}

// Log logs the details of exception
func (exc Exception) Log() (log string) {
	format := moment.NewFormat().YYYY("-").MM("-").DD(" ").HH24(":").Min(":").SS().Str
	display := "Exception (" + exc.ID + ") (" + moment.Now().Display(format) + ")\n"
	display += " - Time: " + moment.NewUnix(exc.Time, moment.MilliSec).Display(format) + "\n"
	display += " - Remark: " + exc.Remark + "\n"
	display += " - Trace:\n"
	display += "    " + strings.Join(LogTraces(exc.Traces), "\n    ")
	display += " - Error: " + exc.Err + "\n"
	display += " - Attachment: "
	return fmt.Sprint(display, exc.Attachment)
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
func LogTraces(traces []Trace) (logs []string) {
	var length = len(traces)
	for i := 0; i < length; i++ {
		trace := traces[i]
		logs = append(logs, strconv.Itoa(i+1)+". pc = "+fmt.Sprint(trace.Pc)+", file = "+trace.File+", line = "+strconv.Itoa(trace.Line)+"\n")
	}
	return
}

// Recover get the panic content as error, and create a exception
func Recover(exc *Exception, msg string, attachment interface{}) {
	if r := recover(); r != nil && exc != nil {
		*exc = NewException(msg, errors.New(fmt.Sprint(r)), attachment)
	}
}
