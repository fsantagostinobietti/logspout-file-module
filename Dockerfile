FROM gliderlabs/logspout:master

ADD modules.go modules.go
ADD build.sh build.sh

#ENV SYSLOG_FORMAT rfc3164
