package config

var (
	Port int

	EnableHTTPLogs    bool
	EnableMetricsPage bool

	Database struct {
		Engine string
		DSN    string
	}

	Migration struct {
		DownTo int64
	}

	Session struct {
		SecretKey string
		Expire    int
	}

	Debug   bool
	Version string
)
