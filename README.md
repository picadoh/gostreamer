Gostreamer allows one to compose building blocks. A building block can process inputs and generate outputs to the subsequent building blocks through channels. The level of paralellism of each block can also be controlled.

### Building blocks

**Collector** receives a function that collects information and publish it in a channel in the form of a message.

**Processor** picks an input channel with messages and executes some function (passed as parameter) over it

**Demux** demultiplexes the input stream into multiple output streams based on a given number of output channels and a given index function.

### Sample pipeline

Please refer to [Sample Go Pipiline](https://github.com/picadoh/sample-go-pipeline) as an example that uses Gostreamer to read input text from a file or from a socket and processes the words, separating hashtags and counting them.

### Build

    $ go build streamer/*.go

### Run tests

    $ go test test/*.go
