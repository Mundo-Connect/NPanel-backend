package general

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/npanel-dev/NPanel-backend/pkg/adapter/proxy"
)

type v2rayShareLink struct {
	Ps            string `json:"ps"`
	Add           string `json:"add"`
	Port          string `json:"port"`
	ID            string `json:"id"`
	Aid           string `json:"aid"`
	Net           string `json:"net"`
	Type          string `json:"type"`
	Host          string `json:"host"`
	SNI           string `json:"sni"`
	Path          string `json:"path"`
	TLS           string `json:"tls"`
	Flow          string `json:"flow,omitempty"`
	Alpn          string `json:"alpn,omitempty"`
	AllowInsecure bool   `json:"allowInsecure"`
	Fingerprint   string `json:"fp,omitempty"`
	PublicKey     string `json:"pbk,omitempty"`
	ShortId       string `json:"sid,omitempty"`
	SpiderX       string `json:"spx,omitempty"`
	Mode          string `json:"mode,omitempty"`
	Cidr          string `json:"cidr,omitempty"`
	V             string `json:"v"`
}

// GenerateBase64General  will output node URLs split by '\n' and then encode into base64
func GenerateBase64General(data []proxy.Proxy, uuid string) []byte {
	var links []string
	for _, v := range data {
		p := buildProxy(v, uuid)
		if p == "" {
			continue
		}
		links = append(links, p)
	}
	var rsp []byte
	rsp = base64.RawStdEncoding.AppendEncode(rsp, []byte(strings.Join(links, "\r\n")))
	return rsp
}

func buildProxy(data proxy.Proxy, uuid string) string {
	switch data.Protocol {
	case "shadowsocks":
		return ShadowsocksUri(data, uuid)
	case "vmess":
		return VmessUri(data, uuid)
	case "vless":
		return VlessUri(data, uuid)
	case "trojan":
		return TrojanUri(data, uuid)
	case "hysteria2":
		return Hysteria2Uri(data, uuid)
	case "tuic":
		return TuicUri(data, uuid)
	case "mx":
		return MxUri(data, uuid)
	default:
		return ""
	}
}

func ShadowsocksUri(data proxy.Proxy, uuid string) string {
	ss, ok := data.Option.(proxy.Shadowsocks)
	if !ok {
		return ""
	}

	password := uuid
	// SIP022 AEAD-2022 Ciphers
	if strings.Contains(ss.Method, "2022") {
		serverKey, userKey := proxy.GenerateShadowsocks2022Password(ss, uuid)
		password = fmt.Sprintf("%s:%s", serverKey, userKey)
	}

	u := &url.URL{
		Scheme:   "ss",
		User:     url.User(strings.TrimSuffix(base64.URLEncoding.EncodeToString([]byte(ss.Method+":"+password)), "=")),
		Host:     net.JoinHostPort(data.Server, strconv.Itoa(data.Port)),
		Fragment: data.Name,
	}
	return u.String()
}

func VmessUri(data proxy.Proxy, uuid string) string {
	vmess := data.Option.(proxy.Vmess)

	transport := vmess.TransportConfig

	securityConfig := vmess.SecurityConfig

	var s = v2rayShareLink{
		V:    "2",
		Add:  data.Server,
		Port: fmt.Sprint(data.Port),
		ID:   uuid,
		Aid:  "0",
		Ps:   data.Name,
		Net:  "tcp",
	}

	switch vmess.Transport {
	case "websocket":
		s.Net = "ws"
		s.Path = transport.Path
		s.Host = transport.Host
	case "grpc":
		s.Net = "grpc"
		s.Path = transport.ServiceName
	case "mc1":
		s.Net = "mc1"
		s.Path = defaultString(transport.Path, "/")
		s.Host = defaultString(transport.Host, data.Server)
		s.Mode = defaultString(mc1Mode(transport), "auto")
		s.Cidr = strings.Join(mc1CidrSegments(transport), ",")
	case "httpupgrade":
		s.Net = "http"
		s.Path = transport.Path
		s.Host = transport.Host
	}

	if vmess.Security == "tls" {
		s.TLS = "tls"
		s.SNI = securityConfig.SNI
		s.AllowInsecure = securityConfig.AllowInsecure
		s.Fingerprint = securityConfig.Fingerprint
	}
	b, _ := json.Marshal(s)
	return "vmess://" + strings.TrimSuffix(base64.StdEncoding.EncodeToString(b), "=")
}

