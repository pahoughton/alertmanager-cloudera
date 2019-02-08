/* 2019-02-07 (cc) <paul4hough@gmail.com>
   process cloudera alerts json file
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	amgr "github.com/pahoughton/cloudera-amgr-alert/alertmanager"
	"github.com/pahoughton/cloudera-amgr-alert/config"
)

type AlertHeader struct {
	htype	string		`json:"type"`
	version	int			`json:"version"`
}

type AlertTime struct {
	IsoStr	string	`json:"iso8601"`
	Epoch	uint	`json:"epochMs"`
}

type Alert struct {
	Content	string						`json:"content"`
	When	AlertTime					`json:"timestamp"`
	Source	string						`json:"source"`
	Attrs	map[string][]interface{}	`json:"attributes"`
}

type AlertBody struct {
	Alert	Alert	`json:"alert"`
}

type AlertMsg struct {
	Body	AlertBody		`json:"body"`
	Header	AlertHeader		`json:"header"`
}

func parseCloudera(fn string,cfg *config.Config,debug bool) error {

	aFile, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer aFile.Close()

	var aData []AlertMsg
	//var aData []map[string]interface{}

	b, _ := ioutil.ReadAll(aFile)
	if err := json.Unmarshal(b, &aData); err != nil {
		return fmt.Errorf("json.Unmarshal %s: %s - %v",fn,err.Error(),string(b))
    }

	var amaList []amgr.AmgrAlert

	for _, a := range aData {

		attrs := a.Body.Alert.Attrs
		prevHealth := attrs["PREVIOUS_HEALTH_SUMMARY"][0].(string)
		suppressed := attrs["ALERT_SUPPRESSED"][0].(string)

		title := fmt.Sprintf("%s %s",
			attrs["CLUSTER_DISPLAY_NAME"][0].(string),
			attrs["ALERT_SUMMARY"][0].(string))

		if prevHealth != "GREEN" || suppressed != "false" {
			if debug {
				fmt.Printf("Skip: %s\n",title)
			}
			continue;
		}
		ama := amgr.AmgrAlert{
			StartsAt:		a.Body.Alert.When.IsoStr,
			GeneratorURL:	a.Body.Alert.Source,
		}
		ama.Labels = make(map[string]string)
		ama.Annots = make(map[string]string)

		if len(cfg.Global.Labels) > 0 {
			for lk, lv := range cfg.Global.Labels {
				ama.Labels[lk] = lv
			}
		}
		if len(cfg.Global.Annots) > 0 {
			for ak, av := range cfg.Global.Annots {
				ama.Annots[ak] = av
			}
		}
		ama.Labels["alertname"] = "cloudera-script"
		ama.Labels["uuid"]		= attrs["__uuid"][0].(string)
		ama.Annots["title"]		= title
		if instance, ok := attrs["HOSTS"]; ok {
			ama.Labels["instance"]	= instance[0].(string)
		}
		ama.Annots["description"]	= a.Body.Alert.Content

		amaList = append(amaList, ama)
	}
	if len(amaList) > 0 {
		if err := amgr.SendAlerts(amaList,cfg,debug); err != nil {
			return err
		}
	}
	return nil
}
