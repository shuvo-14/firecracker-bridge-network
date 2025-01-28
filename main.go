package main

import (
	"fmt"
	"github.com/shuvo-14/firecracker-bridge-network/utils"
	"log"
	"path/filepath"
)

const (
	imageURLBase   = "https://s3.amazonaws.com/spec.ccfc.min/img/hello/"
	kernelImageURL = imageURLBase + "kernel/hello-vmlinux.bin"
	rootfsImageURL = imageURLBase + "fsfiles/hello-rootfs.ext4"
	socketPathVM1  = "/tmp/firecracker-vm1.sock"
	socketPathVM2  = "/tmp/firecracker-vm2.sock"
	logFileVM1     = "/tmp/firecracker-vm1.log"
	logFileVM2     = "/tmp/firecracker-vm2.log"
)

var (
	projectRoot     = "./"
	scriptsDir      = filepath.Join(projectRoot, "scripts")
	configDir       = filepath.Join(projectRoot, "config")
	imagesDir       = filepath.Join(projectRoot, "images")
	rootfsImagePath = filepath.Join(imagesDir, "hello-rootfs.ext4")
)

func main() {
	if err := utils.CreateBridgeNetwork(); err != nil {
		log.Fatalf("Error creating bridge network: %v", err)
	}

	if err := utils.DownloadFile(kernelImageURL, filepath.Join(imagesDir, "hello-vmlinux.bin")); err != nil {
		log.Fatalf("Error downloading kernel image: %v", err)
	}
	if err := utils.DownloadFile(rootfsImageURL, rootfsImagePath); err != nil {
		log.Fatalf("Error downloading root filesystem image: %v", err)
	}

	if err := utils.RunShellScript(filepath.Join(scriptsDir, "prepare-rootfs.sh"), rootfsImagePath); err != nil {
		log.Fatalf("Error preparing root filesystem: %v", err)
	}

	if err := utils.StartFirecrackerVM(socketPathVM1, logFileVM1); err != nil {
		log.Fatalf("Error starting Firecracker VM1: %v", err)
	}
	if err := utils.StartFirecrackerVM(socketPathVM2, logFileVM2); err != nil {
		log.Fatalf("Error starting Firecracker VM2: %v", err)
	}

	if err := utils.LoadAndSendConfig(socketPathVM1, "boot-source", filepath.Join(configDir, "vm1", "kernel.json")); err != nil {
		log.Fatalf("Error configuring VM1 boot source: %v", err)
	}
	if err := utils.LoadAndSendConfig(socketPathVM1, "drives", filepath.Join(configDir, "vm1", "rootfs.json")); err != nil {
		log.Fatalf("Error configuring VM1 root filesystem: %v", err)
	}
	if err := utils.LoadAndSendConfig(socketPathVM1, "network-interfaces", filepath.Join(configDir, "vm1", "network.json")); err != nil {
		log.Fatalf("Error configuring VM1 network: %v", err)
	}

	if err := utils.LoadAndSendConfig(socketPathVM2, "boot-source", filepath.Join(configDir, "vm2", "kernel.json")); err != nil {
		log.Fatalf("Error configuring VM2 boot source: %v", err)
	}
	if err := utils.LoadAndSendConfig(socketPathVM2, "drives", filepath.Join(configDir, "vm2", "rootfs.json")); err != nil {
		log.Fatalf("Error configuring VM2 root filesystem: %v", err)
	}
	if err := utils.LoadAndSendConfig(socketPathVM2, "network-interfaces", filepath.Join(configDir, "vm2", "network.json")); err != nil {
		log.Fatalf("Error configuring VM2 network: %v", err)
	}

	fmt.Println("Firecracker VMs started successfully. You can now login and ping between the VMs.")
}
