package main

import (
	"fmt"
	"io"
	"os/exec"
)

type choices struct {
	c []string
}

func isNewLine(candidate byte) bool {
	return candidate == 10
}

func choicesReader(r io.Reader) (choices, error) {
	var err error
	ch := make([]string, 0, 32)
	buf := make([]byte, maxFileNameLen)
	for {
		n, err2 := r.Read(buf)
		if err2 == io.EOF {
			break
		}
		if n > 0 && err2 != nil {
			err = err2
			break
		}
		lastLineFeed := 0
		for i := 0; i < n; i++ {
			if isNewLine(buf[i]) {
				ch = append(ch, string(buf[lastLineFeed:i]))
				lastLineFeed = i + 1
			}
		}
		buf = buf[:n]
	}
	return choices{ch}, err
}

func readFileNamesViaLs() (choices, error) {
	cmd := exec.Command("ls", "--file-type", "--group-directories-first")
	pipe, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return choices{}, fmt.Errorf("init command start: %v", err)
	}
	defer pipe.Close()
	if err != nil {
		return choices{}, fmt.Errorf("model i/o init: %v", err)
	}
	ch, err := choicesReader(pipe)
	if err != nil {
		return choices{}, fmt.Errorf("read command output: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return choices{}, fmt.Errorf("execution init command: %v", err)
	}
	return ch, nil
}
