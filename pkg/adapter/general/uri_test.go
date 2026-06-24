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

func TestMxUriWithMundoSettings(t *testing.T) {
	for _, transport := range []string{"mundordp", "mundosql"} {
		t.Run(transport, func(t *testing.T) {
			link := MxUri(proxy.Proxy{
				Name:     "Mundo",
				Server:   "example.com",
				Port:     443,
				Protocol: "mx",
				Option: proxy.Mx{
					Transport: transport,
					TransportConfig: proxy.TransportConfig{
						Username:               "admin",
						CertificateFingerprint: "sha256:abc",
						FakeTitle:              "Remote Desktop",
						FakeMessage:            "Access denied",
						AcceptProxyProtocol:    true,
						UseTLSCertificate:      true,
					},
					Security: "tls",
				},
			}, "935b33c7-e128-49f2-816b-71070469cac2")

			parsed, err := url.Parse(link)
			if err != nil {
				t.Fatalf("parse mx uri: %v", err)
			}
			query := parsed.Query()
			for key, want := range map[string]string{
				"type":                   transport,
				"username":               "admin",
				"certificateFingerprint": "sha256:abc",
				"fakeTitle":              "Remote Desktop",
				"fakeMessage":            "Access denied",
				"acceptProxyProtocol":    "true",
				"useTLSCertificate":      "true",
			} {
				if got := query.Get(key); got != want {
					t.Fatalf("%s = %q, want %q; link=%s", key, got, want, link)
				}
			}
		})
	}
}

func TestVlessTrojanUriDoesNotEmitUnsupportedCustomTransports(t *testing.T) {
	const uuid = "935b33c7-e128-49f2-816b-71070469cac2"
	transportConfig := proxy.TransportConfig{
		Path:                   "/custom",
		Host:                   "front.example.com",
		Mode:                   "auto",
		CidrSegments:           []string{"127.0.0.0/24"},
		Username:               "admin",
		CertificateFingerprint: "sha256:def",
		FakeTitle:              "SQL Server",
	}

	for _, transport := range []string{"mc1", "mundordp", "mundosql"} {
		for _, tc := range []struct {
			name     string
			protocol string
			option   interface{}
		}{
			{
				name:     "vless",
				protocol: "vless",
				option: proxy.Vless{
					Transport:       transport,
					TransportConfig: transportConfig,
				},
			},
			{
				name:     "trojan",
				protocol: "trojan",
				option: proxy.Trojan{
					Transport:       transport,
					TransportConfig: transportConfig,
				},
			},
		} {
			t.Run(tc.name+"/"+transport, func(t *testing.T) {
				link := buildProxy(proxy.Proxy{
					Name:     "Custom",
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
				if got := query.Get("type"); got == "mc1" || got == "mundosql" || got == "mundordp" {
					t.Fatalf("type = %q, should not expose custom transport in %s uri; link=%s", got, tc.protocol, link)
				}
				for _, key := range []string{"mode", "cidr", "username", "certificateFingerprint", "fakeTitle", "fakeMessage", "acceptProxyProtocol", "useTLSCertificate"} {
					if got := query.Get(key); got != "" {
						t.Fatalf("%s = %q, should not be emitted in %s uri; link=%s", key, got, tc.protocol, link)
					}
				}
			})
		}
	}
}

func TestVmessUriDoesNotEmitUnsupportedCustomTransport(t *testing.T) {
	for _, transport := range []string{"mc1", "mundordp", "mundosql"} {
		t.Run(transport, func(t *testing.T) {
			link := buildProxy(proxy.Proxy{
				Name:     "Custom",
				Server:   "example.com",
				Port:     443,
				Protocol: "vmess",
				Option: proxy.Vmess{
					Transport: transport,
					TransportConfig: proxy.TransportConfig{
						Path:                   "/custom",
						Host:                   "front.example.com",
						Mode:                   "auto",
						CidrSegments:           []string{"127.0.0.0/24"},
						Username:               "admin",
						CertificateFingerprint: "sha256:abc",
						FakeTitle:              "Remote Desktop",
						AcceptProxyProtocol:    true,
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
			if got, _ := payload["net"].(string); got == "mc1" || got == "mundordp" || got == "mundosql" {
				t.Fatalf("net = %q, should not expose custom transport in vmess payload=%v", got, payload)
			}
			for _, key := range []string{"mode", "cidr", "username", "certificateFingerprint", "fakeTitle", "fakeMessage", "acceptProxyProtocol", "useTLSCertificate"} {
				if _, ok := payload[key]; ok {
					t.Fatalf("%s should not be emitted in vmess payload=%v", key, payload)
				}
			}
		})
	}
}
