#!/bin/bash
set -eo pipefail

echo "Setting up variables"
REGION=$(aws configure get region)

echo "Running UTs"
./run_uts.sh

echo "Running ATs"
./run_ats.sh

echo "Finished running all tests!"
