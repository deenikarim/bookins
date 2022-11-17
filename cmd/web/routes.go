package main

import (
	"github.com/deenikarim/bookings/pkg/config"
	"github.com/deenikarim/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	/**/ //installing middleware
	mux.Use(middleware.Recoverer)

	/**/ //how to use the middleware created
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	/**/ //here is where to set up the routes

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/reservation", handlers.Repo.Reservation)

	//fileServer: creating fileServer to enable static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	/**/
	return mux
}
