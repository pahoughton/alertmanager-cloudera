/* 2019-02-24 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package config

import (
	"testing"
	pmod "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)
func TestLoad(t *testing.T) {
	assert.True(t,true)
	got, err := Load("testdata/good-full.yml")
	assert.Nil(t,err)
	assert.NotNil(t,got)
	assert.Equal(t,got.Global.Labels["env"],pmod.LabelValue("sandbox"))
	assert.Equal(t,got.Amgrs[0].SConfigs.Targets[1],"1.2.3.5:9093")
}
