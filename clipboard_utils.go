package main

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// Get image from clipboard
func getImageFromClipboard() ([]byte, error) {
	data, err := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png", "-out").Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Save link on image to clipboard
func putTextToClipboard(str string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-in")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(str)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}

func getImageFromStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}
