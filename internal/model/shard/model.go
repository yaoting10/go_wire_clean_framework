package shard

import "time"

type IModel interface {
	GetID() uint
	GetCreatedAt() *time.Time
	GetUpdatedAt() *time.Time
}

type BaseModel struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BaseModel) GetID() uint {
	return b.ID
}

func (b *BaseModel) GetCreatedAt() *time.Time {
	return &b.CreatedAt
}

func (b *BaseModel) GetUpdatedAt() *time.Time {
	return &b.UpdatedAt
}
