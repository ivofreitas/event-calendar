package v1

import (
	"blankfactor/event-calendar/api/middleware"
	"blankfactor/event-calendar/api/swagger"
	"blankfactor/event-calendar/api/v1/event"
	"blankfactor/event-calendar/api/v1/health"
	"blankfactor/event-calendar/config"
	"blankfactor/event-calendar/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Register(g *echo.Group) {

	doc := config.GetEnv().Doc
	swagger.Register(swagger.Options{
		Title:       doc.Title,
		Description: doc.Description,
		Version:     doc.Version,
		BasePath:    doc.BasePath,
		Group:       g.Group("/swagger"),
	})

	g.GET("/health", health.Handle)

	handler := event.NewHandler()
	ctrl := middleware.NewController(handler.OverlappingEvents, http.StatusCreated, new(model.OverlappingRequest))
	event := g.Group("/event")
	event.POST("", ctrl.Handle)
}
