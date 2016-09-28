package main

import (
	"log"
	"strings"
	"github.com/picadoh/gostreamer/streamer"
	"os"
	"fmt"
	"net"
	"bufio"
)

var countState streamer.Counter = *streamer.NewCounter()

func TextFileCollector(name string, cfg streamer.Config, out chan streamer.Message) {
	lines, _ := streamer.LoadTextFile(cfg.GetString("source.file"))

	for _, line := range lines {
		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Read message from file: %s\n", out_message)
		out <- out_message
	}
}

func TextSocketCollector(name string, cfg streamer.Config, out chan streamer.Message) {
	listener, _ := net.Listen("tcp", ":" + cfg.GetString("source.port"))
	conn, _ := listener.Accept()

	for {
		line, _ := bufio.NewReader(conn).ReadString('\n')
		line = strings.TrimSuffix(line, "\n")

		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)

		log.Printf("Received raw message from socket: %s\n", out_message)

		out <- out_message
	}
}

func WordExtractor(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
	tweet, _ := input.Get("tweet").(string)

	words := strings.Split(tweet, " ")

	for _, word := range words {
		out_message := streamer.NewMessage()
		out_message.Put("word", word)
		log.Printf("Extracted word: %s\n", word)
		out <- out_message
	}
}

func HashTagFilter(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
	word, _ := input.Get("word").(string)

	if (strings.HasPrefix(word, "#")) {
		out_message := streamer.NewMessage()
		out_message.Put("hashtag", word)
		log.Printf("Filtered hashtag %s\n", word)
		out <- out_message
	}
}

func HashTagCounter(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
	hashtag := input.Get("hashtag").(string)

	count := countState.Increment(hashtag)

	out_message := streamer.NewMessage()
	out_message.Put("hashtag", hashtag)
	out_message.Put("count", count)

	out <- out_message
}

func HashTagCountPublisher(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
	hashtag, _ := input.Get("hashtag").(string)
	count, _ := input.Get("count").(int)

	log.Printf("Publishing %s/%d\n", hashtag, count)
}

func RunPipeline(cfg streamer.Config) {
	// read config
	var extractorParallelismHint = cfg.GetInt("parallelism.extractor")
	var filterParallelismHint = cfg.GetInt("parallelism.filter")
	var counterParallelismHint = cfg.GetInt("parallelism.counter")
	var publisherParallelismHint = cfg.GetInt("parallelism.publisher")

	// build pipeline elements
	collector := streamer.NewCollector("collector", cfg,
		defineCollectorFunction(cfg))

	extractor := streamer.NewProcessor("extractor", cfg,
		WordExtractor, streamer.NewDemux(extractorParallelismHint, streamer.NewRandomDemuxCtx()))

	filter := streamer.NewProcessor("filter", cfg,
		HashTagFilter, streamer.NewDemux(filterParallelismHint, streamer.NewRandomDemuxCtx()))

	counter := streamer.NewProcessor("counter", cfg,
		HashTagCounter, streamer.NewDemux(counterParallelismHint, streamer.NewGroupDemuxCtx("hashtag")))

	publisher := streamer.NewProcessor("publisher", cfg,
		HashTagCountPublisher, streamer.NewDemux(publisherParallelismHint, streamer.NewGroupDemuxCtx("hashtag")))

	// execute pipeline
	sequence := collector.Execute()
	extracted := extractor.Execute(sequence)
	filtered := filter.Execute(extracted)
	counted := counter.Execute(filtered)
	<-publisher.Execute(counted)

	// print final report
	log.Printf("final count report: %s\n", countState.ToString())
}

func defineCollectorFunction(cfg streamer.Config) streamer.CollectFunction {
	switch cfg.GetString("source.mode") {
	case "file": return TextFileCollector
	case "socket": return TextSocketCollector
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if (len(os.Args) < 2) {
		fmt.Println("Usage: " + os.Args[0] + " [<path/to/config/file>]")
		os.Exit(2)
	}

	var cfg, err = streamer.LoadProperties(os.Args[1])

	if (err != nil) {
		fmt.Printf("An error ocurred reading the properties file %s [cause: %s]\n", os.Args[1], err)
		os.Exit(2)
	}

	log.Printf("Loaded configuration: %s\n", cfg.ToString())

	RunPipeline(cfg)
}
