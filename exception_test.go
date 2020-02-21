package event

import (
	"errors"
	"testing"
	"time"
)

func TestException_Error(t *testing.T) {
	tests := []struct {
		name    string
		exc     Exception
		wantErr string
	}{
		{
			name: "Normal Case",
			exc: Exception{
				ID:         "normal_case_id",
				Err:        errors.New("normal_case_err"),
				Remark:     "normal_case_remark",
				Attachment: "normal_case_attachment",
				Time:       time.Unix(0, 1294969890000*int64(time.Millisecond)).UnixNano() / time.Hour.Milliseconds(),
				Traces: []Trace{
					Trace{
						Pc:   1,
						File: "testing1",
						Line: 24,
					},
					Trace{
						Pc:   2,
						File: "testing2",
						Line: 16,
					},
				},
			},
			wantErr: `
Exception(normal_case_id)
 - Time: 2011-01-14 09:51:29
 - Remark: normal_case_remark
 - Trace: 
	1. pc = 1, file = testing1, line = 24
	2. pc = 2, file = testing2, line = 16
 - Error: normal_case_err
 - Attachment: normal_case_attachment
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := tt.exc.Error(); gotErr != tt.wantErr {
				t.Errorf("Exception.Error() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
