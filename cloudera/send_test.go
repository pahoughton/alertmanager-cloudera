/* 2019-02-24 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package cloudera

import (
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	amgr "github.com/pahoughton/alertmanager-cloudera/alertmanager"
	"github.com/pahoughton/alertmanager-cloudera/config"
)
func TestSend(t *testing.T) {
	//assert.True(t,true)
	am := &amgr.MockServer{}
	ms := httptest.NewServer(am)
	defer ms.Close()

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
	err := Send("testdata/cloudera-alert.json",cfg,false)
	assert.Nil(t,err)
	assert.Equal(t,1,am.Hits)
	assert.Equal(t,2,am.Alerts)
}
