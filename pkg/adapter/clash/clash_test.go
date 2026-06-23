package clash

import (
	"embed"
	"testing"

	"github.com/npanel-dev/NPanel-backend/pkg/adapter/proxy"
	"github.com/stretchr/testify/assert"
)

//go:embed template/clash.tpl
var testTemplateFS embed.FS

func TestClash_Build(t *testing.T) {
	adapter := proxy.Adapter{
		TemplateFS: &testTemplateFS,
		Proxies: []proxy.Proxy{
			{
				Name:     "test-proxy",
				Protocol: "shadowsocks",
				Server:   "1.2.3.4",
				Port:     8388,
				Option: proxy.Shadowsocks{
					Method: "aes-256-gcm",
				},
			},
		},
		Group: []proxy.Group{
			{
				Name:    "test-group",
				Type:    "select",
				Proxies: []string{"test-proxy"},
			},
		},
		Rules: []string{
			"DOMAIN-SUFFIX,example.com,DIRECT",
			"GEOIP,CN,DIRECT",
			"MATCH,DIRECT",
		},
	}
	clash := NewClash(adapter)
	result, err := clash.Build("test-uuid")
	assert.NoError(t, err)
	assert.NotNil(t, result)

}
