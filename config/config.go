package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	ElasticsearchURL  string
	Check             string
	NodeIP            string
	NodeName          string
	WarningThreshold  int
	CriticalThreshold int
}

func LoadConfig() (*Config, error) {
	flag.String("es_url", "", "Elasticsearch URL")
	flag.String("check", "", "Check to perform")
	flag.String("node_ip", "", "Node IP address for filtering")
	flag.String("node_name", "", "Node Name for filtering")
	flag.Int("w", 0, "Warning threshold")
	flag.Int("c", 0, "Critical threshold")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("es_url"); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("check"); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("node_ip"); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("node_name"); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("w"); err != nil {
		return nil, err
	}

	if err := viper.BindEnv("c"); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	config := &Config{
		ElasticsearchURL:  viper.GetString("es_url"),
		Check:             viper.GetString("check"),
		NodeIP:            viper.GetString("node_ip"),
		NodeName:          viper.GetString("node_name"),
		WarningThreshold:  viper.GetInt("w"),
		CriticalThreshold: viper.GetInt("c"),
	}

	return config, nil
}
