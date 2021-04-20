package types

type Config struct {
	Backend *DatabaseMetadata `yaml:"database_metadata"`
	Slack   *Slack            `yaml:"slack"`
}
