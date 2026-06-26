package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingRouteEvent struct{ ent.Schema }

func (ProxyRoutingRouteEvent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_route_event"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingRouteEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Routing route event ID"),
		field.String("reporter_type").MaxLen(32).Default("client").Comment("client/node/backend"),
		field.String("reporter_id").MaxLen(128).Default("").Comment("Reporter identifier"),
		field.String("profile_code").MaxLen(128).Default("").Comment("Routing profile code"),
		field.String("routing_hash").MaxLen(128).Default("").Comment("Routing hash"),
		field.String("event_type").MaxLen(64).NotEmpty().Comment("route_decision/route_fallback/outbound_health_changed/dns_resolver_health_changed"),
		field.String("subject").MaxLen(256).Default("").Comment("Domain, IP, service or target"),
		field.String("rule_id").MaxLen(128).Default("").Comment("Matched rule ID"),
		field.String("rule_name").MaxLen(128).Default("").Comment("Matched rule name"),
		field.String("action_type").MaxLen(32).Default("").Comment("direct/proxy/reject/dns_resolver/outbound"),
		field.String("outbound_tag").MaxLen(128).Default("").Comment("Outbound tag"),
		field.String("dns_resolver_tag").MaxLen(128).Default("").Comment("DNS resolver tag"),
		field.String("fallback_target").MaxLen(128).Default("").Comment("Fallback target"),
		field.String("status").MaxLen(32).Default("unknown").Comment("matched/fallback/healthy/failed/degraded/unknown"),
		field.Int("latency_ms").Default(0).Comment("Latency in milliseconds"),
		field.Text("error").Optional().Comment("Event error"),
		field.Time("event_at").Default(time.Now).Comment("Event time"),
		field.Text("event_json").Default("{}").Comment("Raw route event JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingRouteEvent) Edges() []ent.Edge { return nil }
