package jobupdate

import (
	"github.com/aeolus3000/lendo-sdk/messaging"
)

type Configuration struct {
	EsConf      ExecutorServiceConf
	SubConf     messaging.RabbitMqConfiguration
}

type ExecutorServiceConf struct {
	Workers     int `default:"2"`
	QueueLength int `default:"50"`
}

