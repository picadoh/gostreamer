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

    2016/08/28 14:05:58 Generated message &{2016-08-28 14:05:58.919532613 +0100 WEST map[tweet:the #world is #awesome]}
    2016/08/28 14:05:58 Generated message &{2016-08-28 14:05:58.919722092 +0100 WEST map[tweet:#welcome to my #world]}
    2016/08/28 14:05:58 Extracted hashtag #welcome
    2016/08/28 14:05:58 Extracted hashtag #world
    2016/08/28 14:05:58 Counted #world/1
    2016/08/28 14:05:58 Publishing #world/1
    2016/08/28 14:05:58 Extracted hashtag #world
    2016/08/28 14:05:58 Generated message &{2016-08-28 14:05:58.919751818 +0100 WEST map[tweet:just exploring #golang]}
    2016/08/28 14:05:58 Extracted hashtag #awesome
    2016/08/28 14:05:58 Counted #world/2
    2016/08/28 14:05:58 Publishing #world/2
    2016/08/28 14:05:58 Counted #awesome/1
    2016/08/28 14:05:58 Extracted hashtag #golang
    2016/08/28 14:05:58 Publishing #awesome/1
    2016/08/28 14:05:58 Counted #golang/1
    2016/08/28 14:05:58 Publishing #golang/1
    2016/08/28 14:05:58 Counted #welcome/1
    2016/08/28 14:05:58 Publishing #welcome/1
    2016/08/28 14:05:58 report: map[#welcome:%!s(int=1) #awesome:%!s(int=1) #golang:%!s(int=1) #world:%!s(int=2)]

### Build
    streamer$ go build tweetpipeline.go

### Running main
    streamer$ ./tweetpipeline tweets.txt

### Running tests
    streamer$ go test test/*.go
