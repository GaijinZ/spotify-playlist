package config

type GlobalEnv struct {
	ClientID     string `envconfig:"client_id"`
	ClientSecret string `envconfig:"client_secret"`
	RedirectURI  string `envconfig:"redirect_uri"`
	AuthURL      string `envconfig:"auth_url"`
	TokenURL     string `envconfig:"token_url"`
	Scope        string `envconfig:"scope"`
	Port         string `envconfig:"port"`
	Host         string `envconfig:"host"`
	RedisHost    string `envconfig:"redis_host"`
	RedisPort    string `envconfig:"redis_port"`
	ClusterIP    string `envconfig:"cluster_ip"`
	KeySpace     string `envconfig:"key_space"`
	BaseHost     string `envconfig:"base_host"`
	SecretKey    string `envconfig:"secret_key"`
	AutoSplitVar string `split_words:"true"`
}
