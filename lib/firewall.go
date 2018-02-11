package nifwall

import (
	"context"
	"sync"

	nifcloud "github.com/tanishi/go-nifcloud"
)

var client *nifcloud.Client

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

			res, _ := client.DescribeInstanceAttribute(ctx, param)

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
	res, err := client.DescribeInstances(ctx, &nifcloud.DescribeInstancesInput{})

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
