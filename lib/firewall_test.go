package nifwall

import (
	"context"
	"flag"
	"os"
	"reflect"
	"testing"

	nifcloud "github.com/tanishi/go-nifcloud"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		client = NewClient(&Mock{})
	} else {
		endpoint := os.Getenv("NIFCLOUD_ENDPOINT")
		accessKey := os.Getenv("NIFCLOUD_ACCESSKEY")
		secretAccessKey := os.Getenv("NIFCLOUD_SECRET_ACCESSKEY")

		c, _ := nifcloud.NewClient(endpoint, accessKey, secretAccessKey)
		client = NewClient(c)
	}

	os.Exit(m.Run())
}

func TestCreateSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName := "nifwallTest"

	defer teardownFirewallGroup(ctx, t, fwName)

	if err := client.CreateSecurityGroup(ctx, fwName, fwName, ""); err != nil {
		t.Error(err)
	}

	param := &nifcloud.DescribeSecurityGroupsInput{
		GroupNames: []string{fwName},
	}

	for {
		res, err := client.C.DescribeSecurityGroups(ctx, param)

		if err != nil {
			t.Errorf("Not Created")
		}

		status := res.SecurityGroupInfo[0].GroupStatus

		if status == "applied" {
			break
		}
	}
}

func TestAddRuleToSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestAddRuleToSecurityGroup(ctx, t)
	defer teardown(ctx, t, fwName)

	permissions := []ipPermission{
		{
			Protocol: "HTTP",
			CidrIP:   []string{"0.0.0.0/0"},
		},
	}

	if err := client.AddRuleToSecurityGroup(ctx, fwName, permissions); err != nil {
		t.Error(err)
	}

	param := &nifcloud.DescribeSecurityGroupsInput{
		GroupNames: []string{fwName},
	}

	for {
		res, err := client.C.DescribeSecurityGroups(ctx, param)

		if err != nil {
			t.Errorf("Not Created")
		}

		status := res.SecurityGroupInfo[0].GroupStatus
		resPermissions := res.SecurityGroupInfo[0].IPPermission

		if status == "applied" && len(resPermissions) > 0 {
			break
		}
	}
}

func TestRegisterInstancesWithSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestRegisterInstancesWithSecurityGroup(ctx, t)
	defer teardown(ctx, t, fwName)

	serverName := "tanishiTest"

	if err := client.RegisterInstancesWithSecurityGroup(ctx, fwName, serverName); err != nil {
		t.Error(err)
	}

	param := &nifcloud.DescribeSecurityGroupsInput{
		GroupNames: []string{fwName},
	}

	for {
		res, err := client.C.DescribeSecurityGroups(ctx, param)

		if err != nil {
			t.Errorf("Not Created")
		}

		status := res.SecurityGroupInfo[0].GroupStatus
		resInstances := res.SecurityGroupInfo[0].Instances

		if status == "applied" && len(resInstances) > 0 {
			break
		}
	}
}

func TestListInstances(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, err := client.ListInstances(ctx); err != nil {
		t.Error(err)
	}
}

func TestListInappropriateInstances(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestListInappropriateInstances(ctx, t)
	defer teardown(ctx, t, fwName)

	appropriateFWNames := []string{fwName}

	if _, err := client.ListInappropriateInstances(ctx, appropriateFWNames); err != nil {
		t.Error(err)
	}
}

func TestConvert(t *testing.T) {
	expected := []nifcloud.IPPermission{
		{
			IPProtocol:  "HTTP",
			FromPort:    "80",
			ToPort:      "81",
			InOut:       "IN",
			IPRanges:    []string{"0.0.0.0/0"},
			Description: "HOGE",
		},
	}
	actual := convert([]ipPermission{
		{
			Protocol:    "HTTP",
			FromPort:    "80",
			ToPort:      "81",
			InOut:       "IN",
			CidrIP:      []string{"0.0.0.0/0"},
			Description: "HOGE",
		},
	})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, but: %v", expected, actual)
	}
}

func TestGenerateAuthorizeSecurityGroupIngressInput(t *testing.T) {
	expected := &nifcloud.AuthorizeSecurityGroupIngressInput{
		GroupName: "HOGE",
		IPPermissions: []nifcloud.IPPermission{
			{
				IPProtocol: "HTTP",
				IPRanges:   []string{"0.0.0.0/0"},
			},
		},
	}

	actual := generateAuthorizeSecurityGroupIngressInput("HOGE", []ipPermission{
		{
			Protocol: "HTTP",
			CidrIP:   []string{"0.0.0.0/0"},
		},
	})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v,\nbut: %v", expected, actual)
	}
}

