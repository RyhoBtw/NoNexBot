package NoiseDcBot

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	Token         string `yaml:"token"`
	JoinMessage   string `yaml:"joinMessage"`
	JoinChannel   string `yaml:"joinChannel"`
	JoinRole      string `yaml:"joinRole"`
	StreamChannel string `yaml:"streamChannel"`
}

func ReadConf(filename string) (*Conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Conf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
