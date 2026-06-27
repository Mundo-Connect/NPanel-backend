package routing

import (
	"sort"
	"strings"
)

func ApplyClientCapabilities(envelope Envelope, supportedFeatures []string) Envelope {
	missing := MissingRequiredFeatures(envelope.CapabilityRequirements.RequiredFeatures, supportedFeatures)
	if len(missing) == 0 {
		return envelope
	}

	missingSet := featureSet(missing)
	gated := envelope
	gated.Mode = ModeObserve
	gated.DNSResolvers = filterSupportedDNSResolvers(envelope.DNSResolvers, missingSet)
	gated.Outbounds = filterSupportedOutbounds(envelope.Outbounds, missingSet)

	resolverTags := resolverTagSet(gated.DNSResolvers)
	outboundTags := outboundTagSet(gated.Outbounds)
	gated.Profile.DefaultDNSResolverTag = supportedDNSResolverTag(gated.Profile.DefaultDNSResolverTag, resolverTags)
	gated.Profile.DefaultAction = supportedRouteAction(gated.Profile.DefaultAction, resolverTags, outboundTags)
	gated.Profile.DefaultFallbackPolicy = supportedFailPolicy(gated.Profile.DefaultFallbackPolicy, missingSet)
	gated.UnlockServices = supportedUnlockServices(envelope.UnlockServices, resolverTags, outboundTags, missingSet)
	gated.Rules = supportedRules(envelope.Rules, resolverTags, outboundTags)
	gated.HealthSnapshot = filterHealthSnapshot(envelope.HealthSnapshot, resolverTags, outboundTags, serviceCodeSet(gated.UnlockServices))
	gated.RoutingHash = StableHash(gated)
	return gated
}

func featureSet(features []string) map[string]struct{} {
	set := make(map[string]struct{}, len(features))
	for _, feature := range features {
		feature = strings.TrimSpace(feature)
		if feature != "" {
			set[feature] = struct{}{}
		}
	}
	return set
}

func featureMissing(missing map[string]struct{}, feature string) bool {
	_, ok := missing[feature]
	return ok
}

