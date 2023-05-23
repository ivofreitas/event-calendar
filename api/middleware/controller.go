package middleware

import (
	"blankfactor/event-calendar/log"
	"blankfactor/event-calendar/model"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

type HandlerFn func(ctx context.Context, param interface{}) (interface{}, error)

type controller struct {
	fn         HandlerFn
	param      interface{}
	httpStatus int
}

func NewController(fn HandlerFn, httpStatus int, param interface{}) *controller {
	c := controller{fn, param, httpStatus}
	return &c
}

func (ctrl *controller) Handle(c echo.Context) error {

	ctx := c.Request().Context()
	httpLog := log.Get(ctx, log.HTTPKey).(*log.HTTP)

	if ctrl.param != nil {
		ctrl.param = reflect.New(reflect.TypeOf(ctrl.param).Elem()).Interface()
		if err := c.Bind(ctrl.param); err != nil {
			responseErr := err.(*model.ResponseError)
			httpLog.Error = responseErr.Error()
			return c.JSON(http.StatusBadRequest, responseErr)
		}

		if err := c.Validate(ctrl.param); err != nil {
			responseErr := err.(*model.ResponseError)
			httpLog.Error = responseErr.Error()
			return c.JSON(http.StatusBadRequest, responseErr)
		}

		b, _ := json.Marshal(ctrl.param)
		httpLog.Request.Param = string(b)
	}

	result, err := ctrl.fn(ctx, ctrl.param)
	if err != nil {
		responseErr := err.(*model.ResponseError)
		httpLog.Error = responseErr.Error()
		return c.JSON(responseErr.StatusCode, responseErr)
	}

	httpLog.Level = logrus.InfoLevel
	if result != nil {
		httpLog.Response.Body = result
		return c.JSON(ctrl.httpStatus, result)
	}

	return c.JSON(http.StatusNoContent, nil)
}
