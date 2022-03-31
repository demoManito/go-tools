package sync

// DataBusConfig redis config.
type DataBusConfig struct {
	Addr string `xml:"addr" yaml:"addr" json:"addr"`
	Pass string `xml:"auth" yaml:"auth" json:"auth"`
	DB   int    `xml:"db" yaml:"db" json:"db"`
	Key  string `xml:"key" yaml:"key" json:"key"`

	Secret   string `xml:"secret" yaml:"secret" json:"secret"`
	Group    string `xml:"group" yaml:"group" json:"group"`
	Topic    string `xml:"topic" yaml:"topic" json:"topic"`
	Action   string `xml:"action" yaml:"action" json:"action"` // shoule be "pub" or "sub" or "pubsub"
	Buffer   int    `xml:"buffer" yaml:"buffer" json:"buffer"`
	PoolSize int    `xml:"pool_size" yaml:"pool_size" json:"pool_size"` // pub pool size, default: 10
}

// mockConfig local redis server
func mockConfig() *DataBusConfig {
	return &DataBusConfig{
		Addr:     "127.0.0.1:6379",
		Pass:     "", // auth password set
		DB:       0,  // use default DB
		PoolSize: 1,
	}
}
