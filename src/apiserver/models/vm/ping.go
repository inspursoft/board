package vm

type EmailPingParam struct {
	Hostname string
	Port     int
	Username string
	Password string
	IsTLS    bool
}
