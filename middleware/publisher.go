package middleware

import (
	"github.com/aeolus3000/lendo-sdk/configuration"
	"github.com/aeolus3000/lendo-sdk/messaging"
	"github.com/gobuffalo/buffalo"
	"lendo_service/utils"
	"log"
)

const (
	CtxPublisher = "publisher"
	pubConfigPrefix = "pubconf"
)

var (
	Pub messaging.Publisher
)

func init() {
	config := configuration.NewDefaultConfiguration()
	pubConfig := messaging.RabbitMqConfiguration{}
	pubErr := config.Process(configPrefix + "_" + pubConfigPrefix, &pubConfig)
	if pubErr != nil {
		log.Fatal(pubErr)
	}

	Pub, pubErr = messaging.NewRabbitMqPublisher(pubConfig, utils.CreateShutdownSignalReceiver())
	if pubErr != nil {
		log.Fatal(pubErr)
	}
}


func PubMiddleware(publisher messaging.Publisher) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			c.Set(CtxPublisher, publisher)
			err := next(c)
			return err
		}
	}
}