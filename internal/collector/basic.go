package collector

import (
	"fmt"
	"log"
	"time"

	"cloudflare-exporter/internal/metrics"
	"cloudflare-exporter/pkg/cloudflare"
)

type Collector struct {
	client  *cloudflare.Client
	metrics *metrics.Metrics
	zoneID  string
}

func NewCollector(client *cloudflare.Client, metrics *metrics.Metrics, zoneID string) *Collector {
	return &Collector{
		client:  client,
		metrics: metrics,
		zoneID:  zoneID,
	}
}

func (c *Collector) CollectAll() error {
	if err := c.CollectBasicMetrics(); err != nil {
		log.Printf(" Basic metrics: %v", err)
	}

	if err := c.CollectStatusMetrics(); err != nil {
		log.Printf(" Status metrics: %v", err)
	}

	if err := c.CollectContentTypeMetrics(); err != nil {
		log.Printf(" Content type metrics: %v", err)
	}

	if err := c.CollectFirewallMetrics(); err != nil {
		log.Printf("ℹ️  Firewall metrics: %v", err)
	}

	return nil
}

func (c *Collector) CollectBasicMetrics() error {
	now := time.Now()
	since := now.Add(-24 * time.Hour)

	query := fmt.Sprintf(`{
		viewer {
			zones(filter: {zoneTag: "%s"}) {
				httpRequests1dGroups(
					limit: 100
					filter: {date_geq: "%s", date_lt: "%s"}
				) {
					sum {
						requests
						cachedRequests
						bytes
						cachedBytes
						encryptedBytes
						encryptedRequests
						pageViews
						threats
						countryMap {
							clientCountryName
							requests
							bytes
						}
					}
				}
			}
		}
	}`, c.zoneID, since.Format("2006-01-02"), now.Add(24*time.Hour).Format("2006-01-02"))

	result, err := c.client.ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return c.processBasicMetrics(result)
}

func (c *Collector) processBasicMetrics(data map[string]interface{}) error {
	zones := data["data"].(map[string]interface{})["viewer"].(map[string]interface{})["zones"].([]interface{})
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	zone := zones[0].(map[string]interface{})
	groups := zone["httpRequests1dGroups"].([]interface{})

	var totalReqs, cachedReqs, encryptedReqs, pageViews int64
	var totalBw, cachedBw, encryptedBw int64
	var totalThreats int64

	countryReqMap := make(map[string]int64)
	countryBwMap := make(map[string]int64)

	for _, g := range groups {
		group := g.(map[string]interface{})
		sum := group["sum"].(map[string]interface{})

		totalReqs += int64(sum["requests"].(float64))
		cachedReqs += int64(sum["cachedRequests"].(float64))
		encryptedReqs += int64(sum["encryptedRequests"].(float64))
		pageViews += int64(sum["pageViews"].(float64))
		totalBw += int64(sum["bytes"].(float64))
		cachedBw += int64(sum["cachedBytes"].(float64))
		encryptedBw += int64(sum["encryptedBytes"].(float64))
		totalThreats += int64(sum["threats"].(float64))

		if countryMap, ok := sum["countryMap"].([]interface{}); ok {
			for _, c := range countryMap {
				country := c.(map[string]interface{})
				name := country["clientCountryName"].(string)
				if name != "" {
					countryReqMap[name] += int64(country["requests"].(float64))
					countryBwMap[name] += int64(country["bytes"].(float64))
				}
			}
		}
	}

	cacheHitRate := float64(0)
	if totalReqs > 0 {
		cacheHitRate = float64(cachedReqs) / float64(totalReqs) * 100
	}

	encryptionRate := float64(0)
	if totalReqs > 0 {
		encryptionRate = float64(encryptedReqs) / float64(totalReqs) * 100
	}

	c.metrics.TotalRequests.WithLabelValues(c.zoneID).Set(float64(totalReqs))
	c.metrics.CachedRequests.WithLabelValues(c.zoneID).Set(float64(cachedReqs))
	c.metrics.UncachedRequests.WithLabelValues(c.zoneID).Set(float64(totalReqs - cachedReqs))
	c.metrics.EncryptedRequests.WithLabelValues(c.zoneID).Set(float64(encryptedReqs))
	c.metrics.PageViews.WithLabelValues(c.zoneID).Set(float64(pageViews))
	c.metrics.TotalBytes.WithLabelValues(c.zoneID).Set(float64(totalBw))
	c.metrics.CachedBytes.WithLabelValues(c.zoneID).Set(float64(cachedBw))
	c.metrics.UncachedBytes.WithLabelValues(c.zoneID).Set(float64(totalBw - cachedBw))
	c.metrics.EncryptedBytes.WithLabelValues(c.zoneID).Set(float64(encryptedBw))
	c.metrics.Threats.WithLabelValues(c.zoneID).Set(float64(totalThreats))
	c.metrics.CacheHitRate.WithLabelValues(c.zoneID).Set(cacheHitRate)
	c.metrics.EncryptionRate.WithLabelValues(c.zoneID).Set(encryptionRate)

	for country, reqs := range countryReqMap {
		c.metrics.CountryRequests.WithLabelValues(c.zoneID, country).Set(float64(reqs))
	}
	for country, bw := range countryBwMap {
		c.metrics.CountryBytes.WithLabelValues(c.zoneID, country).Set(float64(bw))
	}

	log.Printf(" HTTP: %d reqs |  %.1f%% cache |  %.1f%% https |  %.0f MB |  %d countries",
		totalReqs, cacheHitRate, encryptionRate, float64(totalBw)/1024/1024, len(countryReqMap))

	return nil
}
