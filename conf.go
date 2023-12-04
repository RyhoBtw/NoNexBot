package NoiseDcBot

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var config *Configurations

type BotConf struct {
	Token         string `yaml:"token"`
	JoinMessage   string `yaml:"joinMessage"`
	JoinChannel   string `yaml:"joinChannel"`
	JoinRole      string `yaml:"joinRole"`
	StreamChannel string `yaml:"streamChannel"`
}

type Configurations struct {
	Host       string `yaml:"host"`
	DBName     string `yaml:"dbName"`
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
}

func ReadBotConf(filename string) (*BotConf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &BotConf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}

func ReadDBConf(filename string) (*Configurations, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Configurations{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	log.Println(c.Host, c.DBName, c.DBPassword, c.DBUser)
	config = c
	return c, err
}

func GetConf() *Configurations {
	return config
}
