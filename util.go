package nifwall

import (
	"io"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type FirewallGroup struct {
	Name             string         `yaml:"name"`
	Description      string         `yaml:"description"`
	AvailabilityZone string         `yaml:"availability_zone"`
	IPPermissions    []IPPermission `yaml:"ip_permissions"`
}

type IPPermission struct {
	Protocol    string   `yaml:"protocol"`
	FromPort    int      `yaml:"from_port"`
	ToPort      int      `yaml:"to_port"`
	InOut       string   `yaml:"in_out"`
	GroupNames  []string `yaml:"group_names"`
	CidrIP      []string `yaml:"cidrip"`
	Description string   `yaml:"description"`
}

func ParseYaml(r io.Reader) (*FirewallGroup, error) {
	res, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	firewall := &FirewallGroup{}

	err = yaml.Unmarshal([]byte(res), firewall)
	if err != nil {
		return nil, err
	}

	return firewall, nil
}
