#!/bin/bash
#
# Start CRC/OpenShift Local with optimal settings for observability stack
# (otel-operator, tempo-operator, grafana, loki, etc.)
#

# Resource allocation
CRC_CPUS=8
CRC_MEMORY=16384  # 16GB in MB
CRC_DISK_SIZE=60  # GB

crc config set enable-cluster-monitoring true

echo "Starting CRC with:"
echo "  CPUs: ${CRC_CPUS} cores"
echo "  Memory: ${CRC_MEMORY}MB"
echo "  Disk: ${CRC_DISK_SIZE}GB"
echo ""

crc start \
  --cpus "${CRC_CPUS}" \
  --memory "${CRC_MEMORY}" \
  --disk-size "${CRC_DISK_SIZE}"
