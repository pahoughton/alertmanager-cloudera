/* 2019-02-08 (cc) <paul4hough@gmail.com>
   prometheus alertmanager interface
*/
package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pahoughton/cloudera-amgr-alert/config"
)

const (
	amgrEndpoint	= "/api/v1/alerts"
	contTypeJSON	= "application/json"
)

type AmgrAlert struct {
	Labels			map[string]string	`json:"labels"`
	Annots			map[string]string	`json:"annotations"`
	StartsAt		string				`json:"startsAt"`
	GeneratorURL	string				`json:"generatorURL"`
}

func SendAlerts(amaList []AmgrAlert, cfg *config.Config,debug bool) error {

	if debug {
		djson, err := json.MarshalIndent(amaList,"","  ")
		if err != nil {
			return err
		}
		fmt.Println(string(djson))
	}
	amjson, err := json.Marshal(amaList)
	if err != nil {
		return err
	}

	for _, amgr := range cfg.Amgrs {
		for _, targ := range amgr.SConfigs.Targets {
			url := fmt.Sprintf("%s://%s%s",amgr.Scheme, targ, amgrEndpoint)

			if debug {
				fmt.Println("DEBUG: amgr url - ", url)
			} else {
				resp, err := http.Post(
					url,
					contTypeJSON,
					bytes.NewBuffer(amjson))
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					return fmt.Errorf("alertmanager status %s\n%s",
						resp.Status,
						string(amjson))
				} else {
					fmt.Printf("DEBUG: sent %s\n",string(amjson))
				}
			}

		}
	}

	return nil
}
