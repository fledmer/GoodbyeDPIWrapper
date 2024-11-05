package proc

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Async struct {
	Proc *exec.Cmd
}

func (s *Async) Start(path string, args ...string) (*exec.Cmd, error) {
	if err := s.singleProcCheck(); err != nil {
		return nil, fmt.Errorf("failed to start new proc single check failed, %w", err)
	}

	proc := exec.Command(path, args...)
	if proc.Err != nil {
		return nil, fmt.Errorf("failed to start new proc %w", proc.Err)
	}

	outReader, err := proc.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to read stdout pipe: %w", err)
	}
	if err = proc.Start(); err != nil {
		log.Println("failed to start process: ", err)
		return nil, fmt.Errorf("failed to start process: %w", err)
	} else {
		go s.ReadOutput(outReader)
		s.Proc = proc
	}

	return proc, nil
}

func (s *Async) singleProcCheck() (err error) {
	if s.Proc != nil {
		return fmt.Errorf("proc with pid: %v already run", s.Proc.Process.Pid)
	}
	return
}

func (s *Async) Stop() (err error) {
	if s.Proc == nil {
		return errors.New("no running process")
	}

	s.Proc.Process.Kill()
	state, err := s.Proc.Process.Wait()
	if err != nil {
		return fmt.Errorf("failed to stop proc: %w", err)
	}
	if !state.Exited() {
		return fmt.Errorf("process is not exited: %w", err)
	}
	s.Proc = nil
	return nil
}

func (s *Async) ReadOutput(reader io.ReadCloser) {
	bytes := make([]byte, 4096)
	var err error
	var n int
	for err == nil {
		n, err = reader.Read(bytes)
		os.Stdout.Write(bytes[:n])
	}
}

type Sync struct {
}

func StartSync(path string, args ...string) error {
	proc := exec.Command(path, args...)
	if proc.Err != nil {
		return fmt.Errorf("failed to start new proc %w", proc.Err)
	}
	outReader, err := proc.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to read stdout pipe: %w", err)
	}
	go readOutput(outReader)
	outReader, err = proc.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to read stderr pipe: %w", err)
	}
	go readOutput(outReader)
	err = proc.Run()
	return errors.Join(err, proc.Err)
}

func readOutput(reader io.ReadCloser) {
	bytes := make([]byte, 4096)
	var err error
	var n int
	for err == nil {
		n, err = reader.Read(bytes)
		os.Stdout.Write(bytes[:n])
	}
}
