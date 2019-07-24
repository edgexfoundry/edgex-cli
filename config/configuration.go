package config

// Configuration struct will use this to write config file eventually
type Configuration struct {
	Host string
	Security
	Ports
}

// Ports struct for each service's port
type Ports struct {
	CoreData           string
	CoreMetadata       string
	CoreCommand        string
	Notifications      string
	Logging            string
	Scheduling         string
	RulesEngine        string
	ClientRegistration string
	SystemManagement   string
}

// Security struct for security related config
type Security struct {
	Enabled bool
	Token   string
}
