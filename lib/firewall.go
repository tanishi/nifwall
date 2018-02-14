package nifwall

import (
	"context"
	"sync"

	nifcloud "github.com/tanishi/go-nifcloud"
)

// Client.C is for using nifcloud api
var Client client

// ListInappropriateInstances returns inappropriate instances name
func (c *client) ListInappropriateInstances(ctx context.Context, fwName string) ([]string, error) {
	instanceNames, err := c.ListInstances(ctx)

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

			res, _ := c.C.DescribeInstanceAttribute(ctx, param)

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
func (c *client) ListInstances(ctx context.Context) ([]string, error) {
	res, err := c.C.DescribeInstances(ctx, &nifcloud.DescribeInstancesInput{})

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
func (c *client) CreateSecurityGroup(ctx context.Context, name, description, zone string) error {
	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        name,
		GroupDescription: description,
		AvailabilityZone: zone,
	}

	_, err := c.C.CreateSecurityGroup(ctx, param)

	return err
}

// AddRuleToSecurityGroup add rule to firewall group
func (c *client) AddRuleToSecurityGroup(ctx context.Context, name string, permissions []ipPermission) error {
	param := generateAuthorizeSecurityGroupIngressInput(name, permissions)

	_, err := c.C.AuthorizeSecurityGroupIngress(ctx, param)

	return err
}

// RegisterInstancesWithSecurityGroup apply firewall group to instance
func (c *client) RegisterInstancesWithSecurityGroup(ctx context.Context, fwName, serverName string) error {
	param := &nifcloud.RegisterInstancesWithSecurityGroupInput{
		GroupName:   fwName,
		InstanceIDs: []string{serverName},
	}

	_, err := c.C.RegisterInstancesWithSecurityGroup(ctx, param)

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
func (c *client) UpdateFirewall(ctx context.Context, fwPath string) error {
	fg, err := NewFirewallGroup(fwPath)

	if err != nil {
		return err
	}

	if err := c.CreateSecurityGroup(ctx, fg.Name, fg.Description, fg.AvailabilityZone); err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fg.Name},
		}

		for {
			res, _ := c.C.DescribeSecurityGroups(ctx, param)

			status := res.SecurityGroupInfo[0].GroupStatus

			if status == "applied" {
				break
			}
		}

		wg.Done()
	}()

	wg.Wait()

	return c.AddRuleToSecurityGroup(ctx, fg.Name, fg.IPPermissions)
}
