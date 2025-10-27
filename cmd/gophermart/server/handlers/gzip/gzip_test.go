package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCompressReader(t *testing.T) {

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"test 1", "*gzip.compressReader", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err := zb.Write([]byte(`{"id":"met", "type":"counter"}`))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)

			r := httptest.NewRequest("POST", "/value", buf)

			got, err := NewCompressReader(r.Body)
			assert.NoError(t, err, "error in creation")
			assert.Equalf(t, tt.want, reflect.ValueOf(got).Type().String(), "NewCompressReader()")
		})
	}
}

func TestNewCompressWriter(t *testing.T) {

	tests := []struct {
		name string
		want string
	}{
		{"test 1", "*gzip.compressWriter"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCompressWriter(httptest.NewRecorder())
			assert.Equalf(t, tt.want, reflect.ValueOf(got).Type().String(), "NewCompressWriter()")
		})
	}
}

func Test_compressReader_Close(t *testing.T) {
	type fields struct {
		r  io.ReadCloser
		zr *gzip.Reader
	}
	var f fields

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"test 1", f, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err := zb.Write([]byte(`{"id":"met", "type":"counter"}`))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)

			r := httptest.NewRequest("POST", "/value", buf)

			got, _ := NewCompressReader(r.Body)

			assert.NoError(t, got.Close(), "error in close")
		})
	}
}

func Test_compressReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name  string
		args  args
		wantN int
	}{
		{"test 1", args{p: []byte(`{"id":"met", "type":"counter"}`)}, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err := zb.Write([]byte(`{"id":"met", "type":"counter"}`))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)

			r := httptest.NewRequest("POST", "/value", buf)
			got, err := NewCompressReader(r.Body)
			require.NoError(t, err, "error in creation")

			gotN, _ := got.Read(tt.args.p)
			assert.Equalf(t, tt.wantN, gotN, "error in byte count")
		})
	}
}

func Test_compressWriter_Close(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	var f fields
	f.zw = gzip.NewWriter(bytes.NewBuffer(nil))
	tests := []struct {
		name   string
		fields fields
	}{
		{"test 1", f},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}
			assert.NoError(t, c.Close(), "error in close")
		})
	}
}

func Test_compressWriter_Header(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	f := fields{
		w:  httptest.NewRecorder(),
		zw: gzip.NewWriter(bytes.NewBuffer(nil)),
	}

	tests := []struct {
		name   string
		fields fields
		want   http.Header
	}{
		{"tast 1", f, http.Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}

			assert.Equalf(t, tt.want, c.Header(), "error in header")
		})
	}
}

func Test_compressWriter_Write(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	f := fields{
		w:  httptest.NewRecorder(),
		zw: gzip.NewWriter(bytes.NewBuffer(nil)),
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{"test 1", f, args{p: []byte(`{"id":"met", "type":"counter"}`)}, 30, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}
			got, err := c.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressWriter_WriteHeader(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	f := fields{
		w:  httptest.NewRecorder(),
		zw: gzip.NewWriter(bytes.NewBuffer(nil)),
	}
	type args struct {
		statusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test 1", f, args{statusCode: 200}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}
			c.WriteHeader(tt.args.statusCode)
		})
	}
}
