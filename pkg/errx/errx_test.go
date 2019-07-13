package errx

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestAPIError_IsError(t *testing.T) {
	type fields struct {
		Err        error
		StatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should return true on statusCode 500",
			fields: fields{
				Err:        errors.New("test error"),
				StatusCode: http.StatusInternalServerError,
			},
			want: true,
		},
		{
			name: "should return true on statusCode 400",
			fields: fields{
				Err:        errors.New("test error"),
				StatusCode: http.StatusBadRequest,
			},
			want: true,
		},
		{
			name: "should return false on statusCode 200",
			fields: fields{
				Err:        errors.New("test error"),
				StatusCode: http.StatusOK,
			},
			want: false,
		},
		{
			name: "should return false on statusCode 300",
			fields: fields{
				Err:        errors.New("test error"),
				StatusCode: http.StatusPermanentRedirect,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &APIError{
				Err:        tt.fields.Err,
				StatusCode: tt.fields.StatusCode,
			}
			if got := e.IsError(); got != tt.want {
				t.Errorf("APIError.IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	type fields struct {
		Err        error
		StatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return message: xxx",
			fields: fields{
				Err:        errors.New("xxx"),
				StatusCode: http.StatusBadRequest,
			},
			want: "xxx",
		},
		{
			name: "should return message: (empty string)",
			fields: fields{
				Err:        errors.New(""),
				StatusCode: http.StatusOK,
			},
			want: "",
		},
		{
			name: "should return message: Bad request",
			fields: fields{
				Err:        errors.New("Bad request"),
				StatusCode: http.StatusBadRequest,
			},
			want: "Bad request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &APIError{
				Err:        tt.fields.Err,
				StatusCode: tt.fields.StatusCode,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Serialize(t *testing.T) {
	type fields struct {
		Err        error
		StatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return string with error message and status code",
			fields: fields{
				Err:        errors.New("Invalid credentials"),
				StatusCode: http.StatusBadRequest,
			},
			want: `{"message":"Invalid credentials","statusCode":400}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &APIError{
				Err:        tt.fields.Err,
				StatusCode: tt.fields.StatusCode,
			}
			if got := e.Serialize(); got != tt.want {
				t.Errorf("APIError.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPIError(t *testing.T) {
	type args struct {
		err        error
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want APIError
	}{
		{
			name: "should create new APIError with statusCode and error",
			args: args{
				err:        errors.New("xxx"),
				statusCode: http.StatusBadRequest,
			},
			want: APIError{
				Err:        errors.New("xxx"),
				StatusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIError(tt.args.err, tt.args.statusCode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIError() = %v, want %v", got, tt.want)
			}
		})
	}
}
