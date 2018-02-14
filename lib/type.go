package nifwall

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	nifcloud "github.com/tanishi/go-nifcloud"
)

type Client struct {
	C NifCloud
}

func NewClient(n NifCloud) *Client {
	return &Client{
		C: n,
	}
}

// NifCloud is interface for go-nifcloud mock
type NifCloud interface {
	AuthorizeSecurityGroupIngress(context.Context, *nifcloud.AuthorizeSecurityGroupIngressInput) (*nifcloud.AuthorizeSecurityGroupIngressOutput, error)
	CreateSecurityGroup(context.Context, *nifcloud.CreateSecurityGroupInput) (*nifcloud.CreateSecurityGroupOutput, error)
	DescribeSecurityGroups(context.Context, *nifcloud.DescribeSecurityGroupsInput) (*nifcloud.DescribeSecurityGroupsOutput, error)
	DescribeInstances(context.Context, *nifcloud.DescribeInstancesInput) (*nifcloud.DescribeInstancesOutput, error)
	DescribeInstanceAttribute(context.Context, *nifcloud.DescribeInstanceAttributeInput) (*nifcloud.DescribeInstanceAttributeOutput, error)
	DeleteSecurityGroup(context.Context, *nifcloud.DeleteSecurityGroupInput) (*nifcloud.DeleteSecurityGroupOutput, error)
	RegisterInstancesWithSecurityGroup(context.Context, *nifcloud.RegisterInstancesWithSecurityGroupInput) (*nifcloud.RegisterInstancesWithSecurityGroupOutput, error)
}

// FirewallGroup is struct for nifcloud API
type FirewallGroup struct {
	Name             string         `yaml:"name"`
	Description      string         `yaml:"description"`
	AvailabilityZone string         `yaml:"availability_zone"`
	IPPermissions    []ipPermission `yaml:"ip_permissions"`
}

type ipPermission struct {
	Protocol    string   `yaml:"protocol"`
	FromPort    string   `yaml:"from_port"`
	ToPort      string   `yaml:"to_port"`
	InOut       string   `yaml:"in_out"`
	GroupNames  []string `yaml:"group_names"`
	CidrIP      []string `yaml:"cidrip"`
	Description string   `yaml:"description"`
}

// NewFirewallGroup returns FirewallGroup with yaml file
func NewFirewallGroup(fpath string) (*FirewallGroup, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseYaml(file)
}

// ParseYaml returns FirewallGroup with yaml data
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
