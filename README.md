# nifwall

CLI tool for nifcloud firewall api

[![Build Status](https://travis-ci.org/tanishi/nifwall.svg?branch=master)](https://travis-ci.org/tanishi/nifwall)
[![Coverage Status](https://coveralls.io/repos/github/tanishi/nifwall/badge.svg?branch=master)](https://coveralls.io/github/tanishi/nifwall?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/tanishi/nifwall)](https://goreportcard.com/report/github.com/tanishi/nifwall)

list command display instances for which the specified FW has not been applied.

`nifwall list -fw firewall-name`

Also, you can specify more than one.
In that case returns instances that are not applied either.

`nifwall list -fw firewall-name -fw another-firewall-name`

update command create firewall from yml file

`nifwall update -f /path/to/firewall.yml`

apply command apply the specified FW to the instances

`nifwall apply -fw firewall-name instance-name ...`
