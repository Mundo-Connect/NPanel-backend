package proxy

import (
	"embed"
	"encoding/json"
	"strings"
)

// Adapter represents a proxy adapter
type Adapter struct {
	Proxies    []Proxy
	Group      []Group
	Rules      []string  // rule
	Nodes      []string  // all node
	Default    string    // Default Node
	TemplateFS *embed.FS // Template file system
}

// Proxy represents a proxy server
type Proxy struct {
	Name     string   // Name of the proxy
	Server   string   // Server address of the proxy
	Port     int      // Port of the proxy server
	Protocol string   // Protocol type (e.g., shadowsocks, vless, vmess, trojan, hysteria2, tuic, anytls)
	Country  string   // Country of the proxy
	Tags     []string // Tags for the proxy
	Option   any      // Additional options for the proxy configuration
}

// Group represents a group of proxies
type Group struct {
	Name     string
	Type     GroupType
	Proxies  []string
	URL      string
	Interval int
	Reject   bool     // Reject group
	Direct   bool     // Direct group
	Tags     []string // Tags for the group
	Default  bool     // Default group
}

type GroupType string

const (
	GroupTypeSelect   GroupType = "select"
	GroupTypeURLTest  GroupType = "url-test"
	GroupTypeFallback GroupType = "fallback"
)

// Shadowsocks represents a Shadowsocks proxy configuration
type Shadowsocks struct {
	Port      int    `json:"port"`
	Method    string `json:"method"`
	ServerKey string `json:"server_key"`
}

// Vless represents a Vless proxy configuration
type Vless struct {
	Port            int             `json:"port"`
	Flow            string          `json:"flow"`
	Transport       string          `json:"transport"`
	TransportConfig TransportConfig `json:"transport_config"`
	Security        string          `json:"security"`
	SecurityConfig  SecurityConfig  `json:"security_config"`
}

// Vmess represents a Vmess proxy configuration
type Vmess struct {
	Port            int             `json:"port"`
	Flow            string          `json:"flow"`
	Transport       string          `json:"transport"`
	TransportConfig TransportConfig `json:"transport_config"`
	Security        string          `json:"security"`
	SecurityConfig  SecurityConfig  `json:"security_config"`
}

// Trojan represents a Trojan proxy configuration
type Trojan struct {
	Port            int             `json:"port"`
	Flow            string          `json:"flow"`
	Transport       string          `json:"transport"`
	TransportConfig TransportConfig `json:"transport_config"`
	Security        string          `json:"security"`
	SecurityConfig  SecurityConfig  `json:"security_config"`
}

// Hysteria2 represents a Hysteria2 proxy configuration
type Hysteria2 struct {
	Port           int            `json:"port"`
	HopPorts       string         `json:"hop_ports"`
	HopInterval    int            `json:"hop_interval"`
	ObfsPassword   string         `json:"obfs_password"`
	SecurityConfig SecurityConfig `json:"security_config"`
}

// Tuic represents a Tuic proxy configuration
type Tuic struct {
	Port                 int            `json:"port"`
	DisableSNI           bool           `json:"disable_sni"`
	ReduceRtt            bool           `json:"reduce_rtt"`
	UDPRelayMode         string         `json:"udp_relay_mode"`
	CongestionController string         `json:"congestion_controller"`
	SecurityConfig       SecurityConfig `json:"security_config"`
}

// AnyTLS represents an AnyTLS proxy configuration
type AnyTLS struct {
	Port           int            `json:"port"`
	SecurityConfig SecurityConfig `json:"security_config"`
}

// Mx represents a Mundo X proxy configuration.
type Mx struct {
	Port            int             `json:"port"`
	Transport       string          `json:"transport"`
	TransportConfig TransportConfig `json:"transport_config"`
	Security        string          `json:"security"`
	SecurityConfig  SecurityConfig  `json:"security_config"`
}

// TransportConfig represents the transport configuration for a proxy
type TransportConfig struct {
	Path                 string   `json:"path,omitempty"` // ws/httpupgrade
	Host                 string   `json:"host,omitempty"`
	ServiceName          string   `json:"service_name"`          // grpc
	DisableSNI           bool     `json:"disable_sni"`           // Disable SNI for the transport(tuic)
	ReduceRtt            bool     `json:"reduce_rtt"`            // Reduce RTT for the transport(tuic)
	UDPRelayMode         string   `json:"udp_relay_mode"`        // UDP relay mode for the transport(tuic)
	CongestionController string   `json:"congestion_controller"` // Congestion controller for the transport(tuic)
	Mc1Mode              string   `json:"mc1_mode,omitempty"`
	Mc1CidrSegments      []string `json:"mc1_cidr_segments,omitempty"`
	Mode                 string   `json:"mode,omitempty"`
	CidrSegments         []string `json:"cidrSegments,omitempty"`
	Split                string   `json:"split,omitempty"`
}

func (c *TransportConfig) UnmarshalJSON(data []byte) error {
	type transportConfigAlias TransportConfig
	aux := struct {
		*transportConfigAlias
		Mc1CidrSegments mc1CIDRSegments `json:"mc1_cidr_segments,omitempty"`
		CidrSegments    mc1CIDRSegments `json:"cidrSegments,omitempty"`
	}{
		transportConfigAlias: (*transportConfigAlias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Mc1CidrSegments) > 0 {
		c.Mc1CidrSegments = []string(aux.Mc1CidrSegments)
	}
	if len(aux.CidrSegments) > 0 {
		c.CidrSegments = []string(aux.CidrSegments)
	}
	return nil
}

type mc1CIDRSegments []string

func (s *mc1CIDRSegments) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = sanitizeCIDRSegments(strings.Split(single, ","))
		return nil
	}
	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	*s = sanitizeCIDRSegments(values)
	return nil
}

func sanitizeCIDRSegments(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}

// SecurityConfig represents the security configuration for a proxy
type SecurityConfig struct {
	SNI               string `json:"sni"`
	AllowInsecure     bool   `json:"allow_insecure"`
	Fingerprint       string `json:"fingerprint"`
	RealityServerAddr string `json:"reality_server_addr"`
	RealityServerPort int    `json:"reality_server_port"`
	RealityPrivateKey string `json:"reality_private_key"`
	RealityPublicKey  string `json:"reality_public_key"`
	RealityShortId    string `json:"reality_short_id"`
}

// Relay represents a relay configuration
type Relay struct {
	RelayHost    string
	DispatchMode string
	Prefix       string
}