func VlessUri(data proxy.Proxy, uuid string) string {
	vless := data.Option.(proxy.Vless)
	transportConfig := vless.TransportConfig
	securityConfig := vless.SecurityConfig

	var query = make(url.Values)
	setQuery(&query, "flow", vless.Flow)
	setQuery(&query, "security", vless.Security)
	setQuery(&query, "encryption", "none")

	switch vless.Transport {
	case "websocket":
		setQuery(&query, "type", "ws")
		setQuery(&query, "host", transportConfig.Host)
		setQuery(&query, "path", transportConfig.Path)

	case "http2", "httpupgrade":
		setQuery(&query, "type", "http")
		setQuery(&query, "path", transportConfig.Path)
		setQuery(&query, "host", transportConfig.Host)
	case "grpc":
		setQuery(&query, "type", "grpc")
		setQuery(&query, "serviceName", transportConfig.ServiceName)
	case "mc1":
		setQuery(&query, "type", "mc1")
		setQuery(&query, "path", defaultString(transportConfig.Path, "/"))
		setQuery(&query, "host", defaultString(transportConfig.Host, data.Server))
		setQuery(&query, "mode", defaultString(mc1Mode(transportConfig), "auto"))
		if cidrSegments := mc1CidrSegments(transportConfig); len(cidrSegments) > 0 {
			setQuery(&query, "cidr", strings.Join(cidrSegments, ","))
		}
	}

	if vless.Security == "tls" {
		setQuery(&query, "sni", securityConfig.SNI)
		setQuery(&query, "fp", securityConfig.Fingerprint)
	} else if vless.Security == "reality" {
		setQuery(&query, "pbk", securityConfig.RealityPublicKey)
		setQuery(&query, "sid", securityConfig.RealityShortId)
		setQuery(&query, "sni", securityConfig.SNI)
		setQuery(&query, "fp", securityConfig.Fingerprint)
		setQuery(&query, "servername", securityConfig.SNI)
		setQuery(&query, "spx", "/")

	}

	u := url.URL{
		Scheme:   "vless",
		User:     url.User(uuid),
		Host:     net.JoinHostPort(data.Server, fmt.Sprint(data.Port)),
		RawQuery: query.Encode(),
		Fragment: data.Name,
	}
	return u.String()
}

func TrojanUri(data proxy.Proxy, uuid string) string {
	trojan := data.Option.(proxy.Trojan)
	transportConfig := trojan.TransportConfig
	securityConfig := trojan.SecurityConfig

	var query = make(url.Values)
	setQuery(&query, "security", trojan.Security)

	switch trojan.Transport {
	case "websocket":
		setQuery(&query, "type", "ws")
		setQuery(&query, "path", transportConfig.Path)
		setQuery(&query, "host", transportConfig.Host)
	case "grpc":
		setQuery(&query, "type", "grpc")
		setQuery(&query, "serviceName", transportConfig.ServiceName)
	case "mc1":
		setQuery(&query, "type", "mc1")
		setQuery(&query, "path", defaultString(transportConfig.Path, "/"))
		setQuery(&query, "host", defaultString(transportConfig.Host, data.Server))
		setQuery(&query, "mode", defaultString(mc1Mode(transportConfig), "auto"))
		if cidrSegments := mc1CidrSegments(transportConfig); len(cidrSegments) > 0 {
			setQuery(&query, "cidr", strings.Join(cidrSegments, ","))
		}
	default:
		setQuery(&query, "type", "tcp")
		setQuery(&query, "path", transportConfig.Path)
		setQuery(&query, "host", transportConfig.Host)
	}

	if securityConfig.AllowInsecure {
		setQuery(&query, "allowInsecure", "1")
	}

	u := &url.URL{
		Scheme:   "trojan",
		User:     url.User(uuid),
		Host:     net.JoinHostPort(data.Server, strconv.Itoa(data.Port)),
		RawQuery: query.Encode(),
		Fragment: data.Name,
	}
	return u.String()
}

func Hysteria2Uri(data proxy.Proxy, uuid string) string {
	hysteria2 := data.Option.(proxy.Hysteria2)

	var query = make(url.Values)

	setQuery(&query, "sni", hysteria2.SecurityConfig.SNI)

	if hysteria2.SecurityConfig.AllowInsecure {
		setQuery(&query, "insecure", "1")
	}

	if hp := strings.TrimSpace(hysteria2.HopPorts); hp != "" {
		setQuery(&query, "mport", hp)
	}

	if hysteria2.ObfsPassword != "" {
		setQuery(&query, "obfs", "salamander")
		setQuery(&query, "obfs-password", hysteria2.ObfsPassword)
	}

	u := &url.URL{
		Scheme:   "hysteria2",
		User:     url.User(uuid),
		Host:     net.JoinHostPort(data.Server, strconv.Itoa(data.Port)),
		RawQuery: query.Encode(),
		Fragment: data.Name,
	}
	return u.String()
}

