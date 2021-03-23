package jobupdate

import (
	"fmt"
	"github.com/aeolus3000/lendo-sdk/banking"
	"github.com/aeolus3000/lendo-sdk/messaging"
	"github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"lendo_service/models"
	"time"
)

type Workload struct {
	Message         *messaging.Message
	db *pop.Connection
}

func (w *Workload) Execute() {
	err := w.db.Transaction(w.updateStatus)
	if err != nil {
		log.Errorf("Updating job failed")
	}
}

func (w *Workload) updateStatus(tx *pop.Connection) error {
	// setup logging
	start := tx.Elapsed
	defer func() {
		finished := tx.Elapsed
		elapsed := time.Duration(finished - start)
		log.WithField("updateStatus: Updating transaction took: ", elapsed)
	}()

	bankingApplication, err := banking.DeserializeToApplication(w.Message.Body)
	if err != nil {
		w.acknowledgeMessage("unknown id")
		return fmt.Errorf("updateStatus: can not deserialize application; error = %v", err)
	}
	// Allocate an empty Application
	application := &models.Application{}
	if err := tx.Find(application, bankingApplication.Id); err != nil {
		w.acknowledgeMessage(bankingApplication.Id)
		return fmt.Errorf("updateStatus: can not find id %v in database")
	}
	log.Infof("updateStatus: Updating status of application %v from %v to %v",
		bankingApplication.Id, application.Status, bankingApplication.Status)
	application.Status = bankingApplication.Status

	verrs, err := tx.ValidateAndUpdate(application)
	if err != nil {
		w.acknowledgeMessage(application.ID.String())
		return err
	}
	if verrs.HasAny() {
		w.acknowledgeMessage(application.ID.String())
		return fmt.Errorf("updateStatus: application %v had some validation errors: %v",
			application.ID.String(), verrs.String())
	}

	w.acknowledgeMessage(application.ID.String())
	return nil
}

func (w *Workload) acknowledgeMessage(id string) {
	err := w.Message.Acknowledge()
	if err != nil {
		log.Errorf("acknowledgeMessage: Acknowledge failed for application; id = %v", id)
	}
}
