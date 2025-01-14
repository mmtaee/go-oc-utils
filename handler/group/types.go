package ocservGroup

// OcGroupConfig ocserv group config
type OcGroupConfig struct {
	RxDataPerSec         *string   `json:"rx-data-per-sec"`
	TxDataPerSec         *string   `json:"tx-data-per-sec"`
	MaxSameClients       *int      `json:"max-same-clients"`
	IPv4Network          *string   `json:"ipv4-network"`
	DNS                  *[]string `json:"dns"`
	NoUDP                *bool     `json:"no-udp"`
	KeepAlive            *int      `json:"keepalive"`
	DPD                  *int      `json:"dpd"`
	MobileDPD            *int      `json:"mobile-dpd"`
	TunnelAllDNS         *bool     `json:"tunnel-all-dns"`
	RestrictUserToRoutes *bool     `json:"restrict-user-to-routes"`
	StatsReportTime      *int      `json:"stats-report-time"`
	MTU                  *int      `json:"mtu"`
	IdleTimeout          *int      `json:"idle-timeout"`
	MobileIdleTimeout    *int      `json:"mobile-idle-timeout"`
	SessionTimeout       *int      `json:"session-timeout"`
}

// OcGroupConfigInfo ocserv group info with config
type OcGroupConfigInfo struct {
	Name   string         `json:"name"`
	Path   string         `json:"-"`
	Config *OcGroupConfig `json:"config"`
}
