package actions

import (
	"context"

	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ActionsData struct {
	dcl *client.Client
	ctx context.Context
}

type Actions interface {
	Init() *Actions
	GetContainers() string
}

func Init() *ActionsData {
	var err error
	ac := ActionsData{}
	if ac.dcl == nil {
		ac.ctx = context.Background()
		log.Println(client.WithAPIVersionNegotiation())
		ac.dcl, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Println("Error on creating client ", err)
		}
	}
	return &ac
}

func (ad *ActionsData) GetContainers() string {
	containers, err := ad.dcl.ContainerList(ad.ctx, types.ContainerListOptions{})
	var str string
	if err == nil {
		for _, c := range containers {
			var names string
			for _, n := range c.Names {
				names += n + "\n"
			}
			str += names + ": " + c.Image + ",status: " + c.Status + "\n"
		}
	}
	return str
}
