package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MinLen(13).MaxLen(13).NotEmpty().Unique().Immutable().StructTag("json:\"isbn\""),
		field.String("title").NotEmpty(),
		field.Int("number").Positive().StructTag("json:\"tome_number\""),
		field.String("url").NotEmpty(),
		field.Bool("owned"),
		field.Time("updated_at").Default(time.Now),
		field.Time("released_at").Default(time.Now),
		field.Int("series_id").Positive(),
		field.String("cover").Optional(),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("series", Series.Type).Ref("books").Unique().Field("series_id").Required(),
	}
}
