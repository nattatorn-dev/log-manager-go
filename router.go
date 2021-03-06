package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nattatorn-dev/log-manager/common"
	"github.com/nattatorn-dev/log-manager/controllers"
	"github.com/nattatorn-dev/log-manager/models"
	"github.com/urfave/negroni"
)

func InitRoutes(store models.DataStorer, authModule common.Authorizer) http.Handler {
	router := mux.NewRouter().StrictSlash(false)

	/* healthcheck / motd */
	router.Handle("/", http.HandlerFunc(controllers.Healthcheck)).Methods("GET")

	/* User routes */
	router.Handle("/users", controllers.WithDbAndAuth(authModule, store, controllers.Register)).Methods("POST")
	router.Handle("/users/login", controllers.WithDbAndAuth(authModule, store, controllers.Login)).Methods("POST")
	//router.Handle("/users", withDb(store, controllers.GetUser)).Methods("GET")

	/* Task routes  */
	taskRouter := mux.NewRouter().StrictSlash(false)
	taskRouter.Handle("/tasks", controllers.WithDb(store, controllers.GetAllTasks)).Methods("GET")
	taskRouter.Handle("/tasks/{id}", controllers.WithDb(store, controllers.GetTask)).Methods("GET")
	taskRouter.Handle("/tasks/{id}", controllers.WithDb(store, controllers.DeleteTask)).Methods("DELETE")
	taskRouter.Handle("/tasks", controllers.WithDb(store, controllers.CreateTask)).Methods("POST")
	taskRouter.Handle("/tasks/{id}", controllers.WithDb(store, controllers.UpdateTask)).Methods("PUT")

	/* middleware */
	commonMidleware := negroni.New(
		negroni.NewLogger(),
	)

	router.PathPrefix("/tasks").Handler(negroni.New(
		common.WithAuth(authModule),
		negroni.Wrap(taskRouter),
	))

	// common wraps all routes in default middleware
	// this includes all API hits
	commonMidleware.UseHandler(router)

	return commonMidleware
}
