package main

import (
	"github.com/djedjethai/vigilate/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func routes() http.Handler {

	mux := chi.NewRouter()

	// default middleware
	mux.Use(SessionLoad)
	mux.Use(RecoverPanic)
	mux.Use(NoSurf)
	mux.Use(CheckRemember)

	// login
	mux.Get("/", handlers.Repo.LoginScreen)
	mux.Post("/", handlers.Repo.Login)

	mux.Get("/user/logout", handlers.Repo.Logout)

	// mux.Get("/pusher-test", handlers.Repo.TestPusher)

	// routes to our pusher
	mux.Route("/pusher", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Post("/auth", handlers.Repo.PusherAuth)
	})

	// admin routes
	mux.Route("/admin", func(mux chi.Router) {
		// all admin routes are protected
		mux.Use(Auth)

		// overview
		mux.Get("/overview", handlers.Repo.AdminDashboard)

		// events
		mux.Get("/events", handlers.Repo.Events)

		// settings
		mux.Get("/settings", handlers.Repo.Settings)
		mux.Post("/settings", handlers.Repo.PostSettings)

		// service status pages (all hosts)
		mux.Get("/all-healthy", handlers.Repo.AllHealthyServices)
		mux.Get("/all-warning", handlers.Repo.AllWarningServices)
		mux.Get("/all-problems", handlers.Repo.AllProblemServices)
		mux.Get("/all-pending", handlers.Repo.AllPendingServices)

		// users
		mux.Get("/users", handlers.Repo.AllUsers)
		mux.Get("/user/{id}", handlers.Repo.OneUser)
		mux.Post("/user/{id}", handlers.Repo.PostOneUser)
		mux.Get("/user/delete/{id}", handlers.Repo.DeleteUser)

		// schedule
		mux.Get("/schedule", handlers.Repo.ListEntries)

		// preferences
		mux.Post("/preference/ajax/set-system-pref", handlers.Repo.SetSystemPref)

		// hosts
		mux.Get("/host/all", handlers.Repo.AllHosts)
		mux.Get("/host/{id}", handlers.Repo.Host)
		mux.Post("/host/{id}", handlers.Repo.PostHost)
		mux.Post("/host/ajax/toggle-service", handlers.Repo.TogglerServiceForHost)
		mux.Post("/preference/ajax/toggle-monitoring", handlers.Repo.ToggleMonitoring)
		mux.Get("/perform-check/{id}/{oldstatus}", handlers.Repo.TestCheck)
	})

	// static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
