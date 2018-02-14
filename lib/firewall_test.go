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
	endpoint := os.Getenv("NIFCLOUD_ENDPOINT")
	accessKey := os.Getenv("NIFCLOUD_ACCESSKEY")
	secretAccessKey := os.Getenv("NIFCLOUD_SECRET_ACCESSKEY")

	if testing.Short() {
		Client.C = &mock{}
	} else {
		Client.C, _ = nifcloud.NewClient(endpoint, accessKey, secretAccessKey)
	}

	os.Exit(m.Run())
}

func TestCreateSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestCreateSecurityGroup(t)
	defer teardown(ctx, t)

	if err := Client.CreateSecurityGroup(ctx, fwName, fwName, ""); err != nil {
		t.Error(err)
	}

	done := make(chan *nifcloud.DescribeSecurityGroupsOutput, 0)

	go func() {
		defer close(done)
		for {
			param := &nifcloud.DescribeSecurityGroupsInput{
				GroupNames: []string{fwName},
			}

			res, err := Client.C.DescribeSecurityGroups(ctx, param)

			if err != nil {
				t.Errorf("Not Created")
			}

			status := res.SecurityGroupInfo[0].GroupStatus

			if status == "applied" {
				break
			}
		}
	}()

	<-done
}

func TestAddRuleToSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestAddRuleToSecurityGroup(ctx, t)
	defer teardown(ctx, t)

	permissions := []ipPermission{
		{
			Protocol: "HTTP",
			CidrIP:   []string{"0.0.0.0/0"},
		},
	}

	if err := Client.AddRuleToSecurityGroup(ctx, fwName, permissions); err != nil {
		t.Error(err)
	}

	done := make(chan *nifcloud.DescribeSecurityGroupsOutput, 0)

	go func() {
		defer close(done)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		for {
			res, err := Client.C.DescribeSecurityGroups(ctx, param)

			if err != nil {
				t.Errorf("Not Created")
			}

			status := res.SecurityGroupInfo[0].GroupStatus
			resPermissions := res.SecurityGroupInfo[0].IPPermission

			if status == "applied" && len(resPermissions) > 0 {
				break
			}
		}
	}()

	<-done
}

func TestRegisterInstancesWithSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestRegisterInstancesWithSecurityGroup(ctx, t)
	defer teardown(ctx, t)

	serverName := "tanishiTest"

	if err := Client.RegisterInstancesWithSecurityGroup(ctx, fwName, serverName); err != nil {
		t.Error(err)
	}

	done := make(chan *nifcloud.DescribeSecurityGroupsOutput, 0)

	go func() {
		defer close(done)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		for {
			res, err := Client.C.DescribeSecurityGroups(ctx, param)

			if err != nil {
				t.Errorf("Not Created")
			}

			status := res.SecurityGroupInfo[0].GroupStatus
			resInstances := res.SecurityGroupInfo[0].Instances

			if status == "applied" && len(resInstances) > 0 {
				break
			}
		}
	}()

	<-done
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

func setupTestCreateSecurityGroup(t *testing.T) (string, func(context.Context, *testing.T)) {
	fwName := "nifwallTest"

	return fwName, func(ctx context.Context, t *testing.T) {
		param := &nifcloud.DeleteSecurityGroupInput{
			GroupName: fwName,
		}

		if _, err := Client.C.DeleteSecurityGroup(ctx, param); err != nil {
			t.Fatal(err)
		}
	}
}

