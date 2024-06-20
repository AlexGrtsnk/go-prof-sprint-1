package internal

import (
	"net/http"
	"testing"
)

func TestNewProducer(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	_, err := GzipFormatHandlerJSON(nil, request)
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
