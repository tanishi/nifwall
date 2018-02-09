package nifwall

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewFirewallGroup(t *testing.T) {
	t.Run("It should return FirewallGroup", func(t *testing.T) {
		expected := &FirewallGroup{
			Name:             "FirewallGroupName",
			Description:      "FirewallGroupDescription",
			AvailabilityZone: "",
			IPPermissions: []ipPermission{
				ipPermission{
					Protocol:    "HTTP",
					FromPort:    "80",
					ToPort:      "81",
					InOut:       "IN",
					GroupNames:  []string{"hoge"},
					CidrIP:      []string{"0.0.0.0/0"},
					Description: "memomemo",
				},
			},
		}

		fpath := "./examples/example.yml"

		actual, err := NewFirewallGroup(fpath)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v,\nbut: %v", expected, actual)
		}
	})

	t.Run("It should return error", func(t *testing.T) {
		_, err := NewFirewallGroup("")

		if err == nil {
			t.Errorf("NewFirewallGroup should return error")
		}

		_, err = ParseYaml(bytes.NewBufferString("hoge"))

		if err == nil {
			t.Errorf("ParseYaml return error")
		}
	})
}

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
		IPPermissions: []ipPermission{
			ipPermission{
				Protocol:    "HTTP",
				FromPort:    "80",
				ToPort:      "81",
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
