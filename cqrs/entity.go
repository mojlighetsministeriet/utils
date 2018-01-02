package cqrs

import (
	"time"
)

type Entity interface {
	EntityID() string
	EntityUpdatedAt() time.Time
	EntityCreatedAt() time.Time
	EntityDeletedAt() time.Time
	EntityUpdatedBy() string
	EntityCreatedBy() string
	EntityDeletedBy() string
}

type EntityModel struct {
	ID        string
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedBy string
	CreatedBy string
	DeletedBy string
}

func (model EntityModel) EntityID() string {
	return model.ID
}

func (model EntityModel) EntityUpdatedAt() time.Time {
	return model.UpdatedAt
}

func (model EntityModel) EntityCreatedAt() time.Time {
	return model.CreatedAt
}

func (model EntityModel) EntityDeletedAt() time.Time {
	return model.DeletedAt
}

func (model EntityModel) EntityUpdatedBy() string {
	return model.UpdatedBy
}

func (model EntityModel) EntityCreatedBy() string {
	return model.CreatedBy
}

func (model EntityModel) EntityDeletedBy() string {
	return model.DeletedBy
}