func TuicUri(data proxy.Proxy, uuid string) string {
	tuic := data.Option.(proxy.Tuic)
	var query = make(url.Values)

	setQuery(&query, "congestion_control", tuic.CongestionController)
	setQuery(&query, "udp_relay_mode", tuic.UDPRelayMode)

	if tuic.SecurityConfig.SNI != "" {
		setQuery(&query, "sni", tuic.SecurityConfig.SNI)
	} else {
		setQuery(&query, "disable_sni", "1")
	}
	if tuic.SecurityConfig.AllowInsecure {
		setQuery(&query, "allow_insecure", "1")
	}

	u := &url.URL{
		Scheme:   "tuic",
		User:     url.User(uuid + ":" + uuid),
		Host:     net.JoinHostPort(data.Server, strconv.Itoa(data.Port)),
		RawQuery: query.Encode(),
		Fragment: data.Name,
	}
	return u.String()
}

func MxUri(data proxy.Proxy, uuid string) string {
	mx, ok := data.Option.(proxy.Mx)
	if !ok {
		return ""
	}
	transportConfig := mx.TransportConfig
	securityConfig := mx.SecurityConfig

	transport := strings.ToLower(strings.TrimSpace(mx.Transport))
	if transport == "" {
		transport = "tcp"
	}
	security := strings.ToLower(strings.TrimSpace(mx.Security))

	var query = make(url.Values)
	setQuery(&query, "mode", "multi")
	setQuery(&query, "security", security)

	switch transport {
	case "websocket", "ws":
		setQuery(&query, "type", "ws")
		setQuery(&query, "path", transportConfig.Path)
		setQuery(&query, "host", transportConfig.Host)
	case "grpc":
		setQuery(&query, "type", "grpc")
		setQuery(&query, "serviceName", transportConfig.ServiceName)
	case "mc1":
		setQuery(&query, "type", "mc1")
		setQuery(&query, "path", defaultString(transportConfig.Path, "/"))
		setQuery(&query, "host", defaultString(transportConfig.Host, data.Server))
		setQuery(&query, "mode", defaultString(mc1Mode(transportConfig), "auto"))
		if cidrSegments := mc1CidrSegments(transportConfig); len(cidrSegments) > 0 {
			setQuery(&query, "cidr", strings.Join(cidrSegments, ","))
		}
	case "xhttp":
		setQuery(&query, "type", "jatp")
		setQuery(&query, "endpoint_path", transportConfig.Path)
		setQuery(&query, "endpoint", defaultString(transportConfig.Host, data.Server))
		setQuery(&query, "split", transportConfig.Split)
	case "mundordp", "mundosql":
		setQuery(&query, "type", transport)
		applyMundoQuerySettings(&query, transportConfig)
	default:
		setQuery(&query, "type", transport)
	}

	switch security {
	case "tls":
		setQuery(&query, "fp", defaultString(securityConfig.Fingerprint, "chrome"))
		setQuery(&query, "sni", securityConfig.SNI)
		if securityConfig.AllowInsecure {
			setQuery(&query, "allowInsecure", "1")
		}
	case "reality":
		setQuery(&query, "fp", defaultString(securityConfig.Fingerprint, "chrome"))
		setQuery(&query, "sni", securityConfig.SNI)
		setQuery(&query, "servername", securityConfig.SNI)
		setQuery(&query, "pbk", securityConfig.RealityPublicKey)
		setQuery(&query, "sid", securityConfig.RealityShortId)
		setQuery(&query, "spx", "/")
	}

	u := &url.URL{
		Scheme:   "mx",
		User:     url.User(uuid),
		Host:     net.JoinHostPort(data.Server, strconv.Itoa(data.Port)),
		RawQuery: query.Encode(),
		Fragment: data.Name,
	}
	return u.String()
}

func setQuery(q *url.Values, k, v string) {
	if v != "" {
		q.Set(k, v)
	}
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func mc1Mode(config proxy.TransportConfig) string {
	if strings.TrimSpace(config.Mc1Mode) != "" {
		return config.Mc1Mode
	}
	return config.Mode
}

func mc1CidrSegments(config proxy.TransportConfig) []string {
	if len(config.Mc1CidrSegments) > 0 {
		return config.Mc1CidrSegments
	}
	return config.CidrSegments
}

func applyMundoQuerySettings(query *url.Values, config proxy.TransportConfig) {
	setQuery(query, "username", defaultString(config.Username, "MundoUser"))
	setQuery(query, "certificateFingerprint", config.CertificateFingerprint)
	setQuery(query, "fakeTitle", config.FakeTitle)
	setQuery(query, "fakeMessage", config.FakeMessage)
	if config.AcceptProxyProtocol {
		setQuery(query, "acceptProxyProtocol", "true")
	}
	if config.UseTLSCertificate {
		setQuery(query, "useTLSCertificate", "true")
	}
}
