/* 2019-02-08 (cc) <paul4hough@gmail.com>
   prometheus alertmanager interface
*/
package alertmanager

import (
	"encoding/json"
	"fmt"

	"github.com/pahoughton/cloudera-amgr-alert/config"
)

type AmgrAlert struct {
	Labels			map[string]string	`json:"labels"`
	Annots			map[string]string	`json:"annotations"`
	StartsAt		string				`json:"startsAt"`
	GeneratorURL	string				`json:"generatorURL"`
}

func SendAlerts(amaList []AmgrAlert, cfg *config.Config,debug bool) error {

	amjson, err := json.MarshalIndent(amaList,"","  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n",string(amjson))

	return fmt.Errorf("alertmanager.sendAlert error")
}
