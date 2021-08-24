#!/bin/bash
#docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --use --name mybuilder
docker buildx build -t scjtqs/jd_cookie:latest -f Dockerfile --build-arg Version="v3.0.7" --platform linux/amd64,linux/arm64,linux/386,linux/arm/v7 --push .
docker buildx rm mybuilder
