package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	imageURLBase   = "https://s3.amazonaws.com/spec.ccfc.min/img/hello/"
	kernelImageURL = imageURLBase + "kernel/hello-vmlinux.bin"
	rootfsImageURL = imageURLBase + "fsfiles/hello-rootfs.ext4"
)

var (
	projectRoot     = "./"
	firecrackerBin  = "/usr/local/bin/firecracker"
	scriptsDir      = filepath.Join(projectRoot, "scripts")
	configDir       = filepath.Join(projectRoot, "config")
	imagesDir       = filepath.Join(projectRoot, "images")
	rootfsImagePath = filepath.Join(imagesDir, "hello-rootfs.ext4")
	bridgeName      = "br0"
	tap0Name        = "tap0"
	tap1Name        = "tap1"
)

func createBridgeNetwork() error {
	cmd := exec.Command("sudo", "ip", "link", "add", bridgeName, "type", "bridge")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create bridge using ip link: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "addr", "add", "192.168.0.1/24", "dev", bridgeName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to assign IP address to bridge: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", bridgeName, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring up bridge: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "tuntap", "add", tap0Name, "mode", "tap")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tap0: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "addr", "add", "192.168.0.2/24", "dev", tap0Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to assign IP address to tap0: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", tap0Name, "master", bridgeName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add tap0 to bridge: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", tap0Name, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring up tap0: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "tuntap", "add", tap1Name, "mode", "tap")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tap1: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "addr", "add", "192.168.0.3/24", "dev", tap1Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to assign IP address to tap1: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", tap1Name, "master", bridgeName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add tap1 to bridge: %v", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", tap1Name, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring up tap1: %v", err)
	}

	return nil
}

func downloadFile(url, dest string) error {
	cmd := exec.Command("curl", "-o", dest, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func makeExecutable(script string) error {
	cmd := exec.Command("chmod", "+x", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runShellScript(script string, args ...string) error {
	if err := makeExecutable(script); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	cmd := exec.Command("/bin/bash", append([]string{script}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func loadAndSendConfig(socketPath, urlPath, configFile string) error {
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

func startFirecrackerVM(socketPath, logFilePath string) error {
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

func main() {
	if err := createBridgeNetwork(); err != nil {
		log.Fatalf("Error creating bridge network: %v", err)
	}

	if err := downloadFile(kernelImageURL, filepath.Join(imagesDir, "hello-vmlinux.bin")); err != nil {
		log.Fatalf("Error downloading kernel image: %v", err)
	}
	if err := downloadFile(rootfsImageURL, rootfsImagePath); err != nil {
		log.Fatalf("Error downloading root filesystem image: %v", err)
	}

	if err := runShellScript(filepath.Join(scriptsDir, "prepare-rootfs.sh"), rootfsImagePath); err != nil {
		log.Fatalf("Error preparing root filesystem: %v", err)
	}

	socketPathVM1 := "/tmp/firecracker-vm1.sock"
	socketPathVM2 := "/tmp/firecracker-vm2.sock"
	logFileVM1 := "/tmp/firecracker-vm1.log"
	logFileVM2 := "/tmp/firecracker-vm2.log"

	if err := startFirecrackerVM(socketPathVM1, logFileVM1); err != nil {
		log.Fatalf("Error starting Firecracker VM1: %v", err)
	}
	if err := startFirecrackerVM(socketPathVM2, logFileVM2); err != nil {
		log.Fatalf("Error starting Firecracker VM2: %v", err)
	}

	if err := loadAndSendConfig(socketPathVM1, "boot-source", filepath.Join(configDir, "vm1", "kernel.json")); err != nil {
		log.Fatalf("Error configuring VM1 boot source: %v", err)
	}
	if err := loadAndSendConfig(socketPathVM1, "drives", filepath.Join(configDir, "vm1", "rootfs.json")); err != nil {
		log.Fatalf("Error configuring VM1 root filesystem: %v", err)
	}
	if err := loadAndSendConfig(socketPathVM1, "network-interfaces", filepath.Join(configDir, "vm1", "network.json")); err != nil {
		log.Fatalf("Error configuring VM1 network: %v", err)
	}

	if err := loadAndSendConfig(socketPathVM2, "boot-source", filepath.Join(configDir, "vm2", "kernel.json")); err != nil {
		log.Fatalf("Error configuring VM2 boot source: %v", err)
	}
	if err := loadAndSendConfig(socketPathVM2, "drives", filepath.Join(configDir, "vm2", "rootfs.json")); err != nil {
		log.Fatalf("Error configuring VM2 root filesystem: %v", err)
	}
	if err := loadAndSendConfig(socketPathVM2, "network-interfaces", filepath.Join(configDir, "vm2", "network.json")); err != nil {
		log.Fatalf("Error configuring VM2 network: %v", err)
	}

	fmt.Println("Firecracker VMs started successfully. You can now login and ping between the VMs.")
}
