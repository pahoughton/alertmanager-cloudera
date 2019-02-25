/* 2019-02-08 (cc) <paul4hough@gmail.com>
   prometheus alertmanager interface
*/
package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	pmod "github.com/prometheus/common/model"
	"github.com/pahoughton/cloudera-amgr-alert/config"
)

const (
	amgrEndpoint	= "/api/v1/alerts"
	contTypeJSON	= "application/json"
)
type Alert pmod.Alert

/*
type AmgrAlert struct {
	Labels			map[string]string	`json:"labels"`
	Annots			map[string]string	`json:"annotations"`
	StartsAt		string				`json:"startsAt"`
	GeneratorURL	string				`json:"generatorURL"`
}
*/

func Send(amlist []Alert,cfg *config.Config,debug bool) error {

	if debug {
		djson, err := json.MarshalIndent(amlist,"","  ")
		if err != nil {
			fmt.Printf("DEBUG marshal %v\n",err)
			return err
		}
		fmt.Println(string(djson))
	}
	amjson, err := json.Marshal(amlist)
	if err != nil {
		return err
	}

	for _, amgr := range cfg.Amgrs {
		for _, targ := range amgr.SConfigs.Targets {
			// for testing
			var url string
			if len(targ) < 1 {
				url = amgr.Scheme
			} else {
				url = fmt.Sprintf("%s://%s%s",amgr.Scheme, targ, amgrEndpoint)
			}
			if debug { fmt.Println("DEBUG: amgr url - ", url) }

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
			}
			if debug { fmt.Printf("DEBUG: sent %s\n",string(amjson)) }
		}
	}
	return nil
}
