package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func NginxReload(nginx string) error {
	cmd := exec.Command(nginx, "-s", "reload")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("NginxReload error: '%s' - '%s'", err, stderr.String())
	}

	return nil
}
