package utils

import "os/exec"

func RunCommandTree(cmd *exec.Cmd) error {
	finish, err := StartCommandTree(cmd)
	if err != nil {
		return err
	}
	defer finish()
	return cmd.Wait()
}
