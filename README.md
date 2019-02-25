## alertmanager-cloudera

[![Test Build Status](https://travis-ci.org/pahoughton/alertmanager-cloudera.png)](https://travis-ci.org/pahoughton/alertmanager-cloudera)

[alertmanager](https://prometheus.io/docs/alerting/alertmanager/)
alerts from
[cloudera](https://www.cloudera.com/documentation/enterprise/5-14-x/topics/cm_ag_alert_script.html#concept_sfx_lkw_yt)
alert script json input.

## usage

configure as a cloudera alert script as described by the cloudera 5.14
[documentation](https://www.cloudera.com/documentation/enterprise/5-14-x/topics/cm_ag_alert_script.html#concept_sfx_lkw_yt)

### config

config/tesetdata/good-ful.yml
[github](blob/master/config/testdata/good-full.yml)
[gitlab](../master/config/testdata/good-full.yml)

## build

go build -mod=vendor

## validate

go test ./...

## contribute

https://github.com/pahoughton/alertmanager-cloudera

## licenses

2019-01-16 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
