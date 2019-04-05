Gostreamer allows one to compose building blocks in a processing pipeline. A building block can process inputs and generate outputs to the subsequent building blocks through channels. The level of paralellism of each block can also be controlled.

## How does it work?

### Configuration

Configuration is a simple map context that can be used to pass dynamic information to be used by the functions. It can be created by loading the properties from a properties file (file representing each property by a <key>=<value> format).

	cfg := streamer.LoadProperties("samplepipeline.properties")

The samplepipeline.properties file could look as follows:

	greeting=hello world

### Collectors

A collector is a component responsible for gather information from a specific source and publish it into a channel in the form of a keyed message.

**Example:**

	func TextCollector(name string, cfg streamer.Config, out chan streamer.Message) {
		out_message := streamer.NewMessage()
		out_message.Put("greeting", cfg.GetString("greeting"))
		out <- out_message
	}

The above function publishes static "hello world" message keyed with "greeting". With this signature, this function can be used to build a collector, as below:

	collector := streamer.NewCollector("collector", cfg, TextCollector)

### Processors

A processor is responsible for consuming keyed messages from an input channel, do some processing and possibly publish more keyed messages into an output channel. The processor delivers each message from the input channel into a function for processing.

**Example:**

	func WordExtractor(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
		text, _ := input.Get("greeting").(string)

		words := strings.Split(text, " ")

		for _, word := range words {
			out_message := streamer.NewMessage()
			out_message.Put("word", word)
			out <- out_message
		}
	}

The above function picks up messages keyed with "greeting" and splits the message by the single whitespace delimiter, then it publishes a single word to the output channel as a messaged keyed by "word". With this signature this function can be used to build a processor, as below:

	extractor := streamer.NewProcessor("extractor", cfg, WordExtractor, <Demux>)

Another processor could be used to print each individual message to the output:

	func WordPrinter(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
		word, _ := input.Get("word").(string)

		// simulate some processing time
		time.Sleep(1 * time.Second)

		log.Println(word)
	}

### Demux

When creating a processor, one of the arguments is a Demux. The Demux is a special component that allows to build concurrent work inside a processor. It picks the processor input channel and demultiplexes into  multiple output channels that will be each processed by a separate routine. A demux can be created as follows:

	demux := streamer.NewIndexedChannelDemux(2, streamer.RandomIndex)

A Demux receives a parallelism hint. If possible, it will be run in parallel, depending on the parallelism that can be achieved in the underlying system.

The indexed channel demux is a default implementation that creates an array of output channels. The first argument is the parallelism hint, i.e, the number of channels and routines that will be created for each individual message picked up from the input. The second argument is a function that maps the input to a specific output channel.

This function should respect the following signature:

	func <name>(fanOut int, input streamer.Message) int

The default streamer.RandomIndex functions randomly selects the output channel.

An example of a custom static mapping could be:

	func StaticIndex(fanOut int, input streamer.Message) int {
		switch input.Get("word").(string) {
		case "hello":
			return 0
		default:
			return 1
		}
	}

The above function gets an input message keyed with word and routes the word hello to channel at index 0 and every other word to channel at index 1.

### Building the topology

The topology can be built by chaining the multiple components together, as in the following example:

	// build the components
	cfg, _ := streamer.LoadProperties("samplepipeline.properties")
	collector := streamer.NewCollector("collector", cfg, TextCollector)
	extractor := streamer.NewProcessor("extractor", cfg, WordExtractor, streamer.NewIndexedChannelDemux(1, streamer.RandomIndex))
	printer := streamer.NewProcessor("printer", cfg, WordPrinter, streamer.NewIndexedChannelDemux(2, StaticIndex))

	// execute pipeline
	<-printer.Execute(
		extractor.Execute(
			collector.Execute()))

## Sample Pipeline

Please refer to [Sample Go Pipiline](https://github.com/picadoh/sample-go-pipeline) as full-running example that uses Gostreamer to read input text from a file or from a socket and processes the words, separating hashtags and counting them.

## Build

    $ go build streamer/*.go

## Run tests

    $ go test test/*.go
