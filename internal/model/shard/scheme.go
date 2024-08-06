package shard

import (
	"fmt"
	"reflect"
)

var SchemeManager = NewScheme()

type Metadata struct {
	Table     string
	ModelType reflect.Type
	Model     IModel
}

type Scheme struct {
	metadata map[string]*Metadata
}

func (s *Scheme) RegisterType(name string, metadata *Metadata) {
	s.metadata[name] = metadata
}

func (s *Scheme) New(name string) (any, error) {
	t, ok := s.metadata[name]
	if !ok {
		return nil, fmt.Errorf("unrecognized type name: %s", name)
	}
	return reflect.New(t.ModelType.Elem()).Interface(), nil
}

func (s *Scheme) Model(name string) IModel {
	return s.metadata[name].Model
}

func (s *Scheme) Table(name string) string {
	return s.metadata[name].Table
}

func NewScheme() *Scheme {
	return &Scheme{metadata: map[string]*Metadata{}}
}

func NewMetadata(table string, model IModel) *Metadata {
	return &Metadata{Table: table, ModelType: reflect.TypeOf(model), Model: model}
}
