

## Running the Project

1. Build the Go project:

    ```sh
    go build -o firecracker-bridge-network main.go
    ```

2. Run the executable:

    ```sh
    sudo ./firecracker-bridge-network
    ```

    This will:
    - Create a bridge network.
    - Download the kernel and root filesystem images.
    - Prepare the root filesystem.
    - Start the Firecracker VMs.
    - Configure the VMs.

3. Once the VMs are started, you can log in and ping between them.

    ```sh
    ssh -i path/to/ssh/key root@vm1_ip
    ping vm2_ip
    ```

## Cleaning Up

To clean up the created resources, you can manually remove the socket and log files:

```sh
sudo rm /tmp/firecracker-vm1.sock /tmp/firecracker-vm2.sock /tmp/firecracker-vm1.log /tmp/firecracker-vm2.log