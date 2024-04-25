package cookies

import (
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func Write(w http.ResponseWriter, cookie http.Cookie) error {
	// Encode the cookie value using base64.
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// Check the total length of the cookie contents. Return the ErrValueTooLong
	// error if it's more than 4096 bytes.
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	// Write the cookie as normal.
	http.SetCookie(w, &cookie)

	return nil
}

func Read(r *http.Request, name string) (string, error) {
	// Read the cookie as normal.
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded cookie value. If the cookie didn't contain a
	// valid base64-encoded value, this operation will fail and we return an
	// ErrInvalidValue error.
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	// Return the decoded cookie value.
	return string(value), nil
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request, tknm string) (cks *http.Cookie) {
	// Initialize the cookie as normal.
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    tknm,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	// Write the cookie. If there is an error (due to an encoding failure or it
	// being too long) then log the error and send a 500 Internal Server Error
	// response.
	http.SetCookie(w, &cookie)
	/*
		err = Write(w, cookie)
		if err != nil {
			return err
		}
	*/
	return &cookie
	//w.Write([]byte("cookie set!"))
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) (token string, err error) {
	// Use the Read() function to retrieve the cookie value, additionally
	// checking for the ErrInvalidValue error and handling it as necessary.
	cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		return "", err
	}
	//value = cookie.Value
	//fmt.Println("valuie of cookie inside ", value)
	return cookie.Value, nil
}
