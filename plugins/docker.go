package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Docker collects information about docker
type Docker struct{}

// ShouldRun ensures we're on a platform we support
func (p Docker) ShouldRun(ctx PluginContext) {
	if _, err := os.Stat("/var/run/docker.sock"); os.IsNotExist(err) {
		ctx.Finish <- false
		return
	}

	ctx.Finish <- true
}

// Run executes this plugin
func (p Docker) Run(ctx PluginContext) {
	cli, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		panic(err)
	}

	version, err := cli.ServerVersion(context.Background())

	if err != nil {
		panic(err)
	}

	ctx.Logger.Debug("Containers Running: ", len(containers))
	ctx.Logger.Debug("Images found: ", len(images))
	ctx.Logger.Debug("Networks found: ", len(networks))
	ctx.Logger.Debug("Server version:", version)

	containersJSON, err := json.Marshal(containers)

	if err != nil {
		panic(err)
	}

	imagesJSON, err := json.Marshal(images)

	if err != nil {
		panic(err)
	}

	networksJSON, err := json.Marshal(networks)

	if err != nil {
		panic(err)
	}

	versionJSON, err := json.Marshal(version)

	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/containers.json", ctx.OutputPath), containersJSON, 0600); err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/images.json", ctx.OutputPath), imagesJSON, 0600); err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/networks.json", ctx.OutputPath), networksJSON, 0600); err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/version.json", ctx.OutputPath), versionJSON, 0600); err != nil {
		panic(err)
	}

	ctx.Finish <- true
}
