package log

import (
	"blankfactor/event-calendar/config"
	"context"
	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type key string

func (k key) String() string {
	return "key: " + string(k)
}

var (
	log  *logrus.Logger
	once sync.Once
)

const (
	serviceName string = "event-calendar-api"
)

func Init() {
	once.Do(func() {
		logConfig := config.GetEnv().Logger
		log = logrus.New()
		log.Formatter = stackdriver.NewFormatter(
			stackdriver.WithService(serviceName),
		)
		if level, err := logrus.ParseLevel(logConfig.Level); err == nil {
			log.SetLevel(level)
		}
	})
}

func InitParams(ctx context.Context) context.Context {

	httpLog := new(HTTP)
	httpLog.Request = new(Request)
	httpLog.Response = new(Response)

	ctx = Set(ctx, HTTPKey, httpLog)

	return ctx
}

func NewEntry() *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"mutex": &sync.Mutex{},
		"type":  "json",
	})
}

// Get - Get the value associated to the key
func Get(ctx context.Context, key interface{}) (value interface{}) {

	value = ctx.Value(key)

	if value == nil {
		value = reflect.New(reflect.TypeOf(key).Elem()).Interface()
		ctx = context.WithValue(ctx, key, value)
	}

	return value
}

// Set - Put key and value attached to context
func Set(ctx context.Context, key, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}
