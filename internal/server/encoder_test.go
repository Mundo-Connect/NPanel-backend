package server

import (
	"net/http"
	"net/url"
	"testing"
)

func TestShouldEmitUnpopulatedFields(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "admin subscribe config", path: "/v1/admin/system/subscribe_config", want: true},
		{name: "admin subscribe config with api prefix", path: "/api/v1/admin/system/subscribe_config", want: true},
		{name: "public site config", path: "/v1/common/site/config", want: true},
		{name: "other route", path: "/v1/admin/system/site_config", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldEmitUnpopulatedFields(&http.Request{URL: &url.URL{Path: tt.path}})
			if got != tt.want {
				t.Fatalf("shouldEmitUnpopulatedFields(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
