package main

import (
	"github.com/deenikarim/bookings/internal/helpers"
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
}

//Auth to help protect our routes in the route.go file to ensure that only people who are logged in
// actually have access to the routes we want to protect from unauthenticated users
func Auth(next http.Handler) http.Handler {
	//TODO ; one thing that is different from this middleware is that we actually need to have access
	// to the request because we are going to call that helper method call Auth and it requires a pointer
	// to the to the request as an input parameter

	//HandlerFunc is an adapter that allows the use of ordinary functions as http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if a user is not authenticated then I want to do something
		if !helpers.IsAuthenticated(r) {
			app.Session.Put(r.Context(), "error", "login first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		//if the above code doesn't fail then we just pass it to the next middleware if any and
		// the request lifecycle continue in its process
		next.ServeHTTP(w, r)
	})
}
