package redis

import (
	"fmt"
	"github.com/yufeifly/migrator/model"
	"testing"
)

func TestTryMigrate(t *testing.T) {
	migrateOpts := model.MigrateOpts{
		Container:     "85ea0420bb58", // to identify the container in source node
		CheckpointID:  "cp-redis",
		CheckpointDir: "/tmp",
		DestIP:        "127.0.0.1",
		DestPort:      "6789",
	}
	err := TryMigrate(migrateOpts)
	if err != nil {
		fmt.Printf("TestTryMigrate err: %v\n", err)
	} else {
		fmt.Printf("TestTryMigrate pass\n")
	}
}
