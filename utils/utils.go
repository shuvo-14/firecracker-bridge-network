package utils

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	firecrackerBin = "/usr/local/bin/firecracker"
	bridgeName     = "br0"
	tap0Name       = "tap0"
	tap1Name       = "tap1"
)

func DownloadFile(url, dest string) error {
	cmd := exec.Command("curl", "-o", dest, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func MakeExecutable(script string) error {
	cmd := exec.Command("chmod", "+x", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunShellScript(script string, args ...string) error {
	if err := MakeExecutable(script); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	cmd := exec.Command("/bin/bash", append([]string{script}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func LoadAndSendConfig(socketPath, urlPath, configFile string) error {
	cmd := exec.Command("curl",
		"--unix-socket", socketPath,
		"-i", "-X", "PUT",
		fmt.Sprintf("http://localhost/%s", urlPath),
		"-H", "Content-Type: application/json",
		"-d", fmt.Sprintf("@%s", configFile),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StartFirecrackerVM(socketPath, logFilePath string) error {
	cmd := exec.Command(firecrackerBin, "--api-sock", socketPath)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	return cmd.Start()
}
