package cqrs

import (
	"fmt"

	"github.com/streadway/amqp"
	mgo "gopkg.in/mgo.v2"
)

type RepositoryTransport interface {
	Name() string
	Update(id string) error
	Find(query interface{}, entities []interface{}) error
	FindOne(query interface{}, entity interface{}) error
	FindId(is string, entity interface{}) error
}

type AMQPMongoRepositoryTransport struct {
	name            string
	amqp            *amqp.Connection
	mongoEvents     *mgo.Collection
	mongoRepository *mgo.Collection
}

func (transport *AMQPMongoRepositoryTransport) Find(query interface{}, entities []interface{}) error {
	// TODO: convert raw object to bson.M
	return transport.mongoRepository.Find(query).All(entities)
}

func (transport *AMQPMongoRepositoryTransport) FindOne(query interface{}, entity interface{}) error {
	// TODO: convert raw object to bson.M
	return transport.mongoRepository.Find(query).One(entity)
}

func (transport *AMQPMongoRepositoryTransport) FindId(id string, entity interface{}) error {
	return transport.mongoRepository.FindId(id).One(entity)
}

type EntityUpdater func(eventData interface{}, entity Entity) error

func (transport *AMQPMongoRepositoryTransport) Update(id string, entityUpdater EntityUpdater) (err error) {
	var eventsData []interface{}

	err = transport.mongoEvents.FindId(id).All(&eventsData)
	if err != nil {
		return
	}

	var entity Entity
	for eventData := range eventsData {
		applyError := entityUpdater(eventData, entity)
		if applyError != nil {
			// TODO: add proper logging
			fmt.Println(applyError)
		}
	}

	return
}

func (transport *AMQPMongoRepositoryTransport) connectEventListener() (err error) {
	channel, err := transport.amqp.Channel()
	if err != nil {
		return
	}

	queue, err := channel.QueueDeclare(
		"repository.update."+transport.name, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return
	}

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return
	}

	go func(transport *AMQPMongoRepositoryTransport) {
		// TODO: optimize by looking for a certain time period before triggering update to prevent to many simultainius updates
		for message := range messages {
			updateError := transport.Update(string(message.Body))
			if updateError != nil {
				ackError := message.Ack(true)
				if ackError != nil {
					// TODO: add proper logging
					fmt.Println(updateError)
				}
			} else {
				// TODO: add proper logging
				fmt.Println(updateError)
			}
		}
	}(transport)

	return
}

func NewAMQPMongoRepositoryTransport(name string, mongoURL string, amqpURL string) (transport *AMQPMongoRepositoryTransport, err error) {
	amqpConnection, err := amqp.Dial(amqpURL)
	if err != nil {
		return
	}

	mongoConnection, err := mgo.Dial(mongoURL)
	if err != nil {
		return
	}

	mongoEvents := mongoConnection.DB("eventstore").C("events")

	index := mgo.Index{
		Key:        []string{"repositoryid"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = mongoEvents.EnsureIndex(index)
	if err != nil {
		return
	}

	mongoRepository := mongoConnection.DB("repositories").C(name)

	transport = &AMQPMongoRepositoryTransport{
		name:            name,
		amqp:            amqpConnection,
		mongoEvents:     mongoEvents,
		mongoRepository: mongoRepository,
	}

	transport.connectEventListener()

	return
}
