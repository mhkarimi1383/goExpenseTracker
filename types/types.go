// all types are here
package types

// Configuration is used to store the configuration by cleanenv package
type Configuration struct {
	MetricAddress          string `env:"METRIC_ADDRESS" env-default:":9090" yaml:"metric_address"`
	APIAddress             string `env:"API_ADDRESS" env-default:":8080" yaml:"api_address"`
	SentryDsn              string `env:"SENTRY_DSN" env-default:"" yaml:"sentry_dsn"`
	LogFormat              string `env:"LOG_FORMAT" env-default:"text" yaml:"log_format"`
	ApplicationTitle       string `env:"APPLICATION_TITLE" env-default:"goExpenseTracker" yaml:"application_title"`
	ApplicationDescription string `env:"APPLICATION_DESCRIPTION" env-default:"goExpenseTracker" yaml:"application_description"`
	MongoDBConnectionURI   string `env:"MONGODB_CONNECTION_URI" env-default:"" yaml:"mongodb_connection_uri"`
	OpenIDAddress          string `env:"OPENID_ADDRESS" yaml:"openid_address"`
	OpenIDClientSecret     string `env:"OPENID_Client_SECRET" yaml:"openid_client_secret"`
	OpenIDClientID         string `env:"OPENID_Client_ID" yaml:"openid_client_id"`
	BaseURL                string `env:"BASE_URL" yaml:"base_url" env-default:"http://127.0.0.1:8080"`
}

type ApplicationInformation struct {
	Title       string
	Description string
}

type HealthzResponse struct {
	Name    string `json:"name"`
	Message string `json:"msg"`
}

// a type that used to create untyped map (for json)
type UntypedMap map[any]any

// any acceptable response will be here using generic it will accept one of the given types
type Response interface {
	HealthzResponse
}

type Item struct {
	Description string `bson:"description"`
	Operator    string `bson:"operator"`
	Amount      uint   `bson:"amount"`
	Id          uint   `bson:"_id"`
}

type IndexPage struct {
	Title  string `bson:"title"`
	Amount uint   `bson:"amount"`
	Items  []Item `bson:"items"`
}