func setupTestAddRuleToSecurityGroup(ctx context.Context, t *testing.T) (string, func(context.Context, *testing.T, string)) {
	fwName := "nifwallRuleTest"

	setupFirewallGroup(ctx, t, fwName, "")

	return fwName, teardownFirewallGroup
}

func setupTestRegisterInstancesWithSecurityGroup(ctx context.Context, t *testing.T) (string, func(context.Context, *testing.T, string)) {
	fwName := "nifRegister"
	zone := "west-12"

	setupFirewallGroup(ctx, t, fwName, zone)

	return fwName, teardownFirewallGroup
}

func setupTestListInappropriateInstances(ctx context.Context, t *testing.T) (string, func(context.Context, *testing.T, string)) {
	fwName := "nifwall"

	setupFirewallGroup(ctx, t, fwName, "")

	return fwName, teardownFirewallGroup
}

func setupFirewallGroup(ctx context.Context, t *testing.T, fwName, zone string) {
	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        fwName,
		AvailabilityZone: zone,
	}

	client.C.CreateSecurityGroup(ctx, param)

	for {
		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		res, err := client.C.DescribeSecurityGroups(ctx, param)

		if err != nil {
			t.Errorf("Not Created")
		}

		status := res.SecurityGroupInfo[0].GroupStatus

		if status == "applied" {
			break
		}
	}

}

func teardownFirewallGroup(ctx context.Context, t *testing.T, fwName string) {
	param := &nifcloud.DeleteSecurityGroupInput{
		GroupName: fwName,
	}

	if _, err := client.C.DeleteSecurityGroup(ctx, param); err != nil {
		t.Fatal(err)
	}
}

type Mock struct {
	NifCloud
}

func (m *Mock) DescribeInstanceAttribute(ctx context.Context, param *nifcloud.DescribeInstanceAttributeInput) (*nifcloud.DescribeInstanceAttributeOutput, error) {
	return &nifcloud.DescribeInstanceAttributeOutput{
		GroupID: "groupID",
	}, nil
}

func (m *Mock) DescribeInstances(ctx context.Context, param *nifcloud.DescribeInstancesInput) (*nifcloud.DescribeInstancesOutput, error) {
	return &nifcloud.DescribeInstancesOutput{
		InstancesSet: []nifcloud.InstancesItem{
			{
				InstanceID: "test1",
			},
			{
				InstanceID: "test2",
			},
			{
				InstanceID: "test3",
			},
			{
				InstanceID: "test4",
			},
			{
				InstanceID: "test5",
			},
		},
	}, nil
}

func (m *Mock) DescribeSecurityGroups(ctx context.Context, param *nifcloud.DescribeSecurityGroupsInput) (*nifcloud.DescribeSecurityGroupsOutput, error) {
	return &nifcloud.DescribeSecurityGroupsOutput{
		SecurityGroupInfo: []nifcloud.SecurityGroupInfoItem{
			{
				GroupStatus: "applied",
				GroupName:   param.GroupNames[0],
				IPPermission: []nifcloud.IPPermissionItem{
					{
						IPProtocol: "HTTP",
					},
				},
				Instances: []nifcloud.InstanceItem{
					{
						InstanceID: "test",
					},
				},
			},
		},
	}, nil
}

func (m *Mock) CreateSecurityGroup(ctx context.Context, param *nifcloud.CreateSecurityGroupInput) (*nifcloud.CreateSecurityGroupOutput, error) {
	return &nifcloud.CreateSecurityGroupOutput{}, nil
}

func (m *Mock) DeleteSecurityGroup(ctx context.Context, param *nifcloud.DeleteSecurityGroupInput) (*nifcloud.DeleteSecurityGroupOutput, error) {
	return &nifcloud.DeleteSecurityGroupOutput{}, nil
}

func (m *Mock) AuthorizeSecurityGroupIngress(ctx context.Context, param *nifcloud.AuthorizeSecurityGroupIngressInput) (*nifcloud.AuthorizeSecurityGroupIngressOutput, error) {
	return &nifcloud.AuthorizeSecurityGroupIngressOutput{}, nil
}

func (m *Mock) RegisterInstancesWithSecurityGroup(ctx context.Context, param *nifcloud.RegisterInstancesWithSecurityGroupInput) (*nifcloud.RegisterInstancesWithSecurityGroupOutput, error) {
	return &nifcloud.RegisterInstancesWithSecurityGroupOutput{}, nil
}
