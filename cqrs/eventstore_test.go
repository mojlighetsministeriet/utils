package cqrs_test

import (
	"testing"
	"time"

	"github.com/mojlighetsministeriet/utils/cqrs"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type BankAccount struct {
	cqrs.Entity
	Name  string `json:"name",validate:"required"`
	Owner string `json:"owner",validate:"required,uuid"`
}

type BankAccountCreated struct {
	cqrs.EventBase `bson:",inline"`
	Name           string `json:"name",validate:"required"`
	Owner          string `json:"owner",validate:"required,uuid"`
}

func (event *BankAccountCreated) ApplyTo(entity *BankAccount) {
	event.EventBase.ApplyTo(entity)
}

type BankAccountRenamed struct {
	cqrs.EventBase `bson:",inline"`
	Name           string `json:"name",validate:"required"`
}

func TestEventStorePersist(test *testing.T) {
	transport, err := cqrs.NewEventStoreAMQPMongoTransport("localhost", "amqp://guest:guest@localhost:5672/")
	assert.NoError(test, err)

	store := cqrs.NewEventStore(transport)

	at, err := time.Parse(time.RFC3339, "2017-10-12T10:23:12Z")
	assert.NoError(test, err)

	event := BankAccountCreated{
		Name:  "My account",
		Owner: uuid.Must(uuid.NewV4()).String(),
		EventBase: cqrs.EventBase{
			Type:           "BankAccountCreated",
			Version:        1,
			By:             "f040528c-0af8-4457-9228-e4e4793673c7",
			At:             at,
			EntityID:       "287faf1e-c5cc-4844-a1e8-049fb5740d78",
			RepositoryName: "bankaccount",
		},
	}

	validationError, err := store.Persist(event)
	assert.NoError(test, validationError)
	assert.NoError(test, err)

	at, err := time.Parse(time.RFC3339, "2017-10-13T14:41:22Z")
	assert.NoError(test, err)

	event := BankAccountRenamed{
		Name: "The new name",
		EventBase: cqrs.EventBase{
			Type:           "BankAccountRenamed",
			Version:        1,
			By:             "f040528c-0af8-4457-9228-e4e4793673c7",
			At:             at,
			EntityID:       "287faf1e-c5cc-4844-a1e8-049fb5740d78",
			RepositoryName: "bankaccount",
		},
	}

	validationError, err := store.Persist(event)
	assert.NoError(test, validationError)
	assert.NoError(test, err)

	err = store.Transport.RollbackPersist(event)
	assert.NoError(test, err)
}
