package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Series holds the schema definition for the Series entity.
type Series struct {
	ent.Schema
}

// Fields of the Series.
func (Series) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("url").NotEmpty(),
		field.String("author").Optional(),
		field.String("cover").Optional(),
		field.String("description").Optional(),
	}
}

// Edges of the Series.
func (Series) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("books", Book.Type),
	}
}
