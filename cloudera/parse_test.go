/* 2019-02-24 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package cloudera

import (
	"io/ioutil"
	"os"
	"testing"

	pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/cloudera-amgr-alert/config"
)
func TestParse(t *testing.T) {
	cfg := &config.Config{}

	f, err := os.Open("testdata/cloudera-alert.json")
	assert.Nil(t,err)
	assert.NotNil(t,f)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	assert.Nil(t,err)
	assert.NotNil(t,b)

	expuuid := []pmod.LabelValue{
		"89521139-0859-4bef-bf65-eb141e63dbba",
		"67b4d1c4-791b-428e-a9ea-8a09d4885f5d",
	}
	got := parse(b,cfg,false)
	assert.Equal(t,2,len(got))
	assert.Equal(t,got[0].Labels["uuid"],expuuid[0])
	assert.Equal(t,got[1].Labels["uuid"],expuuid[1])
}
