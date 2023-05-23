package api

import (
	"blankfactor/event-calendar/api/middleware"
	"blankfactor/event-calendar/api/v1"
	"blankfactor/event-calendar/config"
	"blankfactor/event-calendar/log"
	"blankfactor/event-calendar/model"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	http   *echo.Echo
	logger *logrus.Entry
	signal chan struct{}
}

func NewServer() *Server {
	log.Init()

	return &Server{
		logger: log.NewEntry(),
		signal: make(chan struct{}),
	}
}

func (s *Server) Run() {
	s.start()
	s.logger.Println("Server started and waiting for the graceful signal...")
	<-s.signal
}

func (s *Server) start() {
	go s.watchStop()

	serverConfig := config.GetEnv().Server

	s.http = echo.New()
	s.logger.Infof("Server is starting in port %s.", serverConfig.Port)

	s.http.Validator = middleware.NewValidator()
	s.http.Binder = middleware.NewBinder()
	s.http.Use(middleware.Logger())
	s.http.Use(echomiddleware.Recover())
	s.http.Pre(echomiddleware.RemoveTrailingSlash())
	s.http.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		httpLog := log.Get(c.Request().Context(), log.HTTPKey).(*log.HTTP)
		httpLog.Level = logrus.WarnLevel

		responseErr := model.ErrorDiscover(err)
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(responseErr.StatusCode)
		} else {
			err = c.JSON(responseErr.StatusCode, responseErr)
		}
		if err != nil {
			s.http.Logger.Error(err)
		}
	}

	v1.Register(s.http.Group("/v1"))

	addr := fmt.Sprintf(":%s", serverConfig.Port)
	go func() {
		if err := s.http.Start(addr); err != nil {
			s.logger.WithError(err).Info("Shutting down the server now")
		}
	}()
}

func (s *Server) watchStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	s.stop()
}

func (s *Server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.logger.Info("Server is stopping...")
	s.http.Shutdown(ctx)
	close(s.signal)
}
