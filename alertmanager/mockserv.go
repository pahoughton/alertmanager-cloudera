/* 2019-02-24 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package alertmanager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)
type MockServer struct {
	Hits	int
	Alerts	int
}

func (m *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Hits += 1
	if b, err := ioutil.ReadAll(r.Body);  err == nil {
		amlist := make([]Alert,0,10)
		if err = json.Unmarshal(b, &amlist); err == nil {
			m.Alerts += len(amlist)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
