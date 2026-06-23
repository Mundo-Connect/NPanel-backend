mode: rule
proxies:
{{.Proxies | toYaml | indent 2}}
proxy-groups:
{{.ProxyGroups | toYaml | indent 2}}
rules:
{{.Rules | toYaml | indent 2}}
