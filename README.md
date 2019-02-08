## cloudera-amgr-alert

[![Test Build Status](https://travis-ci.org/pahoughton/cloudera-amgr-alert.png)](https://travis-ci.org/pahoughton/cloudera-amgr-alert)

generate an
[alertmanager](https://prometheus.io/docs/alerting/alertmanager/)
alerts from
[cloudera](https://www.cloudera.com/documentation/enterprise/5-14-x/topics/cm_ag_alert_script.html#concept_sfx_lkw_yt)
alert json.  Example alert json available in
[cloudera-alert.json](../master/cloudera-alert.json)

## install

install script and binary in alert publisher user's path.

modify script configuration as needed.

## usage

configure as a cloudera alert script as described by the cloudera 5.14
[documentation](https://www.cloudera.com/documentation/enterprise/5-14-x/topics/cm_ag_alert_script.html#concept_sfx_lkw_yt)

## build

go build -mod=vendor

## contribute

https://github.com/pahoughton/cloudera-amgr-alert

## licenses

2019-01-16 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
