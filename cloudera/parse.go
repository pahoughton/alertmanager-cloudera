/* 2019-02-07 (cc) <paul4hough@gmail.com>
   process cloudera alerts json file
*/
package cloudera

import (
	"encoding/json"
	"fmt"
	"time"
	//"strings"
	"net/url"

	pmod "github.com/prometheus/common/model"
	amgr "github.com/pahoughton/alertmanager-cloudera/alertmanager"
	"github.com/pahoughton/alertmanager-cloudera/config"
)

//Top Level (level 0)
type AlertMsg struct {
	AlertMain []interface{}
	Body	AlertBody		`json:"body"`
	Header	AlertHeader		`json:"header"`
}
//Header struct (level 1)
type AlertHeader struct {
	htype	string		`json:"type"`
	version	int			`json:"version"`
}
//Body struct (level 1)
type AlertBody struct {
	Alert	Alert	`json:"alert"`
}
//Alert sctuct (level 2, under Body)
type Alert struct {
	Attrs	map[string][]interface{}	`json:"attributes"`
	Source	string		`json:"source"`
	Content	string		`json:"content"`
	When	AlertTime	`json:"timestamp"`
}

/*
type Attributes struct {
	AlertSuppressed []string `json:"ALERT_SUPPRESSED"`
	HealthTestName []string `json:"HEALTH_TEST_NAME"`
	PrevHealth []string `json:"PREVIOUS_HEALTH_SUMMARY"`
	ClusterName []string `json:"CLUSTER_DISPLAY_NAME"`
	AlertSumm []string `json:"ALERT_SUMMARY"`
	UUID []string `json:"__uuid"`
}
*/

//timestamp struct (level 3, under alert)
type AlertTime struct {
	Iso		time.Time	`json:"iso8601"`
	Epoch	uint		`json:"epochMs"`
}

func parse(dat []byte,cfg *config.Config,debug bool) []amgr.Alert {
	//defines cloudera with above struct
	var cloudera []AlertMsg
	//unmarshals passed in []byte to above defined cloudera variable
	if err := json.Unmarshal(dat, &cloudera); err != nil {
		panic(err)
	}
	//defines the alert
	amaList := make([]amgr.Alert,0,len(cloudera))
	//for the length of the cloudera array
	for _, a := range cloudera {

		attrs := a.Body.Alert.Attrs
		prevHealth := "GREEN"
		if len(attrs["PREVIOUS_HEALTH_SUMMARY"]) > 0 {
		//if len(a.Body.Alert.Attrs.PrevHealth) > 0 {
			prevHealth = attrs["PREVIOUS_HEALTH_SUMMARY"][0].(string)	
			//prevHealth = a.Body.Alert.Attrs.PrevHealth.(string)
		}
		//prevHealth := attrs["PREVIOUS_HEALTH_SUMMARY"][0].(string)
		suppressed := attrs["ALERT_SUPPRESSED"][0].(string)
		//suppressed := a.Body.Alert.Attrs.AlertSuppressed.(string)

		title := fmt.Sprintf("%s %s",
		attrs["CLUSTER_DISPLAY_NAME"][0].(string),
		//a.Body.Alert.Attrs.ClusterName.(string),
		attrs["ALERT_SUMMARY"][0].(string))
		//a.Body.Alert.Attrs.AlertSumm.(string))

	if prevHealth != "GREEN" || suppressed != "false" {
		if debug {
			fmt.Printf("Skip: %s\n",title)
		}
		continue;
	}
/*
		//this section takes attributes in the JSON and makes them strings
		prevHealthString := strings.Join(a.Body.Alert.Attrs.PrevHealth,"")
		clusterNameString := strings.Join(a.Body.Alert.Attrs.ClusterName,"")
		alertSummaryString := strings.Join(a.Body.Alert.Attrs.AlertSumm,"") 
		isSuppressed := strings.Join(a.Body.Alert.Attrs.AlertSuppressed,"")
		uuid := strings.Join(a.Body.Alert.Attrs.UUID,"")
		//defines the title by combining two strings
		title := fmt.Sprintf("%s %s",clusterNameString,alertSummaryString)
		//checks if the alert is NOT green or NOT suppressed.
		if prevHealthString != "GREEN" || isSuppressed != "false" {
			if debug {
				fmt.Printf("Skip: %v : %v\n", clusterNameString, alertSummaryString)
				}
			continue;
		}
*/
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
		//this section parses out the Hostname from the URL found in the Source
		//this is because the HOST field is found so infrequently in the JSON
		s, err := url.Parse(a.Body.Alert.Source)
		if err != nil {
			panic(err)
		}
		ama.Labels["instance"]	= pmod.LabelValue(s.Host)
		ama.Annotations["description"] = pmod.LabelValue(a.Body.Alert.Content)
		
		amaList = append(amaList, ama)
		
	}
	return amaList
}

/*
	var cloudera []AlertMsg
	//var aData []map[string]interface{}

	if err := json.Unmarshal(dat, &cloudera); err != nil {
		panic(err)
    }

	amaList := make([]amgr.Alert,0,len(cloudera))

	for _, a := range cloudera {

		attrs := a.Body.Alert.Attrs
		prevHealth := "GREEN"
		if len(attrs["PREVIOUS_HEALTH_SUMMARY"]) > 0 {
		  prevHealth = attrs["PREVIOUS_HEALTH_SUMMARY"][0].(string)	
		}
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
*/

