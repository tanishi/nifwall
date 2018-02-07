package nifwall

import (
	"context"
	"os"
	"reflect"
	"testing"

	nifcloud "github.com/tanishi/go-nifcloud"
)

func TestCreateSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, teardown := setupTestCreateSecurityGroup(t)
	defer teardown(ctx, t)

	if err := CreateSecurityGroup(ctx, fwName, fwName); err != nil {
		t.Error(err)
	}

	done := make(chan *nifcloud.DescribeSecurityGroupsOutput, 0)

	go func() {
		defer close(done)
		for {
			param := &nifcloud.DescribeSecurityGroupsInput{
				GroupNames: []string{fwName},
			}

			res, err := client.DescribeSecurityGroups(ctx, param)

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
		ipPermission{
			Protocol: "HTTP",
			CidrIP:   []string{"0.0.0.0/0"},
		},
	}

	if err := AddRuleToSecurityGroup(ctx, fwName, permissions); err != nil {
		t.Error(err)
	}

	done := make(chan *nifcloud.DescribeSecurityGroupsOutput, 0)

	go func() {
		defer close(done)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		for {
			res, err := client.DescribeSecurityGroups(ctx, param)

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

func TestConvert(t *testing.T) {
	expected := []nifcloud.IPPermission{
		nifcloud.IPPermission{
			IPProtocol:  "HTTP",
			FromPort:    "80",
			ToPort:      "81",
			InOut:       "IN",
			IPRanges:    []string{"0.0.0.0/0"},
			Description: "HOGE",
		},
	}
	actual := convert([]ipPermission{
		ipPermission{
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
			nifcloud.IPPermission{
				IPProtocol: "HTTP",
				IPRanges:   []string{"0.0.0.0/0"},
			},
		},
	}

	actual := generateAuthorizeSecurityGroupIngressInput("HOGE", []ipPermission{
		ipPermission{
			Protocol: "HTTP",
			CidrIP:   []string{"0.0.0.0/0"},
		},
	})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v,\nbut: %v", expected, actual)
	}
}

func commonSetupTest(t *testing.T) {
	endpoint := os.Getenv("NIFCLOUD_ENDPOINT")
	accessKey := os.Getenv("NIFCLOUD_ACCESSKEY")
	secretAccessKey := os.Getenv("NIFCLOUD_SECRET_ACCESSKEY")

	var err error
	client, err = nifcloud.NewClient(endpoint, accessKey, secretAccessKey)

	if err != nil {
		t.Fatal(err)
	}
}

func setupTestCreateSecurityGroup(t *testing.T) (string, func(context.Context, *testing.T)) {
	commonSetupTest(t)

	fwName := "nifwallTest"

	return fwName, func(ctx context.Context, t *testing.T) {
		param := &nifcloud.DeleteSecurityGroupInput{
			GroupName: fwName,
		}

		if _, err := client.DeleteSecurityGroup(ctx, param); err != nil {
			t.Fatal(err)
		}
	}
}

func setupTestAddRuleToSecurityGroup(ctx context.Context, t *testing.T) (string, func(context.Context, *testing.T)) {
	commonSetupTest(t)

	fwName := "nifwallRuleTest"

	param := &nifcloud.CreateSecurityGroupInput{
		GroupName: fwName,
	}

	client.CreateSecurityGroup(ctx, param)

	done := make(chan struct{}, 0)

	go func() {
		defer close(done)

		param := &nifcloud.DescribeSecurityGroupsInput{
			GroupNames: []string{fwName},
		}

		for {
			res, err := client.DescribeSecurityGroups(ctx, param)

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

		if _, err := client.DeleteSecurityGroup(ctx, param); err != nil {
			t.Fatal(err)
		}
	}
}
