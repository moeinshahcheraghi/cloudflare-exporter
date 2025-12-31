package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	TotalRequests     *prometheus.GaugeVec
	CachedRequests    *prometheus.GaugeVec
	UncachedRequests  *prometheus.GaugeVec
	EncryptedRequests *prometheus.GaugeVec
	PageViews         *prometheus.GaugeVec
	TotalBytes        *prometheus.GaugeVec
	CachedBytes       *prometheus.GaugeVec
	UncachedBytes     *prometheus.GaugeVec
	EncryptedBytes    *prometheus.GaugeVec
	Threats           *prometheus.GaugeVec
	CacheHitRate      *prometheus.GaugeVec
	EncryptionRate    *prometheus.GaugeVec

	ClientWaitTime *prometheus.GaugeVec
	AvgWaitTime    *prometheus.GaugeVec

	CountryRequests *prometheus.GaugeVec
	CountryBytes    *prometheus.GaugeVec

	EdgeResponseStatus *prometheus.GaugeVec
	Status2xx          *prometheus.GaugeVec
	Status3xx          *prometheus.GaugeVec
	Status4xx          *prometheus.GaugeVec
	Status5xx          *prometheus.GaugeVec

	ContentTypeRequests *prometheus.GaugeVec
	ContentTypeBytes    *prometheus.GaugeVec

	FirewallEvents    *prometheus.GaugeVec
	FirewallAction    *prometheus.GaugeVec
	FirewallSource    *prometheus.GaugeVec
	FirewallRuleID    *prometheus.GaugeVec
	FirewallHost      *prometheus.GaugeVec
	FirewallCountry   *prometheus.GaugeVec
	FirewallIP        *prometheus.GaugeVec
	FirewallUserAgent *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		TotalRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_requests_total",
				Help: "Total number of requests to the zone",
			},
			[]string{"zone_id"},
		),
		CachedRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_requests_cached",
				Help: "Number of cached requests",
			},
			[]string{"zone_id"},
		),
		UncachedRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_requests_uncached",
				Help: "Number of uncached requests",
			},
			[]string{"zone_id"},
		),
		EncryptedRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_requests_encrypted",
				Help: "Number of HTTPS requests",
			},
			[]string{"zone_id"},
		),
		PageViews: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_pageviews_total",
				Help: "Total page views",
			},
			[]string{"zone_id"},
		),
		TotalBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_bandwidth_total_bytes",
				Help: "Total bandwidth in bytes",
			},
			[]string{"zone_id"},
		),
		CachedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_bandwidth_cached_bytes",
				Help: "Cached bandwidth in bytes",
			},
			[]string{"zone_id"},
		),
		UncachedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_bandwidth_uncached_bytes",
				Help: "Uncached bandwidth in bytes",
			},
			[]string{"zone_id"},
		),
		EncryptedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_bandwidth_encrypted_bytes",
				Help: "Encrypted bandwidth in bytes",
			},
			[]string{"zone_id"},
		),
		Threats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_threats_total",
				Help: "Number of threats detected",
			},
			[]string{"zone_id"},
		),
		CacheHitRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_cache_hit_rate_percent",
				Help: "Cache hit rate percentage",
			},
			[]string{"zone_id"},
		),
		EncryptionRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_encryption_rate_percent",
				Help: "Encryption rate percentage",
			},
			[]string{"zone_id"},
		),
		ClientWaitTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_client_wait_time_total_ms",
				Help: "Total client wait time in milliseconds",
			},
			[]string{"zone_id"},
		),
		AvgWaitTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_client_wait_time_avg_ms",
				Help: "Average wait time per request in milliseconds",
			},
			[]string{"zone_id"},
		),
		CountryRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_requests_country",
				Help: "Number of requests by country",
			},
			[]string{"zone_id", "country"},
		),
		CountryBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_bandwidth_country_bytes",
				Help: "Bandwidth by country in bytes",
			},
			[]string{"zone_id", "country"},
		),
		EdgeResponseStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_edge_response_status",
				Help: "Number of requests by HTTP status code",
			},
			[]string{"zone_id", "status"},
		),
		Status2xx: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_status_2xx_total",
				Help: "Total number of 2xx success responses",
			},
			[]string{"zone_id"},
		),
		Status3xx: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_status_3xx_total",
				Help: "Total number of 3xx redirect responses",
			},
			[]string{"zone_id"},
		),
		Status4xx: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_status_4xx_total",
				Help: "Total number of 4xx client error responses",
			},
			[]string{"zone_id"},
		),
		Status5xx: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_status_5xx_total",
				Help: "Total number of 5xx server error responses",
			},
			[]string{"zone_id"},
		),
		ContentTypeRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_requests_content_type",
				Help: "Number of requests by content type",
			},
			[]string{"zone_id", "content_type"},
		),
		ContentTypeBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_bandwidth_content_type_bytes",
				Help: "Bandwidth by content type in bytes",
			},
			[]string{"zone_id", "content_type"},
		),
		FirewallEvents: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_events_total",
				Help: "Total number of firewall events",
			},
			[]string{"zone_id"},
		),
		FirewallAction: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_action",
				Help: "Number of firewall events by action",
			},
			[]string{"zone_id", "action"},
		),
		FirewallSource: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_source",
				Help: "Number of firewall events by source",
			},
			[]string{"zone_id", "source"},
		),
		FirewallRuleID: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_rule_id",
				Help: "Number of firewall events by rule ID",
			},
			[]string{"zone_id", "rule_id"},
		),
		FirewallHost: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_host",
				Help: "Number of firewall events by attacked host",
			},
			[]string{"zone_id", "host"},
		),
		FirewallCountry: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_country",
				Help: "Number of firewall events by attacker country",
			},
			[]string{"zone_id", "country"},
		),
		FirewallIP: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_ip",
				Help: "Number of firewall events by attacker IP (top 100)",
			},
			[]string{"zone_id", "ip"},
		),
		FirewallUserAgent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cloudflare_zone_firewall_user_agent",
				Help: "Number of firewall events by user agent",
			},
			[]string{"zone_id", "user_agent"},
		),
	}
}

func (m *Metrics) Register() {
	prometheus.MustRegister(
		m.TotalRequests,
		m.CachedRequests,
		m.UncachedRequests,
		m.EncryptedRequests,
		m.PageViews,
		m.TotalBytes,
		m.CachedBytes,
		m.UncachedBytes,
		m.EncryptedBytes,
		m.Threats,
		m.CacheHitRate,
		m.EncryptionRate,
		m.ClientWaitTime,
		m.AvgWaitTime,
		m.CountryRequests,
		m.CountryBytes,
		m.EdgeResponseStatus,
		m.Status2xx,
		m.Status3xx,
		m.Status4xx,
		m.Status5xx,
		m.ContentTypeRequests,
		m.ContentTypeBytes,
		m.FirewallEvents,
		m.FirewallAction,
		m.FirewallSource,
		m.FirewallRuleID,
		m.FirewallHost,
		m.FirewallCountry,
		m.FirewallIP,
		m.FirewallUserAgent,
	)
}