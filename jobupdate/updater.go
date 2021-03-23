package jobupdate

import (
	"github.com/aeolus3000/lendo-sdk/configuration"
	"github.com/aeolus3000/lendo-sdk/executor"
	"github.com/aeolus3000/lendo-sdk/messaging"
	"github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"lendo_service/utils"
)

const (
	subConfigPrefix = "subconf"
)

var (
	Sub messaging.Subscriber
)

type Updater struct {
	es              *executor.ExecutorService
	subscriber      messaging.Subscriber
	db *pop.Connection
}



func NewUpdater(db *pop.Connection) *Updater {
	config := configuration.NewDefaultConfiguration()
	updaterConfiguration := Configuration{}
	subErr := config.Process(configPrefix, &updaterConfiguration)
	if subErr != nil {
		log.Fatal(subErr)
	}

	es := executor.NewExecutorService(updaterConfiguration.EsConf.Workers, updaterConfiguration.EsConf.QueueLength)
	sub, err := messaging.NewRabbitMqSubscriber(updaterConfiguration.SubConf, utils.CreateShutdownSignalReceiver())
	if err != nil {
		log.Panicf("NewPoller: Can't start subscriber; error = %v", err)
	}
	return &Updater{
		es:              &es,
		subscriber:      sub,
		db: db,
	}
}

func (p *Updater) Poll() {
	msgs, err := p.subscriber.Consume()
	if err != nil {
		log.Errorf("Poll: Not able to consume messages; error: %v", err)
	}
	for msg := range msgs {
		log.Info("Received job")
		workload := Workload{
			Message:         msg,
			db: p.db,
		}
		p.es.QueueJob(&workload)
	}
}

func (p *Updater) Shutdown() {
	errSub := p.subscriber.Close()
	if errSub != nil {
		log.Errorf("Failed to shutdown sub; error = %v", errSub)
	}
}