package collector

import (
	"fmt"
	"log"
	"time"
)

func (c *Collector) CollectContentTypeMetrics() error {
	now := time.Now()
	since := now.Add(-24 * time.Hour)

	query := fmt.Sprintf(`{
		viewer {
			zones(filter: {zoneTag: "%s"}) {
				httpRequests1dGroups(
					limit: 1000
					filter: {date_geq: "%s", date_lt: "%s"}
				) {
					sum {
						requests
						bytes
					}
					dimensions {
						edgeResponseContentTypeName
					}
				}
			}
		}
	}`, c.zoneID, since.Format("2006-01-02"), now.Add(24*time.Hour).Format("2006-01-02"))

	result, err := c.client.ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return c.processContentTypeMetrics(result)
}

func (c *Collector) processContentTypeMetrics(data map[string]interface{}) error {
	zones := data["data"].(map[string]interface{})["viewer"].(map[string]interface{})["zones"].([]interface{})
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	zone := zones[0].(map[string]interface{})
	groups := zone["httpRequests1dGroups"].([]interface{})

	contentTypeReqMap := make(map[string]int64)
	contentTypeBwMap := make(map[string]int64)

	for _, g := range groups {
		group := g.(map[string]interface{})
		sum := group["sum"].(map[string]interface{})
		reqs := int64(sum["requests"].(float64))
		bw := int64(sum["bytes"].(float64))

		if dims, ok := group["dimensions"].(map[string]interface{}); ok {
			if ct, ok := dims["edgeResponseContentTypeName"]; ok && ct != nil {
				ctStr := ct.(string)
				if ctStr != "" {
					contentTypeReqMap[ctStr] += reqs
					contentTypeBwMap[ctStr] += bw
				}
			}
		}
	}

	for ct, reqs := range contentTypeReqMap {
		c.metrics.ContentTypeRequests.WithLabelValues(c.zoneID, ct).Set(float64(reqs))
	}
	for ct, bw := range contentTypeBwMap {
		c.metrics.ContentTypeBytes.WithLabelValues(c.zoneID, ct).Set(float64(bw))
	}

	log.Printf(" ContentType: %d types", len(contentTypeReqMap))

	return nil
}
