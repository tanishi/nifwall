package nifwall

import (
	"context"

	nifcloud "github.com/tanishi/go-nifcloud"
)

var client *nifcloud.Client

func CreateSecurityGroup(ctx context.Context, name, description string) error {
	param := &nifcloud.CreateSecurityGroupInput{
		GroupName:        name,
		GroupDescription: description,
	}

	_, err := client.CreateSecurityGroup(ctx, param)

	return err
}
