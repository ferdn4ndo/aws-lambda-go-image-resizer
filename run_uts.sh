#!/bin/bash
set -eo pipefail

cd image-resizer
AWS_REGION=$REGION go test

