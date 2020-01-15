package spec

import (
	"encoding/json"
	"sync"
)

// Schema is the collection of one or more attributes.
type Schema struct {
	id          string
	name        string
	description string
	attributes  []*Attribute
}

// ID returns the id of the schema.
func (s *Schema) ID() string {
	return s.id
}

// Name returns the name of the schema.
func (s *Schema) Name() string {
	return s.name
}

// Description returns the human-readable text that describes the schema.
func (s *Schema) Description() string {
	return s.description
}

// ForEachAttribute iterate all attributes in this schema and invoke callback function.
func (s *Schema) ForEachAttribute(callback func(attr *Attribute)) {
	for _, attr := range s.attributes {
		callback(attr)
	}
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(schemaJsonAdapter{
		ID:          s.id,
		Name:        s.name,
		Description: s.description,
		Attributes:  s.attributes,
	})
}

func (s *Schema) UnmarshalJSON(raw []byte) error {
	var adapter schemaJsonAdapter
	if err := json.Unmarshal(raw, &adapter); err != nil {
		return err
	}

	s.id = adapter.ID
	s.name = adapter.Name
	s.description = adapter.Description
	s.attributes = adapter.Attributes
	return nil
}

type schemaJsonAdapter struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Attributes  []*Attribute `json:"attributes"`
}

var (
	schemaReg          *schemaRegistry
	schemaRegistryOnce sync.Once
)

type schemaRegistry struct {
	db map[string]*Schema
}

// Register relates the schema with its id in the registry. This method does not check existence of the id and may
// overwrite existing schemas if abused.
func (r *schemaRegistry) Register(schema *Schema) {
	r.db[schema.id] = schema
}

// Get returns the schema that is related to a schemaId, or nil, along with a boolean indicating if the schema exists.
func (r *schemaRegistry) Get(schemaId string) (schema *Schema, ok bool) {
	schema, ok = r.db[schemaId]
	return
}

// Schemas return the schema registry that holds all registered schemas. Use Get and Register to operate the registry.
func Schemas() *schemaRegistry {
	schemaRegistryOnce.Do(func() {
		schemaReg = &schemaRegistry{db: map[string]*Schema{}}
	})
	return schemaReg
}
