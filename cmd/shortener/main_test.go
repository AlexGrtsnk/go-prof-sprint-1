package main

import (
	fun "go-prof-sprint-1/cmd/functions"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			fun.MainPage(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestAPIPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			fun.APIPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestJSONPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			fun.APIPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}
