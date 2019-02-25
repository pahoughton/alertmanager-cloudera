/* 2019-02-07 (cc) <paul4hough@gmail.com>
   process cloudera alerts json file
*/
package cloudera

import (
	"encoding/json"
	"fmt"
	"time"

	pmod "github.com/prometheus/common/model"
	amgr "github.com/pahoughton/cloudera-amgr-alert/alertmanager"
	"github.com/pahoughton/cloudera-amgr-alert/config"
)

type AlertHeader struct {
	htype	string		`json:"type"`
	version	int			`json:"version"`
}

type AlertTime struct {
	Iso		time.Time	`json:"iso8601"`
	Epoch	uint		`json:"epochMs"`
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

func parse(dat []byte,cfg *config.Config,debug bool) []amgr.Alert {
	var cloudera []AlertMsg
	//var aData []map[string]interface{}

	if err := json.Unmarshal(dat, &cloudera); err != nil {
		panic(err)
    }

	amaList := make([]amgr.Alert,0,len(cloudera))

	for _, a := range cloudera {

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
		ama := amgr.Alert{
			StartsAt:		a.Body.Alert.When.Iso,
			GeneratorURL:	a.Body.Alert.Source,
		}
		ama.Labels = make(pmod.LabelSet)
		ama.Annotations = make(pmod.LabelSet)

		if len(cfg.Global.Labels) > 0 {
			for lk, lv := range cfg.Global.Labels {
				ama.Labels[lk] = lv
			}
		}
		if len(cfg.Global.Annots) > 0 {
			for ak, av := range cfg.Global.Annots {
				ama.Annotations[ak] = av
			}
		}
		if _, ok := ama.Labels["alertname"]; ! ok {
			ama.Labels["alertname"] = "cloudera-script"
		}
		ama.Labels["uuid"]		= pmod.LabelValue(attrs["__uuid"][0].(string))
		ama.Annotations["title"] = pmod.LabelValue(title)
		if instance, ok := attrs["HOSTS"]; ok {
			ama.Labels["instance"]	= pmod.LabelValue(instance[0].(string))
		}
		ama.Annotations["description"] = pmod.LabelValue(a.Body.Alert.Content)

		amaList = append(amaList, ama)
	}
	return amaList
}
