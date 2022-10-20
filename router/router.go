// Package router defines request urls of microservice api
package router

import (
	"github.com/labstack/echo/v4"
	middleKit "github.com/oboadagd/kit-go/middleware/echo"
	"github.com/oboadagd/location-mgmt/controller"
)

// Router represents the router layer.
type Router struct {
	server             *echo.Echo
	locationController controller.LocationControllerInterface
	errorMiddleware    middleKit.ErrorHandlerMiddlewareInterface
}

// NewRouter initializes router layer
func NewRouter(
	server *echo.Echo,
	locationController controller.LocationControllerInterface,
	errorMiddleware middleKit.ErrorHandlerMiddlewareInterface,
) *Router {
	return &Router{
		server,
		locationController,
		errorMiddleware,
	}
}

// Init implements request urls definition
func (r *Router) Init() {
	//create a default router with default middlewares
	basePath := r.server.Group("/location-mgmt")

	locations := basePath.Group("/locations", r.errorMiddleware.HandlerError)
	{
		locations.POST("", r.locationController.Save)
		locations.GET("/users/:latitude/:longitude/:radius", r.locationController.GetUsersByLocationAndRadius)
	}
}
