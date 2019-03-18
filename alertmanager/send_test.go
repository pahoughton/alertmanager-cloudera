/* 2019-02-24 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package alertmanager

import (
	"net/http/httptest"
	"testing"
	"time"
	pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/alertmanager-cloudera/config"
)
func TestSend(t *testing.T) {
	//assert.True(t,true)
	am := &MockServer{}
	ms := httptest.NewServer(am)
	defer ms.Close()

	today, _ := time.Parse(time.RFC3339, "2019-02-20T13:12:11Z")

	cfg := &config.Config{
		Amgrs: []config.Amgr{
			config.Amgr{ Scheme: ms.URL,
				SConfigs: config.AmgrSConfig{
					Targets: []string{
						"",
					},
				},
			},
		},
	}

	amlist := make([]Alert,0,5)
	a := Alert{
		GeneratorURL:	"http://agate-nowhere/none",
		StartsAt:		today,
		Labels:			pmod.LabelSet{
			"alertname":	"cloudera-alert",
			"instance":		"localhost:9100",
			"job":			"cloudera",
			"mongrp":		"01-01",
			"component":	"part",
		}}

	amlist = append(amlist,a)
	amlen := len(amlist)
	err := Send(amlist,cfg,false)
	assert.Nil(t,err)
	assert.Equal(t,1,am.Hits)
	assert.Equal(t,amlen,am.Alerts)

	amlist = append(amlist,a,a)
	amlen += len(amlist)
	err = Send(amlist,cfg,false)
	assert.Nil(t,err)
	assert.Equal(t,2,am.Hits)
	assert.Equal(t,amlen,am.Alerts)
}
