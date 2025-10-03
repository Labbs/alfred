package config

type Config struct {
	// Version of the application
	Version string

	// ConfigFile is the path to the configuration file
	ConfigFile string

	Server struct {
		// Port is the server port
		Port int
		// HttpLogs indicates if HTTP logs are enabled
		HttpLogs bool
	}

	// Logger is the configuration for the zerolog logger.
	// Level is the log level for the logger.
	// Pretty enables or disables pretty printing of logs (non JSON logs).
	Logger struct {
		Level  string
		Pretty bool
	}

	// Database is the configuration for the database connection.
	// Dialect is the database engine (sqlite, postgres, etc.).
	// DSN is the Data Source Name for the database connection.
	Database struct {
		Dialect string // Database engine (sqlite, postgres, etc.)
		DSN     string
	}
}
