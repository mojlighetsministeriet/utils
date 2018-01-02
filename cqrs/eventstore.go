package cqrs

import (
	"github.com/streadway/amqp"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EventStoreTransport interface {
	Persist(event Event) error
	RollbackPersist(event Event) error
	QueueRepositoryUpdate(event Event) error
}

type EventStore struct {
	Transport EventStoreTransport
}

func (store *EventStore) Persist(event Event) (validationError error, err error) {
	validationError = event.Validate()
	if validationError != nil {
		return
	}

	err = store.Transport.Persist(event)
	if err != nil {
		return
	}

	err = store.Transport.QueueRepositoryUpdate(event)
	if err != nil {
		err = store.Transport.RollbackPersist(event)
	}

	return
}

func NewEventStore(transport EventStoreTransport) *EventStore {
	return &EventStore{Transport: transport}
}

type AMQPMongoTransport struct {
	amqp  *amqp.Connection
	mongo *mgo.Collection
}

func (transport *AMQPMongoTransport) Persist(event Event) error {
	return transport.mongo.Insert(event)
}

func (transport *AMQPMongoTransport) RollbackPersist(event Event) (err error) {
	query := bson.M{
		"type":           event.EventType(),
		"version":        event.EventVersion(),
		"repositoryid":   event.EventEntityID(),
		"repositoryname": event.EventRepositoryName(),
		"at":             event.EventAt(),
		"by":             event.EventBy(),
	}

	err = transport.mongo.Remove(query)
	return
}

func (transport *AMQPMongoTransport) QueueRepositoryUpdate(event Event) (err error) {
	channel, err := transport.amqp.Channel()
	if err != nil {
		return
	}

	queue, err := channel.QueueDeclare(
		"repository.update."+event.EventRepositoryName(), // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event.EventEntityID()),
		},
	)

	return
}

// TODO: handle closing connections when transport is no longer needed

func NewEventStoreAMQPMongoTransport(mongoURL string, amqpURL string) (transport *AMQPMongoTransport, err error) {
	amqpConnection, err := amqp.Dial(amqpURL)
	if err != nil {
		return
	}

	mongoConnection, err := mgo.Dial(mongoURL)
	if err != nil {
		return
	}

	mongoCollection := mongoConnection.DB("eventstore").C("events")

	index := mgo.Index{
		Key:        []string{"repositoryid"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = mongoCollection.EnsureIndex(index)
	if err != nil {
		return
	}

	transport = &AMQPMongoTransport{
		amqp:  amqpConnection,
		mongo: mongoCollection,
	}

	return
}
