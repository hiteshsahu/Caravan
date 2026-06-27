#!/usr/bin/env bash
#SBATCH --job-name=caravan-gpu-test
#SBATCH --output=caravan-gpu-test.out
#SBATCH --time=00:01:00
#SBATCH --ntasks=1
#SBATCH --gres=gpu:1

echo "Hello from Caravan GPU job on $(hostname)"
nvidia-smi --query-gpu=name,memory.total,memory.used --format=csv
echo "CUDA_VISIBLE_DEVICES=$CUDA_VISIBLE_DEVICES"

echo "Done"
