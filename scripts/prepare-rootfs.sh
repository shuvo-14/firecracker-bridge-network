#!/bin/bash
set -e

ROOTFS_IMAGE=$1
MOUNT_DIR=$(mktemp -d)

if [ -z "$ROOTFS_IMAGE" ]; then
  echo "Usage: $0 <path_to_rootfs_image>"
  exit 1
fi

echo "Mounting root filesystem..."
sudo mount -o loop "$ROOTFS_IMAGE" "$MOUNT_DIR"

echo "Adding user 'admin' with password 'password'..."
sudo chroot "$MOUNT_DIR" /bin/sh -c "
  echo 'root:password' | chpasswd
  adduser -D admin
  echo 'admin:password' | chpasswd
"

echo "Unmounting root filesystem..."
sudo umount "$MOUNT_DIR"
rmdir "$MOUNT_DIR"
echo "Root filesystem prepared successfully."
