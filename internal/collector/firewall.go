package collector

import (
	"fmt"
	"log"
	"time"
)

func (c *Collector) CollectFirewallMetrics() error {
	now := time.Now()
	since := now.Add(-24 * time.Hour)

	query := fmt.Sprintf(`{
		viewer {
			zones(filter: {zoneTag: "%s"}) {
				firewallEventsAdaptiveGroups(
					limit: 10000
					filter: {datetime_geq: "%s", datetime_leq: "%s"}
				) {
					count
					dimensions {
						action
						source
						ruleId
						clientRequestHTTPHost
						clientIP
						clientCountryName
						userAgent
					}
				}
			}
		}
	}`, c.zoneID, since.Format(time.RFC3339), now.Format(time.RFC3339))

	result, err := c.client.ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return c.processFirewallMetrics(result)
}

func (c *Collector) processFirewallMetrics(data map[string]interface{}) error {
	zones := data["data"].(map[string]interface{})["viewer"].(map[string]interface{})["zones"].([]interface{})
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	zone := zones[0].(map[string]interface{})
	groups, ok := zone["firewallEventsAdaptiveGroups"].([]interface{})
	if !ok || len(groups) == 0 {
		return fmt.Errorf("firewall metrics not available (may require Pro/Business plan)")
	}

	var totalEvents int64
	actionMap := make(map[string]int64)
	sourceMap := make(map[string]int64)
	ruleIDMap := make(map[string]int64)
	hostMap := make(map[string]int64)
	countryMap := make(map[string]int64)
	ipMap := make(map[string]int64)
	userAgentMap := make(map[string]int64)

	for _, g := range groups {
		group := g.(map[string]interface{})
		count := int64(group["count"].(float64))
		totalEvents += count

		dims, ok := group["dimensions"].(map[string]interface{})
		if !ok {
			continue
		}

		if action, ok := dims["action"].(string); ok && action != "" {
			actionMap[action] += count
		}
		if source, ok := dims["source"].(string); ok && source != "" {
			sourceMap[source] += count
		}
		if ruleID, ok := dims["ruleId"].(string); ok && ruleID != "" {
			ruleIDMap[ruleID] += count
		}
		if host, ok := dims["clientRequestHTTPHost"].(string); ok && host != "" {
			hostMap[host] += count
		}
		if country, ok := dims["clientCountryName"].(string); ok && country != "" {
			countryMap[country] += count
		}
		if ip, ok := dims["clientIP"].(string); ok && ip != "" {
			ipMap[ip] += count
		}
		if ua, ok := dims["userAgent"].(string); ok && ua != "" {
			userAgentMap[ua] += count
		}
	}

	c.metrics.FirewallEvents.WithLabelValues(c.zoneID).Set(float64(totalEvents))

	for action, count := range actionMap {
		c.metrics.FirewallAction.WithLabelValues(c.zoneID, action).Set(float64(count))
	}
	for source, count := range sourceMap {
		c.metrics.FirewallSource.WithLabelValues(c.zoneID, source).Set(float64(count))
	}
	for ruleID, count := range getTopN(ruleIDMap, 50) {
		c.metrics.FirewallRuleID.WithLabelValues(c.zoneID, ruleID).Set(float64(count))
	}
	for host, count := range getTopN(hostMap, 20) {
		c.metrics.FirewallHost.WithLabelValues(c.zoneID, host).Set(float64(count))
	}
	for country, count := range countryMap {
		c.metrics.FirewallCountry.WithLabelValues(c.zoneID, country).Set(float64(count))
	}
	for ip, count := range getTopN(ipMap, 100) {
		c.metrics.FirewallIP.WithLabelValues(c.zoneID, ip).Set(float64(count))
	}
	for ua, count := range getTopN(userAgentMap, 20) {
		c.metrics.FirewallUserAgent.WithLabelValues(c.zoneID, ua).Set(float64(count))
	}

	log.Printf(" Firewall: %d events | Actions:%d Sources:%d IPs:%d",
		totalEvents, len(actionMap), len(sourceMap), len(getTopN(ipMap, 100)))

	return nil
}

func getTopN(m map[string]int64, n int) map[string]int64 {
	type kv struct {
		Key   string
		Value int64
	}

	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}

	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].Value > pairs[i].Value {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	result := make(map[string]int64)
	limit := n
	if len(pairs) < limit {
		limit = len(pairs)
	}

	for i := 0; i < limit; i++ {
		result[pairs[i].Key] = pairs[i].Value
	}

	return result
}