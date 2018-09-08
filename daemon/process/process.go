package process

import (
	"github.com/huaweicloud/telescope/agent/core/manager"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/huaweicloud/telescope/agent/core/logs"
)

var (
	sendSignalRetryCount = 3
)

//get current agent version, command: agent -version
func GetAgentVersion(binPath string) (string, error) {
	cmd := exec.Command(binPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

//signal agent, command: agent stop/upgrade
func SendAgentSignal(binPath string, signal os.Signal) error {
	var cmd *exec.Cmd
	switch signal {
	case manager.SIG_STOP:
		cmd = exec.Command(binPath, "stop")
	default:
		cmd = exec.Command(binPath, "-version")
	}

	_, err := cmd.Output()
	return err
}

// start new child process
func StartProcess(binPath string) (*os.Process, error) {
	cmd := exec.Command(binPath)
	e := os.Environ()
	cmd.Env = e
	cmd.Args = os.Args
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logs.GetLogger().Errorf("Start process failed, err:%s", err.Error())
		return nil, err
	}

	proc := cmd.Process
	//wait process, otherwise the process becomes zombie after kill
	cmdwait := make(chan error)
	go func() {
		cmdwait <- cmd.Wait()
	}()
	return proc, nil
}

// kill child process
func KillProcess(proc *os.Process) error {
	if proc == nil {
		return nil
	}

	osName := runtime.GOOS
	if osName == "windows" {
		return proc.Kill()
	} else {
		err := proc.Signal(syscall.SIGKILL)
		if err != nil {
			logs.GetLogger().Errorf("Stop(SIG_UPGRADE) linux agent process failed, err:%s", err.Error())
			return err
		}
		// wait old process finished
		proc.Wait()
		return nil
	}
}

// send and kill child process
func SigAndKillProcess(binPath string, signal os.Signal, proc *os.Process) error {
	if proc == nil {
		return nil
	}
	err := proc.Kill()

	if err != nil {
		return err
	}
	tryCount := 0
	for tryCount < sendSignalRetryCount {
		err = SendAgentSignal(binPath, signal)
		if err == nil {
			break
		}
		tryCount = tryCount + 1
	}

	return nil
}

func StopProcess(binPath string, proc *os.Process) error {
	return SigAndKillProcess(binPath, manager.SIG_STOP, proc)
}
