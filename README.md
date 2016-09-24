## Golang Channels Pipeline Example

This is an example on how to use Golang channels to build an execution pipeline

### Components
**streamer** builds the main pipeline file using collectors and processors

**scollector** receives a function that collects information and publish it in a channel in the form of a message.

**sprocessor** picks an input channel with messages and executes some function (passed as parameter) over it

**demuxrnd** generates a set of output channels out of a single input channel for concurrent execution based on a random index (e.g. messages are forwarded to a random output channel)

**demuxgrp** generates a set of output channels out of a single input channel for concurrent execution based on a value from the message (all messages with the same indicated value forward to the same output channel)

### How it works?

    ------------         -----------         ---------         -----------
    | Collector | ----> | Extractor | ----> | Counter | ----> | Publisher |
    ------------         -----------         ---------         -----------
         |                                                          |
         |                                                          |
        \|/                                                        \|/
     tweets.txt                                                   sysout


**Input:**

    go #world is #awesome
    #welcome to go #world
    just exploring #golang

**Output:**

    2016/09/24 16:24:29.056739 Loaded configuration: map[parallelism.counter:5 parallelism.publisher:5 source.mode:file source.file:tweets.txt source.port:9999 parallelism.collector:2 parallelism.extractor:2]
    2016/09/24 16:24:29.056966 count report: map[]
    2016/09/24 16:24:29.057010 Read message from file: &{2016-09-24 16:24:29.0568988 +0100 WEST map[tweet:hello #world]}
    2016/09/24 16:24:29.057025 Read message from file: &{2016-09-24 16:24:29.057019608 +0100 WEST map[tweet:#world #golang]}
    2016/09/24 16:24:29.057036 Read message from file: &{2016-09-24 16:24:29.057030369 +0100 WEST map[tweet:the #world is #awesome]}
    2016/09/24 16:24:29.057075 Extracted hashtag #world
    2016/09/24 16:24:29.057088 Extracted hashtag #world
    2016/09/24 16:24:29.057129 Extracted hashtag #golang
    2016/09/24 16:24:29.057133 Read message from file: &{2016-09-24 16:24:29.057092833 +0100 WEST map[tweet:#welcome to my #world]}
    2016/09/24 16:24:29.057169 Publishing #world/1
    2016/09/24 16:24:29.057174 Extracted hashtag #world
    2016/09/24 16:24:29.057185 Read message from file: &{2016-09-24 16:24:29.057178675 +0100 WEST map[tweet:just exploring #golang]}
    2016/09/24 16:24:29.057190 Extracted hashtag #awesome
    2016/09/24 16:24:29.057178 Publishing #golang/1
    2016/09/24 16:24:29.057174 Publishing #world/2
    2016/09/24 16:24:29.057227 Publishing #world/3
    2016/09/24 16:24:29.057229 Extracted hashtag #welcome
    2016/09/24 16:24:29.057239 Extracted hashtag #world
    2016/09/24 16:24:29.057249 Publishing #world/4
    2016/09/24 16:24:29.057256 Extracted hashtag #golang
    2016/09/24 16:24:29.057264 Publishing #welcome/1
    2016/09/24 16:24:29.057273 Publishing #awesome/1
    2016/09/24 16:24:29.057282 Publishing #golang/2
    2016/09/24 16:24:29.057330 final count report: map[#world:%!s(int=4) #golang:%!s(int=2) #awesome:%!s(int=1) #welcome:%!s(int=1)]

### Build

    streamer$ go build tweetpipeline.go

### Configure

Configuration is made to key/value properties in a text file as below:

    # Parallelism
    parallelism.collector = 2
    parallelism.extractor = 2
    parallelism.counter = 5
    parallelism.publisher = 5

    # Source
    # acceptable values for mode: socket, file
    source.mode = file
    source.file = tweets.txt
    source.port = 9999

### Running (File collector)

    streamer$ ./tweetpipeline pipeline.cfg

### Running (Socket collector)

    streamer$ ./tweetpipeline pipeline.cfg

You may want to produce some test messages, a simple way of doing it is by using nc command

    $ nc localhost 9999
    hello #world

### Running tests

    streamer$ go test test/*.go
