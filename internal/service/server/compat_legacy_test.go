package server

import (
	"context"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/npanel-dev/NPanel-backend/ent"
	"github.com/npanel-dev/NPanel-backend/ent/enttest"
	serverbiz "github.com/npanel-dev/NPanel-backend/internal/biz/server"
	"github.com/npanel-dev/NPanel-backend/internal/conf"
	servermodel "github.com/npanel-dev/NPanel-backend/internal/model/server"
	"github.com/redis/go-redis/v9"

	_ "github.com/mattn/go-sqlite3"
)

type compatTestProvider struct {
	db *ent.Client
}

func (p compatTestProvider) DB() *ent.Client {
	return p.db
}

func (p compatTestProvider) Redis() redis.UniversalClient {
	return nil
}

func (p compatTestProvider) Queue() *asynq.Client {
	return nil
}

func (p compatTestProvider) AppNodeConfig() *conf.Node {
	return nil
}

func (p compatTestProvider) LoadNodeConfig(ctx context.Context, module string) (*CompatLegacyNodeConfig, error) {
	return &CompatLegacyNodeConfig{
		NodeSecret:             "secret",
		NodePullInterval:       10,
		NodePushInterval:       60,
		TrafficReportThreshold: 0,
		IPStrategy:             "prefer_ipv4",
	}, nil
}

