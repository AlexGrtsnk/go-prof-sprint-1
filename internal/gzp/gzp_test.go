package internal

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestGzipFormatHandlerJSON(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	_, err := GzipFormatHandlerJSON(nil, request)
	if err != nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestGzipFormatHandlerJSONBadreq(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	request.Header.Set(`Content-Encoding`, `gzip`)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_, err := GzipFormatHandlerJSON(nil, request)
	if err != nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestGzipFormatHandlerJSONBadreader(t *testing.T) {
	b := new(bytes.Buffer)
	_, _ = io.WriteString(b, "http://localhost:8080/multi")
	request, _ := http.NewRequest("POST", "/", b)
	request.Header.Set(`Content-Encoding`, `gzip`)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_, err := GzipFormatHandlerJSON(nil, request)
	if err != nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestGzipWrite(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	var v gzipWriter
	_, _ = v.Write([]byte("string"))

}

func TestGzipHandle(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()
	GzipHandle(nil)

}
