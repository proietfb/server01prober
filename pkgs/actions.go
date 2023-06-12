package actions

import (
	"context"
	"fmt"

	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var dcl *client.Client
var ctx context.Context

type Actions interface {
	GetContainers() string
}

func dClient() (*client.Client, context.Context) {
	var err error
	if dcl == nil {
		ctx = context.Background()
		dcl, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Println("Error on creating client ", err)
		}
	}
	return dcl, ctx
}

func GetContainers() string {
	containers, err := dcl.ContainerList(ctx, types.ContainerListOptions{})
	if err == nil {
		for _, c := range containers {
			fmt.Println(c.ID)
		}
	}
	return ""
}
