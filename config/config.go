package config

// InfluxDB
const (
	InfluxHost            = "http://localhost:8086" // The host where the Influx DB is (default port: 8086)
	InfluxToken           = ""                      // The access token (if using 1.8>=version<2.0 use user:pass as token)
	InfluxOrg             = ""                      // The organization (if using 1.8>=version<2.0 leave empty)
	InfluxBucket          = ""                      // The bucket (if using 1.8>=version<2.0 use database/retention-policy, or just the db the default rp should be used)
	InfluxMeasurementName = ""                      // The name of the Influx Measurement
)

// Thermostat APIs
const (
	ClientID        = ""
	ClientSecret    = ""
	SubscriptionKey = ""
	PlantID         = ""
	ModuleID        = ""
)

// DO NOT EDIT ANYTHING UNDER THIS COMMENT IF YOU DON'T KNOW
// WHAT YOU'RE DOING

// Thermostat
const (
	AuthEndpoint = "https://partners-login.eliotbylegrand.com/"       // Endpoint used to access the user account
	APIEndpoint  = "https://api.developer.legrand.com/smarther/v2.0/" // Endpoint used to send requests to the Thermostat
)

// Var
const (
	RefreshFileName = "refresh.txt"
)
