package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
	"github.com/gobuffalo/validate/v3/validators"
)
// Application is used by pop to map your applications database table to your go code.
type Application struct {
    ID uuid.UUID `json:"id" db:"id"`
    FirstName string `json:"first_name" db:"first_name"`
    LastName string `json:"last_name" db:"last_name"`
    Status string `json:"status" db:"status"`
    JobID nulls.String `json:"job_id" db:"job_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a Application) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Applications is not required by pop and may be deleted
type Applications []Application

// String is not required by pop and may be deleted
func (a Applications) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Application) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: a.LastName, Name: "LastName"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Application) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Application) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
