package cookies

import (
	"net/http"
	"testing"
)

func TestGetCookieHandler(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	_, err := GetCookieHandler(nil, request)
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}
