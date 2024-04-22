package cookies

import (
	"errors"
	"net/http"
)

/*
	func main() {
		// Start a web server with the two endpoints.
		mux := http.NewServeMux()
		mux.HandleFunc("/set", setCookieHandler)
		mux.HandleFunc("/get", getCookieHandler)

		log.Print("Listening...")
		err := http.ListenAndServe(":3000", mux)
		if err != nil {
			log.Fatal(err)
		}
	}
*/
func SetCookieHandler(w http.ResponseWriter, r *http.Request, token string) (err error) {
	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	// Write a HTTP response as normal.
	return nil
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) (id string, err error) {
	// Retrieve the cookie from the request using its name (which in our case is
	// "exampleCookie"). If no matching cookie is found, this will return a
	// http.ErrNoCookie error. We check for this, and return a 400 Bad Request
	// response to the client.
	cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return "", http.ErrNoCookie
		default:
			//log.Println(err)
			return "", http.ErrAbortHandler
			//http.Error(w, "server error", http.StatusInternalServerError)
		}
	}

	return cookie.Value, nil
}
