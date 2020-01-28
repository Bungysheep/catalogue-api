package routes

import (
	authcontrollerv1 "github.com/bungysheep/catalogue-api/pkg/controllers/v1/authcontroller"
	cataloguecontrollerv1 "github.com/bungysheep/catalogue-api/pkg/controllers/v1/cataloguecontroller"
	"github.com/bungysheep/catalogue-api/pkg/protocols/rest/middlewares"
	"github.com/gorilla/mux"
)

// APIV1RouteHandler builds Api v1 routes
func APIV1RouteHandler() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middlewares.DefaultMiddleware)

	v1Router := router.PathPrefix("/v1").Subrouter()

	authController := authcontrollerv1.NewAuthController()
	v1Router.HandleFunc("/users", authController.GetAll).Methods("GET")

	authRouter := v1Router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signin", authController.SignIn).Methods("POST")
	authRouter.HandleFunc("/register", authController.Register).Methods("POST")

	catalogueController := cataloguecontrollerv1.NewCatalogueController()
	clgRouter := v1Router.PathPrefix("").Subrouter()
	clgRouter.Use(middlewares.AuthenticationMiddleware)
	clgRouter.HandleFunc("/catalogues", catalogueController.GetAll).Methods("GET")
	clgRouter.HandleFunc("/catalogues/{id}", catalogueController.GetByID).Methods("GET")
	clgRouter.HandleFunc("/catalogues", catalogueController.Create).Methods("POST")
	clgRouter.HandleFunc("/catalogues/{id}", catalogueController.Update).Methods("PUT")
	clgRouter.HandleFunc("/catalogues/{id}", catalogueController.Delete).Methods("DELETE")

	return router
}
