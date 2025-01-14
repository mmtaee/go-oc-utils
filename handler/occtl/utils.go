package occtl

import (
	"context"
	"fmt"
	"os/exec"
)

// Exec occtl execute command
func Exec(c context.Context, command string) ([]byte, error) {
	cmd := exec.CommandContext(c, "sh", "-c", fmt.Sprintf("/usr/bin/occtl %s", command))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
