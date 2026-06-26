package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProxyRoutingProfile struct{ ent.Schema }

func (ProxyRoutingProfile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "routing_profile"},
		entsql.WithComments(true),
	}
}

func (ProxyRoutingProfile) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Routing Profile ID"),
		field.String("code").MaxLen(128).NotEmpty().Unique().Comment("Stable profile code"),
		field.String("name").MaxLen(255).NotEmpty().Comment("Profile display name"),
		field.Text("description").Optional().Comment("Profile description"),
		field.String("scope_type").MaxLen(32).Default("global").Comment("user/plan/group/node/global"),
		field.String("scope_id").MaxLen(128).Default("default").Comment("Scope identifier"),
		field.Int("priority").Default(100).Comment("Lower priority matches first"),
		field.String("mode").MaxLen(32).Default("observe").Comment("off/observe/enforce"),
		field.Bool("enabled").Default(true).Comment("Profile enabled"),
		field.Text("profile_json").Default("{}").Comment("routing_profile.v1 profile object JSON"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Created at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Updated at"),
	}
}

func (ProxyRoutingProfile) Edges() []ent.Edge { return nil }
