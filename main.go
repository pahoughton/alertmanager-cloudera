/* 2019-02-07 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"github.com/pahoughton/cloudera-amgr-alert/config"
	"github.com/pahoughton/cloudera-amgr-alert/cloudera"

	"gopkg.in/alecthomas/kingpin.v2"

)

var (
	Version	  string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion = runtime.Version()
)

type CommandArgs struct {
	ConfigFn	*string
	Debug		*bool
	AlertFn		*string
}

func main() {
	version := fmt.Sprintf(`%s: version %s branch: %s, rev: %s
  build: %s %s
`,
		path.Base(os.Args[0]),
		Version,
		Branch,
		Revision,
		BuildDate,
		GoVersion)

	app := kingpin.New(path.Base(os.Args[0]),
		"connect cloudera alerts to alertmanager").
			Version(version)

	args := CommandArgs{
		ConfigFn:	app.Flag("config-fn","config filename").
			Default("cloudera-amgr-alert.yml").String(),
		Debug:		app.Flag("debug","debug output to stdout").Bool(),
		AlertFn:	app.Arg("json", "json alert filename").String(),
	}
	debug := false
	if args.Debug != nil && *args.Debug {
		debug = true
	}
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if debug { fmt.Println("loading ",*args.ConfigFn) }

	cfg, err := config.Load(*args.ConfigFn)
	if err != nil {
		panic(err)
	}

	if err := cloudera.Send(*args.AlertFn,cfg,debug); err != nil {
		panic(err)
	}
	if debug { fmt.Println(*args.AlertFn) }
}
