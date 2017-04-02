#!/bin/bash -ex

docker pull golang:onbuild
docker build -t localhost:5000/benlubar/wtdwtf-science .
docker push localhost:5000/benlubar/wtdwtf-science
docker service update --image localhost:5000/benlubar/wtdwtf-science wtdwtf-science
