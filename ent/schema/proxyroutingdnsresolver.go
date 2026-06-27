package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingDNSResolver struct{ ent.Schema }

func (ProxyRoutingDNSResolver) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_dns_resolver"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingDNSResolver) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("DNS Resolver ID"),
		field.String("tag").MaxLen(128).NotEmpty().Unique().Comment("Stable resolver tag"),
		field.String("name").MaxLen(255).NotEmpty().Comment("Resolver display name"),
		field.String("proto").MaxLen(32).Default("doh").Comment("doh/dot/udp/tcp"),
		field.String("address").MaxLen(512).NotEmpty().Comment("Resolver address"),
		field.Int("port").Default(443).Comment("Resolver port"),
		field.Bool("enabled").Default(true).Comment("Resolver enabled"),
		field.Text("resolver_json").Default("{}").Comment("routing_profile.v1 DNS resolver JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingDNSResolver) Edges() []ent.Edge { return nil }
