#!/usr/bin/env bash
#SBATCH --job-name=caravan-test
#SBATCH --output=caravan-test.out
#SBATCH --time=00:01:00
#SBATCH --ntasks=1

echo "Hello from Caravan job on $(hostname)"
sleep 10

echo "Done"
