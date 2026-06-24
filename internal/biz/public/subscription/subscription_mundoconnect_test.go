package subscription

import (
	"context"
	"encoding/json"
	"testing"

	v1 "github.com/npanel-dev/NPanel-backend/api/public/subscription/v1"
)

type subscribeMundoRepoStub struct {
	nodes []*NodeInfo
}

func (s *subscribeMundoRepoStub) ValidateTokenAndGetSubscribe(ctx context.Context, token string) (*UserSubscribe, error) {
	return &UserSubscribe{
		ID:            99,
		UserID:        1,
		SubscribeID:   10,
		Token:         token,
		UUID:          "935b33c7-e128-49f2-816b-71070469cac2",
		SubscribeName: "Mundo",
		Status:        1,
	}, nil
}

func (s *subscribeMundoRepoStub) GetAvailableNodes(ctx context.Context, userSubscribe *UserSubscribe) ([]*NodeInfo, error) {
	return s.nodes, nil
}

func (s *subscribeMundoRepoStub) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	return &UserInfo{ID: userID, Email: "user@example.com"}, nil
}

func (s *subscribeMundoRepoStub) GetSubscribeInfo(ctx context.Context, userSubscribe *UserSubscribe) string {
	return ""
}

func (s *subscribeMundoRepoStub) UpdateSubscribeLog(ctx context.Context, userSubscribe *UserSubscribe, userAgent, clientIP string) error {
	return nil
}

func (s *subscribeMundoRepoStub) GetSubscribeApplications(ctx context.Context) ([]*SubscribeApplication, error) {
	return []*SubscribeApplication{
		{
			Name:              "General",
			UserAgent:         "clash",
			IsDefault:         true,
			OutputFormat:      "json",
			SubscribeTemplate: `{{ toJson .Proxies }}`,
		},
		{
			Name:              "Mundo Connect",
			UserAgent:         "mundoconnect",
			OutputFormat:      "json",
			SubscribeTemplate: `{{ toJson .Proxies }}`,
		},
	}, nil
}

func (s *subscribeMundoRepoStub) GetSubscribeDomain(ctx context.Context) string {
	return ""
}

func (s *subscribeMundoRepoStub) GetSubscribePath(ctx context.Context) string {
	return "/s"
}

func (s *subscribeMundoRepoStub) GetSiteName(ctx context.Context) string {
	return "Mundo"
}

func (s *subscribeMundoRepoStub) GetSubscribeRuntimeConfig(ctx context.Context) (*SubscribeRuntimeConfig, error) {
	return &SubscribeRuntimeConfig{}, nil
}

func TestGetSubscribeConfigMundoConnectKeepsMxAndMundoTransports(t *testing.T) {
	uc := NewSubscriptionUseCase(&subscribeMundoRepoStub{nodes: mundoConnectTestNodes()})

	reply, err := uc.GetSubscribeConfig(
		context.Background(),
		&v1.GetSubscribeConfigRequest{Token: "token"},
		"mundoconnect/1.0",
		"127.0.0.1",
		"/s/token",
		"sub.example.com",
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("GetSubscribeConfig returned error: %v", err)
	}

	nodes := decodeRenderedNodes(t, reply.Config)
	if len(nodes) != 5 {
		t.Fatalf("rendered nodes = %d, want 5; nodes=%v", len(nodes), nodes)
	}
	assertRenderedNode(t, nodes, "Mundo X", "mx", "mc1")
	assertRenderedNode(t, nodes, "MC1 VLESS", "vless", "mc1")
	assertRenderedNode(t, nodes, "Mundo RDP", "mx", "mundordp")
	assertRenderedNode(t, nodes, "Mundo SQL", "mx", "mundosql")
	assertRenderedNode(t, nodes, "Trojan", "trojan", "tcp")
}

func TestGetSubscribeConfigThirdPartyHidesMxAndMundoTransports(t *testing.T) {
	uc := NewSubscriptionUseCase(&subscribeMundoRepoStub{nodes: mundoConnectTestNodes()})

	reply, err := uc.GetSubscribeConfig(
		context.Background(),
		&v1.GetSubscribeConfigRequest{Token: "token"},
		"clash-verge/1.0",
		"127.0.0.1",
		"/s/token",
		"sub.example.com",
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("GetSubscribeConfig returned error: %v", err)
	}

	nodes := decodeRenderedNodes(t, reply.Config)
	if len(nodes) != 1 {
		t.Fatalf("rendered nodes = %d, want only regular node; nodes=%v", len(nodes), nodes)
	}
	assertRenderedNode(t, nodes, "Trojan", "trojan", "tcp")
}

func mundoConnectTestNodes() []*NodeInfo {
	return []*NodeInfo{
		{
			ID:        1,
			Name:      "Mundo X",
			Server:    "mx.example.com",
			Port:      443,
			Type:      "mx",
			Transport: "mc1",
			Path:      "/mc1",
			Host:      "front.example.com",
			Mc1Mode:   "auto",
			Security:  "tls",
			SNI:       "sni.example.com",
		},
		{
			ID:        2,
			Name:      "MC1 VLESS",
			Server:    "vless.example.com",
			Port:      443,
			Type:      "vless",
			Transport: "mc1",
			Path:      "/mc1",
			Host:      "front.example.com",
		},
		{
			ID:            3,
			Name:          "Mundo RDP",
			Server:        "rdp.example.com",
			Port:          443,
			Type:          "mx",
			Transport:     "mundordp",
			MundoUsername: "MundoUser",
		},
		{
			ID:        4,
			Name:      "Mundo SQL",
			Server:    "sql.example.com",
			Port:      3306,
			Type:      "mx",
			Transport: "mundosql",
		},
		{
			ID:        5,
			Name:      "Trojan",
			Server:    "trojan.example.com",
			Port:      443,
			Type:      "trojan",
			Transport: "tcp",
		},
	}
}

func decodeRenderedNodes(t *testing.T, payload []byte) []map[string]interface{} {
	t.Helper()
	var nodes []map[string]interface{}
	if err := json.Unmarshal(payload, &nodes); err != nil {
		t.Fatalf("unmarshal rendered nodes: %v; payload=%s", err, string(payload))
	}
	return nodes
}

func assertRenderedNode(t *testing.T, nodes []map[string]interface{}, name, protocolType, transport string) {
	t.Helper()
	for _, node := range nodes {
		if node["Name"] != name {
			continue
		}
		if got, _ := node["Type"].(string); got != protocolType {
			t.Fatalf("%s Type = %q, want %q; node=%v", name, got, protocolType, node)
		}
		if got, _ := node["Transport"].(string); got != transport {
			t.Fatalf("%s Transport = %q, want %q; node=%v", name, got, transport, node)
		}
		return
	}
	t.Fatalf("node %q not rendered; nodes=%v", name, nodes)
}
