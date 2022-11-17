package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

/* writeToConsole illustrating the format of how to write middleware
func writeToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//inside the return statement
		fmt.Println("hit the page")

		//need at the end of the code above, moves to the next which could be the next middleware
		next.ServeHTTP(w, r)
	})
}
*/

//NoSurf create CSRFToken: adds CSRF protection to all POST Request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next) //will create a new handler

	//need to set some values for it
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", //the way to refer to entire site for a cookie
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad create middleware that loads  and saves the sessions on every request to make web server state aware
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
	//return session.LoadAndSave(next)
}
