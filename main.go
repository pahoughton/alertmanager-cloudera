/* 2019-02-07 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import (
	"fmt"
	"os"
	"path"

	"github.com/pahoughton/cloudera-amgr-alert/config"

	"gopkg.in/alecthomas/kingpin.v2"

)

type CommandArgs struct {
	ConfigFn	*string
	Debug		*bool
	AlertFn		*string
}

func main() {

	app := kingpin.New(path.Base(os.Args[0]),
		"cloudera alertmanager alert script").
			Version("1.0.1")

	args := CommandArgs{
		ConfigFn:	app.Flag("config-fn","config filename").
			Default("cloudera-amgr-alert.yml").ExistingFile(),
		Debug:		app.Flag("debug","debug output to stdout").
			Default("true").Bool(),
		AlertFn:	app.Arg("json", "json alert filename").
			Required().ExistingFile(),
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
	fmt.Println("loading ",*args.ConfigFn)

	cfg, err := config.LoadFile(*args.ConfigFn)
	if err != nil {
		panic(err)
	}

	if err := parseCloudera(*args.AlertFn,cfg,*args.Debug); err != nil {
		panic(err)
	}
	fmt.Println(*args.AlertFn)
}
