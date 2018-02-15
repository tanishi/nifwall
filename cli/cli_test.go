package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	nifcloud "github.com/tanishi/go-nifcloud"
	nifwall "github.com/tanishi/nifwall/lib"
)

func TestRun(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		outStream := new(bytes.Buffer)

		cli := &CLI{
			OutStream: outStream,
			ErrStream: outStream,
			Client:    nifwall.NewClient(&Mock{}),
		}

		status := cli.Run(strings.Split("nifwall list", " "))

		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		actual := outStream.String()
		expected := []string{"[", "test1", "test2", "test3", "test4", "test5", "]"}
		for _, e := range expected {
			if !strings.Contains(actual, e) {
				t.Errorf("expected: %v, but: %v", expected, actual)
			}
		}
	})
	t.Run("update", func(t *testing.T) {
		outStream := new(bytes.Buffer)

		cli := &CLI{
			OutStream: outStream,
			ErrStream: outStream,
			Client:    nifwall.NewClient(&Mock{}),
		}

		status := cli.Run(strings.Split("nifwall update -f ../examples/example.yml", " "))

		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		actual := outStream.String()
		expected := ""

		if expected != actual {
			t.Errorf("expected: %v, but: %v", expected, actual)
		}
	})
	t.Run("apply", func(t *testing.T) {
		outStream := new(bytes.Buffer)

		cli := &CLI{
			OutStream: outStream,
			ErrStream: outStream,
			Client:    nifwall.NewClient(&Mock{}),
		}

		status := cli.Run(strings.Split("nifwall apply test", " "))

		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		actual := outStream.String()
		expected := ""

		if expected != actual {
			t.Errorf("expected: %v, but: %v", expected, actual)
		}
	})

	cases := []struct {
		args     string
		expected string
	}{
		{"nifwall -version", fmt.Sprintf("nifwall version %s\n", version)},
		{"nifwall", fmt.Sprintf("list or update or apply")},
	}

	for _, c := range cases {
		outStream := new(bytes.Buffer)

		cli := &CLI{
			OutStream: outStream,
			ErrStream: outStream,
			Client:    nifwall.NewClient(&Mock{}),
		}

		status := cli.Run(strings.Split(c.args, " "))

		if status != exitCodeOK {
			t.Errorf("ExitStatus=%d, want %d", status, exitCodeOK)
		}

		actual := outStream.String()

		if c.expected != actual {
			t.Errorf("expected: %v, but: %v", c.expected, actual)
		}
	}
}

type Mock struct {
	nifwall.NifCloud
}

func (m *Mock) DescribeInstanceAttribute(ctx context.Context, param *nifcloud.DescribeInstanceAttributeInput) (*nifcloud.DescribeInstanceAttributeOutput, error) {
	return m.MockDescribeInstanceAttribute(ctx, param)
}

func (m *Mock) MockDescribeInstanceAttribute(ctx context.Context, param *nifcloud.DescribeInstanceAttributeInput) (*nifcloud.DescribeInstanceAttributeOutput, error) {
	return &nifcloud.DescribeInstanceAttributeOutput{
		GroupID: "groupID",
	}, nil
}

func (m *Mock) DescribeInstances(ctx context.Context, param *nifcloud.DescribeInstancesInput) (*nifcloud.DescribeInstancesOutput, error) {
	return m.MockDescribeInstances(ctx, param)
}

func (m *Mock) MockDescribeInstances(ctx context.Context, param *nifcloud.DescribeInstancesInput) (*nifcloud.DescribeInstancesOutput, error) {
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
