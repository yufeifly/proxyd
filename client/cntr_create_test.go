package client

import (
	"encoding/json"
	"fmt"
	"github.com/yufeifly/migrator/api/types"
	"testing"
)

func TestCli_SendContainerCreate(t *testing.T) {
	cli := NewClient(types.Address{
		IP:   "127.0.0.1",
		Port: "6789",
	})
	cmdSlice := []string{"/bin/sh", "-c", "i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done"}
	cmd, err := json.Marshal(&cmdSlice)
	fmt.Printf("cmd: %v\n", string(cmd))

	opts := types.CreateReqOpts{
		CreateOpts: types.CreateOpts{
			ContainerName: "bb22",
			ImageName:     "busybox",
			HostPort:      "",
			ContainerPort: "",
			Cmd:           string(cmd),
		},
		Address: types.Address{
			IP:   "127.0.0.1",
			Port: "6789",
		},
	}
	got, err := cli.SendContainerCreate(opts)
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		var ans map[string]interface{}
		json.Unmarshal(got, &ans)
		fmt.Printf("create result: %v\n", ans["containerId"])
	}
}