func filterSupportedDNSResolvers(items []DNSResolver, missing map[string]struct{}) []DNSResolver {
	result := make([]DNSResolver, 0, len(items))
	for _, item := range items {
		if dnsResolverBlocked(item, missing) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func dnsResolverBlocked(item DNSResolver, missing map[string]struct{}) bool {
	if featureMissing(missing, "route_dns_resolver") {
		return true
	}
	proto := strings.ToLower(strings.TrimSpace(item.Proto))
	address := strings.ToLower(strings.TrimSpace(item.Address))
	if (proto == "doh" || proto == "https" || strings.HasPrefix(address, "https://")) && featureMissing(missing, "doh") {
		return true
	}
	if (proto == "dot" || proto == "tls" || strings.HasPrefix(address, "tls://")) && featureMissing(missing, "dot") {
		return true
	}
	if proto == "udp" && featureMissing(missing, "dns_udp") {
		return true
	}
	if proto == "tcp" && featureMissing(missing, "dns_tcp") {
		return true
	}
	return false
}

func filterSupportedOutbounds(items []RouteOutbound, missing map[string]struct{}) []RouteOutbound {
	result := make([]RouteOutbound, 0, len(items))
	for _, item := range items {
		if outboundBlocked(item, missing) {
			continue
		}
		item.FailPolicy = supportedFailPolicy(item.FailPolicy, missing)
		if featureMissing(missing, "route_fallback") {
			item.FallbackPoolTags = nil
		}
		result = append(result, item)
	}
	return result
}

func outboundBlocked(item RouteOutbound, missing map[string]struct{}) bool {
	if featureMissing(missing, "route_outbound") {
		return true
	}
	if strings.ToLower(strings.TrimSpace(item.Type)) != "external" && item.External == nil {
		return false
	}
	protocol := ""
	if item.External != nil {
		protocol = strings.ToLower(strings.TrimSpace(item.External.Protocol))
	}
	switch protocol {
	case "wireguard", "wg":
		return featureMissing(missing, "external_wireguard")
	case "socks", "socks4", "socks5":
		return featureMissing(missing, "external_socks")
	case "http", "https":
		return featureMissing(missing, "external_http")
	default:
		return featureMissing(missing, "external_wireguard") ||
			featureMissing(missing, "external_socks") ||
			featureMissing(missing, "external_http")
	}
}

func supportedRouteAction(action RouteAction, resolverTags, outboundTags map[string]struct{}) RouteAction {
	if action.Type == "outbound" {
		if _, ok := outboundTags[action.OutboundTag]; !ok {
			return RouteAction{Type: "proxy"}
		}
	}
	if action.Type == "dns_resolver" {
		if _, ok := resolverTags[action.DNSResolverTag]; !ok {
			return RouteAction{Type: "proxy"}
		}
	}
	if action.DNSResolverTag != "" {
		action.DNSResolverTag = supportedDNSResolverTag(action.DNSResolverTag, resolverTags)
	}
	return action
}

func supportedDNSResolverTag(tag string, resolverTags map[string]struct{}) string {
	tag = strings.TrimSpace(tag)
	if tag == "" || tag == "dns:system" {
		return tag
	}
	if _, ok := resolverTags[tag]; ok {
		return tag
	}
	return "dns:system"
}

func supportedFailPolicy(policy string, missing map[string]struct{}) string {
	if featureMissing(missing, "route_fail_policy") || featureMissing(missing, "route_fallback") {
		return ""
	}
	return policy
}

func supportedUnlockServices(items []UnlockService, resolverTags, outboundTags map[string]struct{}, missing map[string]struct{}) []UnlockService {
	result := make([]UnlockService, 0, len(items))
	for _, item := range items {
		item.DefaultDNSResolverTag = supportedDNSResolverTag(item.DefaultDNSResolverTag, resolverTags)
		if item.DefaultOutboundTag != "" && item.DefaultOutboundTag != "proxy:default" {
			if _, ok := outboundTags[item.DefaultOutboundTag]; !ok {
				item.DefaultOutboundTag = "proxy:default"
			}
		}
		item.DefaultFailPolicy = supportedFailPolicy(item.DefaultFailPolicy, missing)
		result = append(result, item)
	}
	return result
}

func supportedRules(items []Rule, resolverTags, outboundTags map[string]struct{}) []Rule {
	result := make([]Rule, 0, len(items))
	for _, item := range items {
		action := supportedRouteAction(item.Action, resolverTags, outboundTags)
		if action.Type != item.Action.Type && item.Action.Type != "" {
			continue
		}
		item.Action = action
		result = append(result, item)
	}
	return result
}

func resolverTagSet(items []DNSResolver) map[string]struct{} {
	result := map[string]struct{}{
		"":           {},
		"dns:system": {},
	}
	for _, item := range items {
		result[item.Tag] = struct{}{}
	}
	return result
}

func outboundTagSet(items []RouteOutbound) map[string]struct{} {
	result := map[string]struct{}{
		"":              {},
		"proxy:default": {},
	}
	for _, item := range items {
		result[item.Tag] = struct{}{}
	}
	return result
}

func serviceCodeSet(items []UnlockService) map[string]struct{} {
	result := map[string]struct{}{}
	for _, item := range items {
		result[item.Code] = struct{}{}
	}
	return result
}

func filterHealthSnapshot(snapshot HealthSnapshot, resolverTags, outboundTags, serviceCodes map[string]struct{}) HealthSnapshot {
	filtered := snapshot
	filtered.Outbounds = nil
	for _, item := range snapshot.Outbounds {
		if _, ok := outboundTags[item.Tag]; ok && item.Tag != "" {
			filtered.Outbounds = append(filtered.Outbounds, item)
		}
	}
	filtered.DNSResolvers = nil
	for _, item := range snapshot.DNSResolvers {
		if _, ok := resolverTags[item.Tag]; ok && item.Tag != "" {
			filtered.DNSResolvers = append(filtered.DNSResolvers, item)
		}
	}
	filtered.Services = nil
	for _, item := range snapshot.Services {
		if _, ok := serviceCodes[item.Code]; ok {
			filtered.Services = append(filtered.Services, item)
		}
	}
	sort.SliceStable(filtered.Outbounds, func(i, j int) bool { return filtered.Outbounds[i].Tag < filtered.Outbounds[j].Tag })
	sort.SliceStable(filtered.DNSResolvers, func(i, j int) bool { return filtered.DNSResolvers[i].Tag < filtered.DNSResolvers[j].Tag })
	sort.SliceStable(filtered.Services, func(i, j int) bool { return filtered.Services[i].Code < filtered.Services[j].Code })
	return filtered
}
