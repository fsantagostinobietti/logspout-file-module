FROM gliderlabs/logspout:master

ADD custom/modules.go modules.go
ADD custom/build.sh build.sh

