package main

import (
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/handlers"
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
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)

	//route to handle the POST request method
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.CheckAvailabilityJSON)

	//fileServer: creating fileServer to enable static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	/**/
	return mux
}
