package errs

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestErrAsUtil(t *testing.T) {
	type args struct {
		i interface{}
	}

	validUtil := &Util{
		StatusCode: http.StatusInternalServerError,
		List: []error{
			errors.New("mock err 1"),
		},
	}

	tests := []struct {
		name  string
		args  args
		want  *Util
		want1 bool
	}{
		{
			name:  "valid",
			args:  args{i: error(validUtil)},
			want:  validUtil,
			want1: true,
		},
		{
			name:  "invalid",
			args:  args{i: true},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ErrAsUtil(tt.args.i)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrAsUtil() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ErrAsUtil() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
