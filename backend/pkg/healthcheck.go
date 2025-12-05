package pkg

const (
	ServiceHealthy   ServiceHealthState = "healthy"
	ServiceUnhealthy ServiceHealthState = "unhealthy"
	ServiceUnknown   ServiceHealthState = "unknown"
)

type ServiceHealthState string

type ServiceHealthError struct {
	Service string `json:"service"`
	Error   string `json:"error"`
}

type ServiceHealth interface {
	Name() string
	Check() error
}
