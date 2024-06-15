package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"sherry/shr/config"
	"sherry/shr/helpers"
	"strings"
	"time"
)

func getServicePath() string {
	return helpers.PreparePath(path.Join(config.GetConfigPath(), "bin", helpers.If(runtime.GOOS == "windows", func() string {
		return "sherry-demon.exe"
	}, func() string {
		return "sherry-demon"
	})))
}

func getPidPath() string {
	return helpers.PreparePath(path.Join(config.GetConfigPath(), "pid"))
}

func getWindowsStartServiceCommand() []string {
	ps, _ := exec.LookPath("powershell.exe")
	servicePath := getServicePath()
	configPath := config.GetConfigPath()
	return []string{
		ps,
		"-NoProfile", "-NonInteractive", "-Command",
		fmt.Sprintf(`
$JOB_NAME = "SherryManualStart"
Start-Job -Name $JOB_NAME -ScriptBlock {
	(Start-Process -FilePath "%s" -NoNewWindow -PassThru -WorkingDirectory "%s" -ArgumentList @("-c", "%s", "-s")).Id
} | out-null
echo (Receive-Job -AutoRemoveJob -Wait $JOB_NAME)
exit
`,
			servicePath, configPath, configPath,
		),
	}
}

func StartService(yes bool) bool {
	pid, e := os.ReadFile(getPidPath())
	if e == nil && string(pid) != "" {
		if yes {
			StopService()
		} else {
			helpers.PrintErr(fmt.Sprintf("Service is already started, PID: %s", pid))
			return false
		}
	}

	servicePath := getServicePath()

	helpers.PrintMessage(fmt.Sprintf("Starting service at %s", servicePath))

	var args []string
	switch runtime.GOOS {
	case "windows":
		args = getWindowsStartServiceCommand()
	default:
		helpers.PrintErr("not support")
		return false
	}

	var cmdOut, cmdErr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	err := cmd.Start()

	helpers.PrintMessage("Starting...")
	time.Sleep(time.Second * 1)
	if cmdErr.Len() > 0 {
		helpers.PrintErr(cmdErr.String())
		return false
	}
	if err != nil {
		helpers.PrintErr(err.Error())
		return false
	}

	pid = regexp.MustCompile("[0-9]+").Find(cmdOut.Bytes())

	helpers.PrintMessage(fmt.Sprintf("The pid is %s", pid))
	_ = os.WriteFile(getPidPath(), pid, 0644)

	return false
}

func StopService() bool {
	pid, err := os.ReadFile(getPidPath())
	if err != nil {
		helpers.PrintErr("Service is not started")
		return false
	}

	_ = os.Remove(getPidPath())

	var out []byte
	switch runtime.GOOS {
	case "windows":
		ps, _ := exec.LookPath("powershell.exe")
		out, err = exec.Command(
			ps, "-NoProfile", "-NonInteractive",
			fmt.Sprintf(`kill %s`, pid),
		).Output()
	default:
		helpers.PrintErr("not support")
		return false
	}
	if err != nil {
		helpers.PrintErr(err.Error())
		return false
	}

	if strings.TrimSpace(string(out)) != "" {
		helpers.PrintMessage(string(out))
	}
	helpers.PrintMessage("Service stopped")

	return false
}
