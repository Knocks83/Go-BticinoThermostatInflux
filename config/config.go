package config

// InfluxDB
const (
	InfluxHost            = "http://localhost:8086" // The host where the Influx DB is (default port: 8086)
	InfluxToken           = ""                      // The access token (if using 1.8>=version<2.0 use user:pass as token)
	InfluxOrg             = ""                      // The organization (if using 1.8>=version<2.0 leave empty)
	InfluxBucket          = ""                      // The bucket (if using 1.8>=version<2.0 use database/retention-policy, or just the db if the default rp should be used)
	InfluxMeasurementName = ""                      // The name of the Influx Measurement
)

// Thermostat APIs
const (
	ClientID        = ""
	ClientSecret    = ""
	Redirect        = ""
	SubscriptionKey = ""
	PlantID         = ""
	ModuleID        = ""
)

// Var
const (
	CalculateAbsolutePath = true // Whether the software should calculate the executable path (disable if using the `go run filename` syntax)
)

// DO NOT EDIT ANYTHING UNDER THIS COMMENT IF YOU DON'T KNOW WHAT YOU'RE DOING

// Thermostat
const (
	AuthEndpoint = "https://partners-login.eliotbylegrand.com/"       // Endpoint used to access the user account
	APIEndpoint  = "https://api.developer.legrand.com/smarther/v2.0/" // Endpoint used to send requests to the Thermostat
)

// Var
const (
	RefreshFileName = "refresh.txt" // The filename of the file where the RefreshToken'll be stored
	RequestDelay    = 182           // Delay in seconds between the requests (182 is fine-tuned for the free 500 requests/day, edit it ONLY if you have a custom plan)
)
