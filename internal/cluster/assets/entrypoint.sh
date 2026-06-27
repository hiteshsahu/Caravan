#!/usr/bin/env bash
set -euo pipefail

# munge (cluster auth)
install -d -o munge -g munge -m 0755 /run/munge
install -d -o munge -g munge -m 0700 /var/lib/munge /var/log/munge
chown munge:munge /etc/munge/munge.key
chmod 400 /etc/munge/munge.key
runuser -u munge -- /usr/sbin/munged
sleep 1

install -d -m 0755 /var/spool/slurmctld /var/spool/slurmd /var/log/slurm

case "${1:-}" in
  slurmctld) echo "[entrypoint] starting slurmctld"; exec /usr/sbin/slurmctld -D ;;
  slurmd)    echo "[entrypoint] starting slurmd as $(hostname)"; exec /usr/sbin/slurmd -D ;;
  *)         exec "$@" ;;
esac
