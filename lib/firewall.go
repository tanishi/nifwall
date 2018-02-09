package nifwall

import (
	"context"

	nifcloud "github.com/tanishi/go-nifcloud"
)

var client *nifcloud.Client

// CreateSecurityGroup create firewall group
func CreateSecurityGroup(ctx context.Context, name, description string) error {
	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        name,
		GroupDescription: description,
	}

	_, err := client.CreateSecurityGroup(ctx, param)

	return err
}

// AddRuleToSecurityGroup add rule to firewall group
func AddRuleToSecurityGroup(ctx context.Context, name string, permissions []ipPermission) error {
	param := generateAuthorizeSecurityGroupIngressInput(name, permissions)

	_, err := client.AuthorizeSecurityGroupIngress(ctx, param)

	return err
}

// RegisterInstancesWithSecurityGroup apply firewall group to instance
func RegisterInstancesWithSecurityGroup(ctx context.Context, fwName, serverName string) error {
	param := &nifcloud.RegisterInstancesWithSecurityGroupInput{
		GroupName:   fwName,
		InstanceIDs: []string{serverName},
	}

	_, err := client.RegisterInstancesWithSecurityGroup(ctx, param)

	return err
}

func convert(permissions []ipPermission) []nifcloud.IPPermission {
	res := make([]nifcloud.IPPermission, 0, len(permissions))

	for _, p := range permissions {
		res = append(res, nifcloud.IPPermission{
			IPProtocol:  p.Protocol,
			FromPort:    p.FromPort,
			ToPort:      p.ToPort,
			InOut:       p.InOut,
			IPRanges:    p.CidrIP,
			Description: p.Description,
		})
	}

	return res
}

func generateAuthorizeSecurityGroupIngressInput(name string, permissions []ipPermission) *nifcloud.AuthorizeSecurityGroupIngressInput {
	return &nifcloud.AuthorizeSecurityGroupIngressInput{
		GroupName:     name,
		IPPermissions: convert(permissions),
	}
}
