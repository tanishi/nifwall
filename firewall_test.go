package nifwall

import (
	"context"
	"os"
	"testing"

	nifcloud "github.com/tanishi/go-nifcloud"
)

func TestCreateSecurityGroup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fwName, Teardown := SetupTest(t)
	defer Teardown(ctx, t)

	if err := CreateSecurityGroup(ctx, fwName, fwName); err != nil {
		t.Fatal(err)
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

func SetupTest(t *testing.T) (string, func(context.Context, *testing.T)) {
	endpoint := os.Getenv("NIFCLOUD_ENDPOINT")
	accessKey := os.Getenv("NIFCLOUD_ACCESSKEY")
	secretAccessKey := os.Getenv("NIFCLOUD_SECRET_ACCESSKEY")

	var err error
	client, err = nifcloud.NewClient(endpoint, accessKey, secretAccessKey)

	if err != nil {
		t.Fatal(err)
	}

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