func TestCompatQueryServerProtocolConfigFiltersEnabledProtocols(t *testing.T) {
	ctx := context.Background()
	client := enttest.Open(t, "sqlite3", "file:compat_query_server_protocols?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	client.ProxyServer.Create().
		SetID(2).
		SetName("HK").
		SetServerAddr("103.214.22.2").
		SetProtocol(`[{"type":"simnet","port":443,"enable":false,"ratio":1,"cert_mode":"none"},{"type":"omniflow","port":443,"enable":false,"security":"tls","ratio":1},{"type":"vmess","port":1566,"enable":true,"security":"none","fingerprint":"chrome","transport":"tcp","ratio":1,"cert_mode":"none"},{"type":"vless","port":1886,"enable":true,"security":"none","fingerprint":"chrome","transport":"tcp","flow":"none","xhttp_mode":"auto","encryption":"none","ratio":1,"cert_mode":"none"}]`).
		SaveX(ctx)

	service := &ServerService{}
	resp, err := service.CompatQueryServerProtocolConfig(ctx, compatTestProvider{db: client}, &CompatLegacyQueryServerConfigRequest{
		ServerID:  2,
		Protocols: []string{"vmess", "vless"},
	})
	if err != nil {
		t.Fatalf("CompatQueryServerProtocolConfig returned error: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("Total = %d, want 2", resp.Total)
	}
	if len(resp.Protocols) != 2 {
		t.Fatalf("Protocols length = %d, want 2", len(resp.Protocols))
	}
	if resp.Protocols[0].Type != "vmess" || resp.Protocols[0].Port != 1566 {
		t.Fatalf("first protocol = %s:%d, want vmess:1566", resp.Protocols[0].Type, resp.Protocols[0].Port)
	}
	if resp.Protocols[1].Type != "vless" || resp.Protocols[1].Port != 1886 {
		t.Fatalf("second protocol = %s:%d, want vless:1886", resp.Protocols[1].Type, resp.Protocols[1].Port)
	}
	if resp.IPStrategy != "prefer_ipv4" {
		t.Fatalf("IPStrategy = %q, want prefer_ipv4", resp.IPStrategy)
	}
}

func TestCompatLegacyProtocolConfigMapIncludesMxMc1TransportConfig(t *testing.T) {
	protocols, err := servermodel.UnmarshalProtocols(`[{"type":"mx","port":443,"enable":true,"security":"reality","sni":"sni.example.com","fingerprint":"chrome","reality_public_key":"public-key","reality_short_id":"short-id","transport":"mc1","path":"/mc1","host":"front.example.com","mode":"auto","cidrSegments":["127.0.0.0/24"]}]`)
	if err != nil {
		t.Fatalf("UnmarshalProtocols returned error: %v", err)
	}
	config := compatLegacyProtocolConfigMap(protocols[0])
	if got, _ := config["transport"].(string); got != "mc1" {
		t.Fatalf("transport = %q, want mc1; config=%v", got, config)
	}
	transportConfig, ok := config["transport_config"].(map[string]interface{})
	if !ok {
		t.Fatalf("transport_config type = %T, want map[string]interface{}", config["transport_config"])
	}
	networkSettings, ok := config["networkSettings"].(map[string]interface{})
	if !ok {
		t.Fatalf("networkSettings type = %T, want map[string]interface{}", config["networkSettings"])
	}
	if got, _ := transportConfig["mode"].(string); got != "auto" {
		t.Fatalf("transport_config.mode = %q, want auto; config=%v", got, transportConfig)
	}
	if got, _ := networkSettings["mode"].(string); got != "auto" {
		t.Fatalf("networkSettings.mode = %q, want auto; config=%v", got, networkSettings)
	}
	cidrSegments, ok := transportConfig["cidrSegments"].([]interface{})
	if !ok || len(cidrSegments) != 1 || cidrSegments[0] != "127.0.0.0/24" {
		t.Fatalf("transport_config.cidrSegments = %#v, want [127.0.0.0/24]", transportConfig["cidrSegments"])
	}
	if _, ok := config["tlsSettings"].(map[string]interface{}); !ok {
		t.Fatalf("tlsSettings type = %T, want map[string]interface{}", config["tlsSettings"])
	}
}

func TestCompatLegacyProtocolConfigMapIncludesMundoTransportConfig(t *testing.T) {
	protocols, err := servermodel.UnmarshalProtocols(`[{"type":"mx","port":443,"enable":true,"security":"tls","transport":"mundordp","username":"admin","certificateFingerprint":"sha256:abc","fakeTitle":"Remote Desktop","fakeMessage":"Access denied","acceptProxyProtocol":true,"useTLSCertificate":true}]`)
	if err != nil {
		t.Fatalf("UnmarshalProtocols returned error: %v", err)
	}
	config := compatLegacyProtocolConfigMap(protocols[0])
	if got, _ := config["transport"].(string); got != "mundordp" {
		t.Fatalf("transport = %q, want mundordp; config=%v", got, config)
	}
	transportConfig, ok := config["transport_config"].(map[string]interface{})
	if !ok {
		t.Fatalf("transport_config type = %T, want map[string]interface{}", config["transport_config"])
	}
	networkSettings, ok := config["networkSettings"].(map[string]interface{})
	if !ok {
		t.Fatalf("networkSettings type = %T, want map[string]interface{}", config["networkSettings"])
	}
	for key, want := range map[string]string{
		"username":               "admin",
		"certificateFingerprint": "sha256:abc",
		"fakeTitle":              "Remote Desktop",
		"fakeMessage":            "Access denied",
	} {
		if got, _ := transportConfig[key].(string); got != want {
			t.Fatalf("transport_config.%s = %q, want %q; config=%v", key, got, want, transportConfig)
		}
		if got, _ := networkSettings[key].(string); got != want {
			t.Fatalf("networkSettings.%s = %q, want %q; config=%v", key, got, want, networkSettings)
		}
	}
	if got, _ := transportConfig["acceptProxyProtocol"].(bool); !got {
		t.Fatalf("transport_config.acceptProxyProtocol = %v, want true; config=%v", got, transportConfig)
	}
	if got, _ := transportConfig["useTLSCertificate"].(bool); !got {
		t.Fatalf("transport_config.useTLSCertificate = %v, want true; config=%v", got, transportConfig)
	}
}

func TestNormalizeMundoProtocolForResponseOnlyKeepsMxMundo(t *testing.T) {
	vless := normalizeMundoProtocolForResponse(&serverbiz.Protocol{
		Type:                        "vless",
		Transport:                   "mundordp",
		MundoUsername:               "admin",
		MundoCertificateFingerprint: "sha256:abc",
		MundoFakeTitle:              "Remote Desktop",
		MundoFakeMessage:            "Denied",
		MundoAcceptProxyProtocol:    true,
		MundoUseTLSCertificate:      true,
	})
	if vless.MundoUsername != "" ||
		vless.MundoCertificateFingerprint != "" ||
		vless.MundoFakeTitle != "" ||
		vless.MundoFakeMessage != "" ||
		vless.MundoAcceptProxyProtocol ||
		vless.MundoUseTLSCertificate {
		t.Fatalf("non-mx protocol kept mundo fields: %+v", vless)
	}

	mx := normalizeMundoProtocolForResponse(&serverbiz.Protocol{
		Type:      "mx",
		Transport: "mundosql",
	})
	if mx.MundoUsername != "MundoUser" {
		t.Fatalf("mx mundo username = %q, want MundoUser", mx.MundoUsername)
	}
}
