package clipboard

import "os/exec"

const (
	copyCmd = "pbcopy"
)

func Write(text string) error {
	cmd := exec.Command(copyCmd)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}
