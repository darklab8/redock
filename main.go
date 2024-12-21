package main

import (
	"flag"
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/darklab8/go-typelog/examples/logus"
	"github.com/darklab8/go-typelog/typelog"
)

func main() {
	logus.Log.Info("entrypoint of cli interface")

	container_id := flag.String("ctr", "", "Container ID or Name you wish to restart")

	image_name := flag.String("image_name", "", "Image name to repull and use for recreation. Default to use old one")

	strict := flag.Bool("strict_pull", false, "Strict mod. fails if failing to pull immage.")
	flag.Parse()

	if *container_id == "" {
		logus.Log.Fatal("ctr is required param, like `redock --ctr=darkbot`")
	}

	redock(*container_id, *strict, *image_name)
}

func redock(container_id string, strict bool, image_name string) {
	logus.Log.Info("deploy is called")

	client, _ := docker.NewClientFromEnv()

	// get container info
	old_container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: container_id,
	})
	logus.Log.CheckPanic(err, "get container info")

	//TODO delete _new if an error occores?

	target_image_name := ""
	if image_name != "" {
		target_image_name = image_name
	} else {
		image_name = old_container.Config.Image
	}

	// pull new image
	fmt.Printf("Image: %s\n", target_image_name)
	if err := client.PullImage(docker.PullImageOptions{
		Repository: target_image_name,
	}, docker.AuthConfiguration{}); err != nil {
		logus.Log.CheckError(err, "failed repulling the img")
		if strict {
			logus.Log.Fatal("strict mod. failing repulling image is not alllowed")
		}
	} else {
		logus.Log.Info("succesfully repulled image", typelog.Any("image", target_image_name))
	}

	//TODO handle image tags/labels?

	// naming
	name := old_container.Name
	tmp_name := name + "_new"

	// copy container
	var options docker.CreateContainerOptions
	options.Name = tmp_name
	options.Config = old_container.Config
	options.Config.Image = target_image_name
	options.HostConfig = old_container.HostConfig
	// get all vomumes
	options.HostConfig.VolumesFrom = []string{old_container.ID}

	// get all links
	// this is a hack to fix the way links are returned and sent
	links := old_container.HostConfig.Links
	for i, link := range links {
		parts := strings.SplitN(link, ":", 2)
		if len(parts) != 2 {
			logus.Log.Error("Unable to parse link ", typelog.Any("link", link))
			return
		}
		container_name := strings.TrimPrefix(parts[0], "/")
		alias_parts := strings.Split(parts[1], "/")
		alias := alias_parts[len(alias_parts)-1]
		links[i] = fmt.Sprintf("%s:%s", container_name, alias)
	}
	options.HostConfig.Links = links

	fmt.Println("Creating...")
	new_container, err := client.CreateContainer(options)
	logus.Log.CheckPanic(err, "creating container")

	// rename
	err = client.RenameContainer(docker.RenameContainerOptions{ID: old_container.ID, Name: name + "_old"})
	logus.Log.CheckPanic(err, "renaming old ctr to old")
	err = client.RenameContainer(docker.RenameContainerOptions{ID: new_container.ID, Name: name})
	logus.Log.CheckPanic(err, "renaming new to default name")

	if old_container.State.Running {
		fmt.Printf("Stopping old container\n")
		err = client.StopContainer(old_container.ID, 10)
		logus.Log.CheckPanic(err, "stopping old container")
		fmt.Printf("Starting new container\n")
		err = client.StartContainer(new_container.ID, new_container.HostConfig)
		logus.Log.CheckPanic(err, "starting new container")
	}

	// add optionn to rm old container on sucsess
	fmt.Printf("Migrated from %s to %s\n", old_container.ID, new_container.ID)
	client.RemoveContainer(docker.RemoveContainerOptions{ID: old_container.ID})

	fmt.Println("Done")
}
