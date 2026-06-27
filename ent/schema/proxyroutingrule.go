package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingRule struct{ ent.Schema }

func (ProxyRoutingRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_rule"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingRule) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Routing Rule ID"),
		field.Int64("profile_id").Default(0).Comment("Bound routing_profile.id, 0 means default P1 profile"),
		field.String("name").MaxLen(255).NotEmpty().Comment("Rule display name"),
		field.Int("priority").Default(100).Comment("Lower priority matches first"),
		field.Bool("enabled").Default(true).Comment("Rule enabled"),
		field.String("service_code").MaxLen(128).Default("").Comment("Unlock service code"),
		field.Text("matcher_json").Default("{}").Comment("routing_profile.v1 matcher JSON"),
		field.Text("action_json").Default("{}").Comment("routing_profile.v1 action JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingRule) Edges() []ent.Edge { return nil }
