package internal

import (
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
