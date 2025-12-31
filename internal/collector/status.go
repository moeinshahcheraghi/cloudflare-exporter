package collector

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func (c *Collector) CollectStatusMetrics() error {
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
					}
					dimensions {
						edgeResponseStatus
					}
				}
			}
		}
	}`, c.zoneID, since.Format("2006-01-02"), now.Add(24*time.Hour).Format("2006-01-02"))

	result, err := c.client.ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return c.processStatusMetrics(result)
}

func (c *Collector) processStatusMetrics(data map[string]interface{}) error {
	zones := data["data"].(map[string]interface{})["viewer"].(map[string]interface{})["zones"].([]interface{})
	if len(zones) == 0 {
		return fmt.Errorf("no zones found")
	}

	zone := zones[0].(map[string]interface{})
	groups := zone["httpRequests1dGroups"].([]interface{})

	statusMap := make(map[string]int64)
	var status2xxTotal, status3xxTotal, status4xxTotal, status5xxTotal int64

	for _, g := range groups {
		group := g.(map[string]interface{})
		sum := group["sum"].(map[string]interface{})
		reqs := int64(sum["requests"].(float64))

		if dims, ok := group["dimensions"].(map[string]interface{}); ok {
			if status, ok := dims["edgeResponseStatus"]; ok && status != nil {
				statusInt := int(status.(float64))
				statusStr := strconv.Itoa(statusInt)
				statusMap[statusStr] += reqs

				switch {
				case statusInt >= 200 && statusInt < 300:
					status2xxTotal += reqs
				case statusInt >= 300 && statusInt < 400:
					status3xxTotal += reqs
				case statusInt >= 400 && statusInt < 500:
					status4xxTotal += reqs
				case statusInt >= 500 && statusInt < 600:
					status5xxTotal += reqs
				}
			}
		}
	}

	c.metrics.Status2xx.WithLabelValues(c.zoneID).Set(float64(status2xxTotal))
	c.metrics.Status3xx.WithLabelValues(c.zoneID).Set(float64(status3xxTotal))
	c.metrics.Status4xx.WithLabelValues(c.zoneID).Set(float64(status4xxTotal))
	c.metrics.Status5xx.WithLabelValues(c.zoneID).Set(float64(status5xxTotal))

	for status, count := range statusMap {
		c.metrics.EdgeResponseStatus.WithLabelValues(c.zoneID, status).Set(float64(count))
	}

	log.Printf("✅ Status: %d codes | 2xx:%d 3xx:%d 4xx:%d 5xx:%d",
		len(statusMap), status2xxTotal, status3xxTotal, status4xxTotal, status5xxTotal)

	return nil
}