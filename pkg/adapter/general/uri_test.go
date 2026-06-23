package general

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/npanel-dev/NPanel-backend/pkg/adapter/proxy"
)

func createServer() proxy.Proxy {
	return proxy.Proxy{
		Name:     "Meta",
		Server:   "127.0.0.1",
		Port:     13092,
		Protocol: "shadowsocks",
		Option: proxy.Shadowsocks{
			Method:    "aes-256-gcm",
			ServerKey: "",
		},
	}
}

func TestGenerateBase64General(t *testing.T) {
	s := createServer()
	p := buildProxy(s, "935b33c7-e128-49f2-816b-71070469cac2")
	t.Log(p)
}

func TestMxUriWithMc1Reality(t *testing.T) {
	node := proxy.Proxy{
		Name:     "Mundo X",
		Server:   "example.com",
		Port:     443,
		Protocol: "mx",
		Option: proxy.Mx{
			Port:      443,
			Transport: "mc1",
			TransportConfig: proxy.TransportConfig{
				Path:            "/path",
				Host:            "front.example.com",
				Mc1Mode:         "auto",
				Mc1CidrSegments: []string{"127.0.0.0/24", "10.0.0.0/8"},
			},
			Security: "reality",
			SecurityConfig: proxy.SecurityConfig{
				SNI:              "sni.example.com",
				Fingerprint:      "chrome",
				RealityPublicKey: "public-key",
				RealityShortId:   "short-id",
			},
		},
	}

	link := MxUri(node, "935b33c7-e128-49f2-816b-71070469cac2")
	if !strings.HasPrefix(link, "mx://935b33c7-e128-49f2-816b-71070469cac2@example.com:443?") {
		t.Fatalf("unexpected mx uri prefix: %s", link)
	}
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse mx uri: %v", err)
	}
	query := parsed.Query()
	assertQuery := func(key, want string) {
		t.Helper()
		if got := query.Get(key); got != want {
			t.Fatalf("query %s = %q, want %q; link=%s", key, got, want, link)
		}
	}
	assertQuery("type", "mc1")
	assertQuery("mode", "auto")
	assertQuery("path", "/path")
	assertQuery("host", "front.example.com")
	assertQuery("cidr", "127.0.0.0/24,10.0.0.0/8")
	assertQuery("security", "reality")
	assertQuery("sni", "sni.example.com")
	assertQuery("servername", "sni.example.com")
	assertQuery("pbk", "public-key")
	assertQuery("sid", "short-id")
	assertQuery("spx", "/")
	assertQuery("fp", "chrome")
}

func TestMxUriWithXhttpDefaultsEndpoint(t *testing.T) {
	link := MxUri(proxy.Proxy{
		Name:     "Mundo X",
		Server:   "example.com",
		Port:     443,
		Protocol: "mx",
		Option: proxy.Mx{
			Transport: "xhttp",
			TransportConfig: proxy.TransportConfig{
				Path:  "/xhttp",
				Split: "rand",
			},
			Security: "tls",
			SecurityConfig: proxy.SecurityConfig{
				SNI: "sni.example.com",
			},
		},
	}, "935b33c7-e128-49f2-816b-71070469cac2")

	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse mx uri: %v", err)
	}
	query := parsed.Query()
	for key, want := range map[string]string{
		"type":          "jatp",
		"endpoint_path": "/xhttp",
		"endpoint":      "example.com",
		"split":         "rand",
		"security":      "tls",
		"sni":           "sni.example.com",
	} {
		if got := query.Get(key); got != want {
			t.Fatalf("%s = %q, want %q; link=%s", key, got, want, link)
		}
	}
}

func TestMxConfigUnmarshalMc1CidrStringAlias(t *testing.T) {
	var mx proxy.Mx
	if err := json.Unmarshal([]byte(`{"port":443,"transport":"mc1","transport_config":{"path":"/mc1","host":"front.example.com","mode":"auto","cidrSegments":"127.0.0.0/24, 10.0.0.0/8"},"security":"none"}`), &mx); err != nil {
		t.Fatalf("unmarshal mx config: %v", err)
	}
	node := proxy.Proxy{
		Name:     "Mundo X",
		Server:   "example.com",
		Port:     mx.Port,
		Protocol: "mx",
		Option:   mx,
	}
	link := MxUri(node, "935b33c7-e128-49f2-816b-71070469cac2")
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse mx uri: %v", err)
	}
	if got := parsed.Query().Get("cidr"); got != "127.0.0.0/24,10.0.0.0/8" {
		t.Fatalf("cidr = %q, want 127.0.0.0/24,10.0.0.0/8; link=%s", got, link)
	}
}

func TestVlessTrojanUriWithMc1Aliases(t *testing.T) {
	const uuid = "935b33c7-e128-49f2-816b-71070469cac2"
	transportConfig := proxy.TransportConfig{
		Path:         "/mc1",
		Host:         "front.example.com",
		Mode:         "auto",
		CidrSegments: []string{"127.0.0.0/24"},
	}

	for _, tc := range []struct {
		name     string
		protocol string
		option   interface{}
	}{
		{
			name:     "vless",
			protocol: "vless",
			option: proxy.Vless{
				Transport:       "mc1",
				TransportConfig: transportConfig,
			},
		},
		{
			name:     "trojan",
			protocol: "trojan",
			option: proxy.Trojan{
				Transport:       "mc1",
				TransportConfig: transportConfig,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			link := buildProxy(proxy.Proxy{
				Name:     "MC1",
				Server:   "example.com",
				Port:     443,
				Protocol: tc.protocol,
				Option:   tc.option,
			}, uuid)
			parsed, err := url.Parse(link)
			if err != nil {
				t.Fatalf("parse %s uri: %v", tc.protocol, err)
			}
			query := parsed.Query()
			if got := query.Get("type"); got != "mc1" {
				t.Fatalf("type = %q, want mc1; link=%s", got, link)
			}
			if got := query.Get("mode"); got != "auto" {
				t.Fatalf("mode = %q, want auto; link=%s", got, link)
			}
			if got := query.Get("cidr"); got != "127.0.0.0/24" {
				t.Fatalf("cidr = %q, want 127.0.0.0/24; link=%s", got, link)
			}
		})
	}
}

func TestVmessUriWithMc1Aliases(t *testing.T) {
	link := buildProxy(proxy.Proxy{
		Name:     "MC1",
		Server:   "example.com",
		Port:     443,
		Protocol: "vmess",
		Option: proxy.Vmess{
			Transport: "mc1",
			TransportConfig: proxy.TransportConfig{
				Path:         "/mc1",
				Host:         "front.example.com",
				Mode:         "auto",
				CidrSegments: []string{"127.0.0.0/24"},
			},
		},
	}, "935b33c7-e128-49f2-816b-71070469cac2")

	encoded := strings.TrimPrefix(link, "vmess://")
	if padding := len(encoded) % 4; padding != 0 {
		encoded += strings.Repeat("=", 4-padding)
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode vmess uri: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(decoded, &payload); err != nil {
		t.Fatalf("unmarshal vmess uri payload: %v", err)
	}
	for key, want := range map[string]string{
		"net":  "mc1",
		"path": "/mc1",
		"host": "front.example.com",
		"mode": "auto",
		"cidr": "127.0.0.0/24",
	} {
		if got, _ := payload[key].(string); got != want {
			t.Fatalf("%s = %q, want %q; payload=%v", key, got, want, payload)
		}
	}
}
