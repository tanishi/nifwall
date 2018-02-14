package nifwall

import (
	"context"
	"sync"

	nifcloud "github.com/tanishi/go-nifcloud"
)

type client struct {
	c *nifcloud.Client
}

// Client.c is for using nifcloud api
var Client client

// ListInappropriateInstances returns inappropriate instances name
func ListInappropriateInstances(ctx context.Context, fwName string) ([]string, error) {
	instanceNames, err := ListInstances(ctx)

	if err != nil {
		return nil, err
	}

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)
	result := make([]string, 0)

	for _, name := range instanceNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			param := &nifcloud.DescribeInstanceAttributeInput{
				InstanceID: name,
				Attribute:  "groupId",
			}

			res, _ := Client.c.DescribeInstanceAttribute(ctx, param)

			if res.GroupID != fwName {
				mutex.Lock()
				result = append(result, name)
				mutex.Unlock()
			}
		}(name)
	}

	wg.Wait()

	return result, nil
}

// ListInstances returns instances name
func ListInstances(ctx context.Context) ([]string, error) {
	res, err := Client.c.DescribeInstances(ctx, &nifcloud.DescribeInstancesInput{})

	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(res.InstancesSet))

	for _, instance := range res.InstancesSet {
		result = append(result, instance.InstanceID)
	}

	return result, nil
}

// CreateSecurityGroup create firewall group
func CreateSecurityGroup(ctx context.Context, name, description, zone string) error {
	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        name,
		GroupDescription: description,
		AvailabilityZone: zone,
	}

	_, err := Client.c.CreateSecurityGroup(ctx, param)

	return err
}

// AddRuleToSecurityGroup add rule to firewall group
func AddRuleToSecurityGroup(ctx context.Context, name string, permissions []ipPermission) error {
	param := generateAuthorizeSecurityGroupIngressInput(name, permissions)

	_, err := Client.c.AuthorizeSecurityGroupIngress(ctx, param)

	return err
}

// RegisterInstancesWithSecurityGroup apply firewall group to instance
func RegisterInstancesWithSecurityGroup(ctx context.Context, fwName, serverName string) error {
	param := &nifcloud.RegisterInstancesWithSecurityGroupInput{
		GroupName:   fwName,
		InstanceIDs: []string{serverName},
	}

	_, err := Client.c.RegisterInstancesWithSecurityGroup(ctx, param)

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

// UpdateFirewall create firewall with rule
func UpdateFirewall(ctx context.Context, fwPath string) error {
	fg, err := NewFirewallGroup(fwPath)

	if err != nil {
		return err
	}

	if err := CreateSecurityGroup(ctx, fg.Name, fg.Description, fg.AvailabilityZone); err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fg.Name},
		}

		for {
			res, _ := Client.c.DescribeSecurityGroups(ctx, param)

			status := res.SecurityGroupInfo[0].GroupStatus

			if status == "applied" {
				break
			}
		}

		wg.Done()
	}()

	wg.Wait()

	return AddRuleToSecurityGroup(ctx, fg.Name, fg.IPPermissions)
}
