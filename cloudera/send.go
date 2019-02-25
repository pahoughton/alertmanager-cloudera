/* 2019-02-24 (cc) <paul4hough@gmail.com>
   Send cloudera alerts to alertmanager
*/
package cloudera

import (
	"os"
	"io/ioutil"
	"github.com/pahoughton/cloudera-amgr-alert/config"
	amgr "github.com/pahoughton/cloudera-amgr-alert/alertmanager"
)

func Send(fn string,cfg *config.Config,debug bool) error {
	if f, err := os.Open(fn); err == nil {
		defer f.Close()
		if b, err := ioutil.ReadAll(f); err == nil {
			if amlist := parse(b,cfg,debug); len(amlist) > 0 {
				if err := amgr.Send(amlist,cfg,debug); err != nil {
					return err
				}
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	return nil
}
