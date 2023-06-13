package actions

import (
	"context"

	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ActionsData struct {
	dcl *client.Client
	ctx context.Context
}

type Actions interface {
	Init() *Actions
	GetStatus() string
	GetContainersName() []string
	GetContainerID(name string) string
	GetContainerLog(name string) string
	RestartContainer(name string) bool
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

func (ad *ActionsData) GetContainerID(name string) string {
	containers, err := ad.dcl.ContainerList(ad.ctx, types.ContainerListOptions{})
	var id string
	if err == nil {
		for _, c := range containers {
			if c.Names[0] == name {
				id = name
			}
		}
	}
	return id
}

func (ad *ActionsData) GetContainersName() []string {
	containers, err := ad.dcl.ContainerList(ad.ctx, types.ContainerListOptions{})
	var names []string
	if err == nil {
		for _, c := range containers {
			names = append(names, c.Names[0])
		}
	}
	return names

}

func (ad *ActionsData) GetContainerlog(name string) string {
	id := ad.GetContainerID(name)
	if io, err := ad.dcl.ContainerLogs(ad.ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: true,
	}); err == nil {
		defer io.Close()
		buff := make([]byte, 1024*1024*512)
		n, e := io.Read(buff)
		log.Println("Err: ", e, "\n", string(buff[:n]))
		return string(buff[:n])

	}
	return ""
}

func (ad *ActionsData) RestartContainer(name string) bool {
	ret := false
	id := ad.GetContainerID(name)
	if err := ad.dcl.ContainerRestart(ad.ctx, id, container.StopOptions{}); err == nil {
		ret = true
	}

	return ret
}

func (ad *ActionsData) GetStatus() string {
	containers, err := ad.dcl.ContainerList(ad.ctx, types.ContainerListOptions{})
	var status string
	if err == nil {
		for _, c := range containers {
			var names string
			for _, n := range c.Names {
				names += "- " + n[1:] + ":\n"
			}
			status += names + "\t\t\t*Status*: " + c.Status + "\n\t\t\t*State*: " + c.State + "\n\n"
		}
	}
	return status
}
