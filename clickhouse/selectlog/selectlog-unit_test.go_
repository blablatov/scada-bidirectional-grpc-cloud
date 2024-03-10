// Модульное (unit) тестирование выборки лога.
// Определены тестовые сигнатуры методов пакета.
// Объявлены тестовые типы соответствующие интерфейсу, его экземпляры.
// go test -v selectlog-unit_test.go

package selectlog

import (
	"fmt"
	"reflect"
	"testing"
)

type Selecter interface {
	ReadFile(interface{}) (interface{}, error)
	NewCertPool() interface{}
	AppendCertsFromPEM(interface{}) bool
	NewRequest(interface{}, string, interface{}) (interface{}, error)

	LoadX509KeyPair(...interface{}) (interface{}, error)
	Get(string) (interface{}, error)
	Set(...string)
	Do(interface{}) (interface{}, error)
	ReadAll(interface{}) (interface{}, error)
}

type TestSelect struct {
	Selecter
	success bool
}

func (ts *TestSelect) ReadFile(interface{}) (interface{}, error) {
	if ts.success {
		return nil, nil
	}
	return nil, fmt.Errorf("ReadFile test error")
}

func TestReadFile(t *testing.T) {
	type args struct {
		sel Selecter
	}
	tests := []struct {
		file       string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			file: "file exists",
			args: args{
				sel: &TestSelect{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			file: "file not exists",
			args: args{
				sel: &TestSelect{success: false},
			},
			wantErr:    fmt.Errorf("file test error"),
			wantExists: false,
		},
	}
	var f TestSelect
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			gotExists, gotErr := f.ReadFile(tt.args.sel)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func ReadFile() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func ReadFile() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestSelect) NewCertPool() interface{} {
	if tс.success {
		return nil
	}
	return fmt.Errorf("NewRequest test error")
}

func TestNewCertPool(t *testing.T) {
	type args struct {
		cert Selecter
	}
	tests := []struct {
		cert    string
		args    args
		wantErr error
	}{
		{
			cert: "cert exists",
			args: args{
				cert: &TestSelect{success: true},
			},
			wantErr: nil,
		}, {
			cert: "cert not exists",
			args: args{
				cert: &TestSelect{success: false},
			},
			wantErr: fmt.Errorf("cert test error"),
		},
	}

	var f TestSelect

	for _, tt := range tests {
		t.Run(tt.cert, func(t *testing.T) {
			gotErr := f.NewCertPool()
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func NewCertPool() gotExists = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestSelect) AppendCertsFromPEM(interface{}) bool {
	if tс.success {
		return true
	}
	return false
}

func TestAppendCertsFromPEM(t *testing.T) {
	type args struct {
		acert Selecter
	}
	tests := []struct {
		acert      string
		args       args
		wantExists bool
	}{
		{
			acert: "acert exists",
			args: args{
				acert: &TestSelect{success: true},
			},
			wantExists: true,
		}, {
			acert: "acert not exists",
			args: args{
				acert: &TestSelect{success: false},
			},
			wantExists: false,
		},
	}

	var f TestSelect

	for _, tt := range tests {
		t.Run(tt.acert, func(t *testing.T) {
			gotExists := f.NewCertPool()
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func AppendCertsFromPEM() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func (tс *TestSelect) NewRequest(interface{}, string, interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("NewRequest test error")
}

func TestNewRequest(t *testing.T) {
	type args struct {
		req Selecter
	}
	tests := []struct {
		apiUrl     string
		method     string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			apiUrl: "apiUrl exists",
			args: args{
				req: &TestSelect{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			apiUrl: "apiUrl not exists",
			args: args{
				req: &TestSelect{success: false},
			},
			wantErr:    fmt.Errorf("apiUrl test error"),
			wantExists: false,
		}, {
			method: "method exists",
			args: args{
				req: &TestSelect{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			method: "method not exists",
			args: args{
				req: &TestSelect{success: false},
			},
			wantErr:    fmt.Errorf("method test error"),
			wantExists: false,
		},
	}
	var f TestSelect
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			gotExists, gotErr := f.NewRequest("GET", "apiUrl", nil)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func NewRequest() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func NewRequest() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestSelect) Do(interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("Do test error")
}

func TestDo(t *testing.T) {
	type args struct {
		do Selecter
	}
	tests := []struct {
		request    string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			request: "do exists",
			args: args{
				do: &TestSelect{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			request: "do not exists",
			args: args{
				do: &TestSelect{success: false},
			},
			wantErr:    fmt.Errorf("do test error"),
			wantExists: false,
		},
	}
	var f TestSelect
	for _, tt := range tests {
		t.Run(tt.request, func(t *testing.T) {
			gotExists, gotErr := f.Do(tt.args.do)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func Do() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func Do() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (tс *TestSelect) ReadAll(interface{}) (interface{}, error) {
	if tс.success {
		return nil, nil
	}
	return nil, fmt.Errorf("ReadAll test error")
}

func TestReadAll(t *testing.T) {
	type args struct {
		rall Selecter
	}
	tests := []struct {
		read       string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			read: "read exists",
			args: args{
				rall: &TestSelect{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			read: "read not exists",
			args: args{
				rall: &TestSelect{success: false},
			},
			wantErr:    fmt.Errorf("read do test error"),
			wantExists: false,
		},
	}

	var f TestSelect

	for _, tt := range tests {
		t.Run(tt.read, func(t *testing.T) {
			gotExists, gotErr := f.ReadAll(tt.args.rall)
			if reflect.DeepEqual(gotExists, tt.wantExists) {
				t.Errorf("Check func ReadAll() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func ReadAll() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
