package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingOutbound struct{ ent.Schema }

func (ProxyRoutingOutbound) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_outbound"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingOutbound) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Route Outbound ID"),
		field.String("tag").MaxLen(128).NotEmpty().Unique().Comment("Stable outbound tag"),
		field.String("name").MaxLen(255).NotEmpty().Comment("Outbound display name"),
		field.String("type").MaxLen(32).Default("node_group").Comment("node/node_group/external"),
		field.String("region").MaxLen(32).Default("").Comment("Region code"),
		field.Bool("enabled").Default(true).Comment("Outbound enabled"),
		field.Text("outbound_json").Default("{}").Comment("routing_profile.v1 outbound JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingOutbound) Edges() []ent.Edge { return nil }
