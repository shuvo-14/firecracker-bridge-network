package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateBridgeNetwork() error {
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
