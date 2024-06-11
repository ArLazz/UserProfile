// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"userprofile/database"
	"userprofile/handlers"
	"userprofile/models"
	"userprofile/restapi/operations"
)

//go:generate swagger generate server --target ../../goproject --name Userprofile --spec ../swagger.yaml --principal interface{}

func configureFlags(api *operations.UserprofileAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.UserprofileAPI) http.Handler {

	api.ServeError = errors.ServeError
	api.Logger = log.Printf
	api.UseSwaggerUI()
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	handler := handlers.NewHandler(database.NewInMemoryDB())
	api.BasicAuthAuth = func(username string, pass string) (interface{}, error) {
		user, err := handler.AuthenticateUser(username, pass)
		if err != nil {
			err = errors.New(401, err.Error())
			return nil, err
		}
		
		return user, err
	}

	api.DeleteUserIDHandler = operations.DeleteUserIDHandlerFunc(func(params operations.DeleteUserIDParams, principal interface{}) middleware.Responder {
		user := principal.(*models.User)
		if !user.Admin {
			api.Logger("user %s attempted to delete user but is not an admin", user.Username)
			return operations.NewDeleteUserIDForbidden().WithPayload(&models.ErrorResponse{Message: "you are not admin"})
		}

		err := handler.DeleteUserByID(params)
		if err != nil {
			api.Logger("failed to delete user: %s", err.Error())
			return operations.NewDeleteUserIDNotFound().WithPayload(&models.ErrorResponse{Message: err.Error()})
		}
		api.Logger("deleted user with id = %s", params.ID)

		return operations.NewDeleteUserIDNoContent()
	})

	api.GetUserHandler = operations.GetUserHandlerFunc(func(params operations.GetUserParams, principal interface{}) middleware.Responder {
		payload := handler.GetUser(params)
		api.Logger("get users list")
		return operations.NewGetUserOK().WithPayload(payload)
	})

	api.GetUserIDHandler = operations.GetUserIDHandlerFunc(func(params operations.GetUserIDParams, principal interface{}) middleware.Responder {
		payload, err := handler.GetUserByID(params)
		if err != nil {
			api.Logger("failed to get user by id(%s): %s", params.ID, err.Error())
			return operations.NewGetUserIDNotFound().WithPayload(&models.ErrorResponse{Message: err.Error()})
		}
		api.Logger("get user with id = %s", params.ID)
		return operations.NewGetUserIDOK().WithPayload(payload)
	})

	api.PostUserHandler = operations.PostUserHandlerFunc(func(params operations.PostUserParams, principal interface{}) middleware.Responder {
		user := principal.(*models.User)
		if !user.Admin {
			api.Logger("user %s attempted to create user but is not an admin", user.Username)
			return operations.NewDeleteUserIDForbidden().WithPayload(&models.ErrorResponse{Message: "you are not admin"})
		}

		err := handler.AddUser(params)
		if err != nil {
			api.Logger("failed to create user: %s", err.Error())
			return operations.NewPostUserBadRequest().WithPayload(&models.ErrorResponse{Message: err.Error()})
		}

		api.Logger("created user with name: %s", params.User.Username)
		return operations.NewPostUserCreated()
	})

	api.PutUserIDHandler = operations.PutUserIDHandlerFunc(func(params operations.PutUserIDParams, principal interface{}) middleware.Responder {
		user := principal.(*models.User)
		if !user.Admin {
			api.Logger("user %s attempted to update user but is not an admin", user.Username)
			return operations.NewDeleteUserIDForbidden().WithPayload(&models.ErrorResponse{Message: "you are not admin"})
		}

		err := handler.UpdateUserByID(params)
		if err != nil {
			api.Logger("failed to update user: %s", err.Error())
			return operations.NewPutUserIDNotFound().WithPayload(&models.ErrorResponse{Message: err.Error()})
		}

		api.Logger("user updated")
		return operations.NewPutUserIDOK()
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
