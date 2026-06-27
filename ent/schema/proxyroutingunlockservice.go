package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingUnlockService struct{ ent.Schema }

func (ProxyRoutingUnlockService) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_unlock_service"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingUnlockService) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Unlock Service ID"),
		field.String("code").MaxLen(128).NotEmpty().Unique().Comment("Stable service code"),
		field.String("name").MaxLen(255).NotEmpty().Comment("Service display name"),
		field.String("category").MaxLen(128).Default("").Comment("Service category"),
		field.Bool("enabled").Default(true).Comment("Service enabled"),
		field.Text("service_json").Default("{}").Comment("routing_profile.v1 unlock service JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingUnlockService) Edges() []ent.Edge { return nil }
