/* 2019-02-08 (cc) <paul4hough@gmail.com>
   parse config file
*/
package config

import (
	"io/ioutil"

	pmod "github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

type AmgrSConfig struct {
	Targets	[]string	`yaml:"targets"`
}

type Amgr struct {
	Scheme		string		`yaml:"scheme"`
	SConfigs	AmgrSConfig	`yaml:"static_configs"`
}

type GlobalConfig struct {
	Labels		pmod.LabelSet	`yaml:"labels,omitempty"`
	Annots		pmod.LabelSet	`yaml:"annotations,omitempty"`
}

type Config struct {
	Global		GlobalConfig		`yaml:"global"`
	Amgrs		[]Amgr				`yaml:"alertmanagers"`
}

func Load(fn string) (*Config, error) {

	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(dat, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
