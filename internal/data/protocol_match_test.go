package data

import (
	"testing"

	servermodel "github.com/npanel-dev/NPanel-backend/internal/model/server"
)

func TestMatchNodeProtocolConfigPrefersSameTypeAndPort(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1"},
		{Type: "mx", Port: 3389, Enable: true, Transport: "mundordp"},
		{Type: "mx", Port: 3306, Enable: true, Transport: "mundosql"},
	}

	for _, tc := range []struct {
		name      string
		port      uint16
		transport string
	}{
		{name: "mc1", port: 443, transport: "mc1"},
		{name: "rdp", port: 3389, transport: "mundordp"},
		{name: "sql", port: 3306, transport: "mundosql"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			matched, _, _ := matchNodeProtocolConfig(protocols, "mx", tc.port)
			if matched == nil {
				t.Fatal("matched protocol is nil")
			}
			if matched.Transport != tc.transport {
				t.Fatalf("matched transport = %q, want %q", matched.Transport, tc.transport)
			}
		})
	}
}

func TestMatchNodeProtocolConfigFallsBackToSameType(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1"},
		{Type: "vless", Port: 8443, Enable: true, Transport: "tcp"},
	}

	matched, firstEnabled, firstAvailable := matchNodeProtocolConfig(protocols, "mx", 9443)
	if matched == nil || matched.Transport != "mc1" {
		t.Fatalf("matched = %+v, want mx/mc1 fallback", matched)
	}
	if firstEnabled != protocols[0] || firstAvailable != protocols[0] {
		t.Fatalf("unexpected fallbacks: enabled=%+v available=%+v", firstEnabled, firstAvailable)
	}
}

func TestMatchNodeProtocolConfigReturnsFallbacksWhenTypeMissing(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "vless", Port: 443, Enable: false},
		{Type: "trojan", Port: 8443, Enable: true},
	}

	matched, firstEnabled, firstAvailable := matchNodeProtocolConfig(protocols, "mx", 443)
	if matched != nil {
		t.Fatalf("matched = %+v, want nil for missing type", matched)
	}
	if firstEnabled != protocols[1] {
		t.Fatalf("firstEnabled = %+v, want trojan", firstEnabled)
	}
	if firstAvailable != protocols[0] {
		t.Fatalf("firstAvailable = %+v, want vless", firstAvailable)
	}
}
