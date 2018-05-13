package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"docker.io/go-docker"
	types "docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"golang.org/x/net/context"

	noobTypes "github.com/peterdeme/noobernetes/types"
)

// StartContainers starts a container by task def
func StartContainers(w http.ResponseWriter, req *http.Request) {
	taskDef := new(noobTypes.TaskDefinition)
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, taskDef)

	if err != nil {
		http.Error(w, "damn it", http.StatusBadRequest)
		return
	}

	err = pullContainerAndStart(req.Context(), taskDef.ImageURL)

	if err != nil {
		res, _ := json.Marshal(err)
		w.Write(res)
	}
}

func pullContainerAndStart(ctx context.Context, imageName string) error {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	return runContainer(ctx, cli, imageName, 5)
}

func runContainer(ctx context.Context, cli *docker.Client, imageName string, count int) error {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")

	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

		if err != nil {
			return err
		}
	}

	return nil
}
