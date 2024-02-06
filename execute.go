package discordyetanotherremoterunner

import (
	"fmt"
	"os/exec"
)

func execute(cmd ConfigCommand) (string, error) {
	proc := exec.Command(cmd.Cmd, cmd.Args...)
	out, err := proc.CombinedOutput()
	if err != nil {
		return "", err
	}
	code := proc.ProcessState.ExitCode()
	outStr := string(out)
	if code != 0 {
		return "", fmt.Errorf("failed to execute %s: %s", cmd.Cmd, outStr)
	}
	return outStr, nil
}
