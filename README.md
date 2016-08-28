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

    2016/08/28 13:16:04 Generated message &{2016-08-28 13:16:04.386522533 +0100 WEST map[tweet:go #world is #awesome]}
    2016/08/28 13:16:04 Generated message &{2016-08-28 13:16:04.386753983 +0100 WEST map[tweet:#welcome to go #world]}
    2016/08/28 13:16:04 Extracted hashtag #welcome
    2016/08/28 13:16:04 Extracted hashtag #world
    2016/08/28 13:16:04 Extracted hashtag #world
    2016/08/28 13:16:04 Counted #world/1
    2016/08/28 13:16:04 Extracted hashtag #awesome
    2016/08/28 13:16:04 Publishing #world/1
    2016/08/28 13:16:04 Counted #awesome/1
    2016/08/28 13:16:04 Publishing #awesome/1
    2016/08/28 13:16:04 Counted #world/2
    2016/08/28 13:16:04 Publishing #world/2
    2016/08/28 13:16:04 Generated message &{2016-08-28 13:16:04.386777564 +0100 WEST map[tweet:just exploring #golang]}
    2016/08/28 13:16:04 Counted #welcome/1
    2016/08/28 13:16:04 Extracted hashtag #golang
    2016/08/28 13:16:04 Publishing #welcome/1
    2016/08/28 13:16:04 Counted #golang/1
    2016/08/28 13:16:04 Publishing #golang/1
    2016/08/28 13:16:04 report: map[#awesome:%!s(int=1) #golang:%!s(int=1) #world:%!s(int=2) #welcome:%!s(int=1)]

### Build
    streamer$ go build -o streamer src/*.build
    streamer$ ./streamer
