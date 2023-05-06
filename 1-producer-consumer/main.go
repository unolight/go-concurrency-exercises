//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return tweets
		}

		tweets = append(tweets, tweet)
	}
}

func producerByChan(stream Stream, ch chan<- *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			break
			// return tweets
		}

		// tweets = append(tweets, tweet)
		ch <- tweet
	}
}

func consumer(tweets []*Tweet) {
	for _, t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func consumerByChan(ch <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweetChan := make(chan *Tweet, 100)
	tweetWg := sync.WaitGroup{}

	// Producer
	// tweets := producer(stream)
	go producerByChan(stream, tweetChan)
	tweetWg.Add(1)

	// Consumer
	// go consumer(tweets)
	go consumerByChan(tweetChan, &tweetWg)

	tweetWg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
