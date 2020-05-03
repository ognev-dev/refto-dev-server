package dataimport

type Data struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	Preview     string     `yaml:"preview"`
	Links       []DataLink `yaml:"links"`
	Topics      []string   `yaml:"topics"`
}

type DataLink struct {
	Icon string `yaml:"icon"`
	Text string `yaml:"text"`
	Addr string `yaml:"addr"`
}