func setupTestAddRuleToSecurityGroup(ctx context.Context, t *testing.T) (string, func(context.Context, *testing.T)) {
	fwName := "nifwallRuleTest"

	param := &nifcloud.CreateSecurityGroupInput{
		GroupName: fwName,
	}

	Client.C.CreateSecurityGroup(ctx, param)

	done := make(chan struct{}, 0)

	go func() {
		defer close(done)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		for {
			res, err := Client.C.DescribeSecurityGroups(ctx, param)

			if err != nil {
				t.Errorf("Not Created")
			}

			status := res.SecurityGroupInfo[0].GroupStatus

			if status == "applied" {
				break
			}
		}
	}()

	<-done

	return fwName, func(ctx context.Context, t *testing.T) {
		param := &nifcloud.DeleteSecurityGroupInput{
			GroupName: fwName,
		}

		if _, err := Client.C.DeleteSecurityGroup(ctx, param); err != nil {
			t.Fatal(err)
		}
	}
}

func setupTestRegisterInstancesWithSecurityGroup(ctx context.Context, t *testing.T) (string, func(context.Context, *testing.T)) {
	fwName := "nifRegister"

	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        fwName,
		AvailabilityZone: "west-12",
	}

	Client.C.CreateSecurityGroup(ctx, param)

	done := make(chan struct{}, 0)

	go func() {
		defer close(done)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		for {
			res, err := Client.C.DescribeSecurityGroups(ctx, param)

			if err != nil {
				t.Errorf("Not Created")
			}

			status := res.SecurityGroupInfo[0].GroupStatus

			if status == "applied" {
				break
			}
		}
	}()

	<-done

	return fwName, func(ctx context.Context, t *testing.T) {
		param := &nifcloud.DeleteSecurityGroupInput{
			GroupName: fwName,
		}

		if _, err := Client.C.DeleteSecurityGroup(ctx, param); err != nil {
			t.Fatal(err)
		}
	}
}

func TestListInstances(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, err := Client.ListInstances(ctx); err != nil {
		t.Error(err)
	}
}

func TestListInappropriateInstances(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appropriateFWName := "nifwall"

	if _, err := Client.ListInappropriateInstances(ctx, appropriateFWName); err != nil {
		t.Error(err)
	}
}

type mock struct {
	NifCloud
}

func (m *mock) DescribeInstanceAttribute(ctx context.Context, param *nifcloud.DescribeInstanceAttributeInput) (*nifcloud.DescribeInstanceAttributeOutput, error) {
	return m.MockDescribeInstanceAttribute(ctx, param)
}

func (m *mock) MockDescribeInstanceAttribute(ctx context.Context, param *nifcloud.DescribeInstanceAttributeInput) (*nifcloud.DescribeInstanceAttributeOutput, error) {
	return &nifcloud.DescribeInstanceAttributeOutput{
		GroupID: "groupID",
	}, nil
}

func (m *mock) DescribeInstances(ctx context.Context, param *nifcloud.DescribeInstancesInput) (*nifcloud.DescribeInstancesOutput, error) {
	return m.MockDescribeInstances(ctx, param)
}

func (m *mock) MockDescribeInstances(ctx context.Context, param *nifcloud.DescribeInstancesInput) (*nifcloud.DescribeInstancesOutput, error) {
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

func (m *mock) DescribeSecurityGroups(ctx context.Context, param *nifcloud.DescribeSecurityGroupsInput) (*nifcloud.DescribeSecurityGroupsOutput, error) {
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
			},
		},
	}, nil
}

func (m *mock) CreateSecurityGroup(ctx context.Context, param *nifcloud.CreateSecurityGroupInput) (*nifcloud.CreateSecurityGroupOutput, error) {
	return &nifcloud.CreateSecurityGroupOutput{}, nil
}

func (m *mock) DeleteSecurityGroup(ctx context.Context, param *nifcloud.DeleteSecurityGroupInput) (*nifcloud.DeleteSecurityGroupOutput, error) {
	return &nifcloud.DeleteSecurityGroupOutput{}, nil
}

func (m *mock) AuthorizeSecurityGroupIngress(ctx context.Context, param *nifcloud.AuthorizeSecurityGroupIngressInput) (*nifcloud.AuthorizeSecurityGroupIngressOutput, error) {
	return &nifcloud.AuthorizeSecurityGroupIngressOutput{}, nil
}
