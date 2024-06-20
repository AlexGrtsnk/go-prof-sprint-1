package internal

import (
	"net/http"
	"testing"
)

func TestGzipFormatHandlerJSON(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	_, err := GzipFormatHandlerJSON(nil, request)
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestGzipWrite(t *testing.T) {
	//request, _ := http.NewRequest("POST", "/", nil)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	var v gzipWriter
	_, _ = v.Write([]byte("string"))
	//_, err := GzipFormatHandlerJSON(nil, request)

}

func TestGzipHandle(t *testing.T) {
	//request, _ := http.NewRequest("POST", "/", nil)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()
	GzipHandle(nil)
	//_, err := GzipFormatHandlerJSON(nil, request)

}
