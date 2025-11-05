#!/bin/sh


sudo su
mkdir -p /mnt/data/compose

cat > /etc/default/minio << EOF
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_VOLUMES=/mnt/data
EOF