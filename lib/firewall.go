package nifwall

import (
	"context"
	"sync"

	nifcloud "github.com/tanishi/go-nifcloud"
)

// Client is for using nifcloud api
var client *Client

// ListInappropriateInstances returns inappropriate instances name
func (c *Client) ListInappropriateInstances(ctx context.Context, fwNames []string) ([]string, error) {
	list := make(chan []string, 1)

	go func() {
		l, _ := c.ListInstances(ctx)
		list <- l
	}()

	ichan := make(chan string, 1)

	go func() {
		defer close(ichan)
		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: fwNames,
		}

		res, _ := c.C.DescribeSecurityGroups(ctx, param)

		for _, info := range res.SecurityGroupInfo {
			for _, instance := range info.Instances {
				ichan <- instance.InstanceID
			}
		}
	}()

	l := <-list
	for i := range ichan {
		l = func(list []string, target string) []string {
			r := make([]string, 0)
			for _, name := range l {
				if name != target {
					r = append(r, name)
				}
			}
			return r
		}(l, i)
	}

	return l, nil
}

// ListInstances returns instances name
func (c *Client) ListInstances(ctx context.Context) ([]string, error) {
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
func (c *Client) CreateSecurityGroup(ctx context.Context, name, description, zone string) error {
	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        name,
		GroupDescription: description,
		AvailabilityZone: zone,
	}

	_, err := c.C.CreateSecurityGroup(ctx, param)

	return err
}

// AddRuleToSecurityGroup add rule to firewall group
func (c *Client) AddRuleToSecurityGroup(ctx context.Context, name string, permissions []ipPermission) error {
	param := generateAuthorizeSecurityGroupIngressInput(name, permissions)

	_, err := c.C.AuthorizeSecurityGroupIngress(ctx, param)

	return err
}

// RegisterInstancesWithSecurityGroup apply firewall group to instance
func (c *Client) RegisterInstancesWithSecurityGroup(ctx context.Context, fwName, serverName string) error {
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
func (c *Client) UpdateFirewall(ctx context.Context, fwPath string) error {
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
