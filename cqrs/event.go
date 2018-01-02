package cqrs

import (
	"encoding/json"
	"time"

	"github.com/mojlighetsministeriet/utils/jsonvalidator"
)

type Event interface {
	EventType() string
	EventVersion() uint
	EventEntityID() string
	EventRepositoryName() string
	EventAt() time.Time
	EventBy() string
	JSON() []byte
	Validate() error
	ApplyTo(entity Entity)
}

type EventBase struct {
	Type           string    `json:"type",validate:"required"`
	Version        uint      `json:"version",validate:"required,min=1"`
	EntityID       string    `json:"entityId",validate:"required,uuid"`
	RepositoryName string    `json:"repositoryName",validate:"required"`
	At             time.Time `json:"at",validate:"required,date-time"`
	By             string    `json:"by",validate:"required,uuid"`
}

func (event EventBase) EventType() string {
	return event.Type
}

func (event EventBase) EventVersion() uint {
	return event.Version
}

func (event EventBase) EventEntityID() string {
	return event.EntityID
}

func (event EventBase) EventRepositoryName() string {
	return event.RepositoryName
}

func (event EventBase) EventAt() time.Time {
	return event.At
}

func (event EventBase) EventBy() string {
	return event.By
}

func (event EventBase) JSON() []byte {
	data, _ := json.Marshal(event)
	return data
}

func (event EventBase) Validate() error {
	validator := jsonvalidator.NewValidator()
	return validator.Validate(event)
}

func (event EventBase) ApplyTo(entity Entity) {
	model := entity.(EntityModel)

	if model.CreatedAt.IsZero() == true {
		model.CreatedAt = event.At
	}

	if model.CreatedBy == "" {
		model.CreatedBy = event.By
	}

	model.UpdatedAt = event.At
	model.UpdatedBy = event.By
}
