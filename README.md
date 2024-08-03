# AWS Lambda Go-lang Image Resizer

A lightweight go-lang API (with auto-deploy) to use on serverless AWS lambda functions to resize images based on S3 buckets.

## Table of contents

[1. Usage](#usage)

[2. Credits and References](#credits-and-references)

## Usage

```
https://<your-endpoint>/<input-folder>/heightxwidth_crop_name-of-image.ext
```

For example

```
https://kb2sf4lrd3.execute-api.ap-southeast-1.amazonaws.com/production/600x600_center_beeketing.jpg
```

## Credits and References

This repository is based on the work done by [Duc Ho (@ducmeit1)](https://github.com/ducmeit1) at [golang-resize-image-tool](https://github.com/ducmeit1/golang-resize-image-tool). He has also created an awesome tutorial available at [Build a resize images tool with AWS S3](https://medium.com/@ducmeit/build-a-resize-images-tool-with-aws-s3-lambda-api-gateway-at-golang-7569c72c3e8a).
