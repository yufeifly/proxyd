package container

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yufeifly/migrator/model"
	"github.com/yufeifly/migrator/utils"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

// Start ContainerStart handler
func Start(c *gin.Context) {
	containerID := c.Query("containerId")
	checkpointID := c.Query("checkpointID")
	checkpointDir := c.Query("checkpointDir")

	startOpts := model.StartOpts{
		CStartOpts: types.ContainerStartOptions{
			CheckpointID:  checkpointID,
			CheckpointDir: checkpointDir,
		},
		ContainerID: containerID,
	}

	err := StartContainer(startOpts)
	if err != nil {
		utils.ReportErr(c, err)
		logrus.Panic(err)
	}
	c.JSON(200, gin.H{
		"result": "success",
	})
}

// StartContainer start a container with opts
func StartContainer(startOpts model.StartOpts) error {
	header := "container.StartContainer"
	err := cli.ContainerStart(context.Background(), startOpts.ContainerID, startOpts.CStartOpts)
	if err != nil {
		logrus.Errorf("%s, start container failed, err: %v", header, err)
		return err
	}
	return nil
}
