package cqrs

type Command interface {
	Do() (events []Event, errors []error)
}
