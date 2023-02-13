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
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom) //chi router
	// handles the connecting of book now URL to make-reservation page(should have the same URL path as the one defined in General page)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	//route to handle the POST request method
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.CheckAvailabilityJSON)

	//routes for displaying login page
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.LogOut)

	//fileServer: creating fileServer to enable static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	//here where we want to protect our routes
	//todo Route() creates a new mux with fresh middleware stack and mounts it along the pattern
	//  as a subRouter
	//  WE SETUP A ROUTE TO THE ADMIN-DASHBOARD THAT WILL ONLY BE AVAILABLE TO USERS WHO HAVE LOGGED IN
	mux.Route("/admin", func(mux chi.Router) {
		//what we need to do here is

		// add a middleware
		//mux.Use(Auth) //make sure a user is authenticated

		//now let's create a new route
		//The pattern is going to have (/admin) pre-attach to it so we want to add "dashboard" which will be
		//(/admin/dashboard) and is also going to have a handler
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		//	THE DASHBOARD TEMPLATE TO ENABLE FUNCTIONALITY
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
		//mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)
		//.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		//mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		//mux.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		//mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
	})

	/**/
	return mux
}
