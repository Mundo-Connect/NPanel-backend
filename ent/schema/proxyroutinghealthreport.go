package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingHealthReport struct{ ent.Schema }

func (ProxyRoutingHealthReport) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_health_report"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingHealthReport) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Routing health report ID"),
		field.String("reporter_type").MaxLen(32).Default("client").Comment("client/node/backend"),
		field.String("reporter_id").MaxLen(128).Default("").Comment("Reporter identifier"),
		field.String("profile_code").MaxLen(128).Default("").Comment("Routing profile code"),
		field.String("routing_hash").MaxLen(128).Default("").Comment("Routing hash"),
		field.String("subject_type").MaxLen(32).NotEmpty().Comment("outbound/dns_resolver/service"),
		field.String("subject_key").MaxLen(128).NotEmpty().Comment("Outbound tag, DNS resolver tag or service code"),
		field.String("region").MaxLen(32).Default("").Comment("Region code"),
		field.String("status").MaxLen(32).Default("unknown").Comment("healthy/ok/failed/degraded/stale/disabled/unknown"),
		field.String("source").MaxLen(64).Default("health_report").Comment("Health source"),
		field.Int("rtt_ms").Default(0).Comment("RTT in milliseconds"),
		field.Int("consecutive_failures").Default(0).Comment("Consecutive failures"),
		field.Text("last_error").Optional().Comment("Last health error"),
		field.String("outbound_tag").MaxLen(128).Default("").Comment("Related outbound tag"),
		field.String("dns_resolver_tag").MaxLen(128).Default("").Comment("Related DNS resolver tag"),
		field.Time("checked_at").Default(time.Now).Comment("Checked at"),
		field.Text("report_json").Default("{}").Comment("Raw health report JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingHealthReport) Edges() []ent.Edge { return nil }
