package nifwall

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseYaml(t *testing.T) {
	data :=
		`name: HOGE
description: HOGEHOGE
availability_zone: west-11
ip_permissions:
  - protocol: HTTP
    from_port: 80
    to_port: 81
    in_out: IN
    group_names:
      - "HTTPAllow"
    cidrip:
      - "0.0.0.0/0"
    description: "HTTP"`

	expected := &FirewallGroup{
		Name:             "HOGE",
		Description:      "HOGEHOGE",
		AvailabilityZone: "west-11",
		IPPermissions: []IPPermission{
			IPPermission{
				Protocol:    "HTTP",
				FromPort:    80,
				ToPort:      81,
				InOut:       "IN",
				GroupNames:  []string{"HTTPAllow"},
				CidrIP:      []string{"0.0.0.0/0"},
				Description: "HTTP",
			},
		},
	}
	actual, err := ParseYaml(bytes.NewBufferString(data))

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v,\nbut: %v", expected, actual)
	}
}
