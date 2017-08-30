# logspout-file
A minimalistic adapter for github.com/gliderlabs/logspout to write to file.
It is based on instruction in [https://github.com/gliderlabs/logspout/tree/master/custom].

It can be useful if you want to use it in conjunction with Splunk for application metrics.

## features
 - it uses a log rotation strategy to save logs: file is renamed when it reaches a custom defined size
 - provided as a small docker container
 
## build docker image
```
$ docker build -t logspout-file:latest .
```

## usage example
```
# start 'logspout-file' instance 
$ docker run -d --name="logspout-file" --volume=/var/run/docker.sock:/var/run/docker.sock -v $(pwd)/log:/var/log logspout-file:latest file://filename.log?maxfilesize=10240

# start applications you want to collect logs
$ docker run ...
```

If your applications produce log lines you'll see them in mounted volume (in this example '$(pwd)/log').
```
$ la -1 $(pwd)/log
filename.log
filename.log.2017-08-29T11:49:12Z
```

## custum options
You can customize some features using options:

option |  description   | default value
---------|----------------|--------------
maxfilesize | max size of rotated file | 100 Mbyte

