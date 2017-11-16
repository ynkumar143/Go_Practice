package main

import (
	"fmt"
	"time"
	"github.com/nats-io/go-nats"
	"log"
)

type userID string
type method string
type topic string

type topicChannel chan userInfo
type userChannel string

type subscriptions map[userID]string

type message struct {
	method method
	topic   topic
	user  userID
}

type userInfo struct {
	ID      userID
	Channel userChannel
}

//This program explains usage of storing subscription details and nats package for publishing messages.


func PerformSubscription(topics map[topic]topicChannel) chan message {

	//This route will store all subscriptions information based on topic.

	routes := make(map[topic]subscriptions)

	for topic := range topics {
		routes[topic] = make(subscriptions)
	}

	fmt.Println("Route Structure Before Subscription Begins", routes)
	messageChan := make(chan message)

	go func() {
		for msg := range messageChan {
			var latest_topic string
			latest_topic = string(msg.topic) + "." + string(msg.user)

			switch string(msg.method) {
			case "Unsubscribe":
				// delete the client from the topic
				delete(routes[msg.topic], msg.user)
				fmt.Println("Unsubscribing Topic: X ", latest_topic)
				fmt.Println("Current Status of subscription is:\n", routes)
				fmt.Println(" ")
			case "Subscribe":
				// add the client to the topic
				nc , nats_err := nats.Connect(nats.DefaultURL)
				if nats_err != nil {
					log.Println("Error Connecting to " + nats.DefaultURL , nats_err)
				}

				nc.Subscribe(latest_topic, func(m *nats.Msg){
					routes[msg.topic][msg.user] = string(m.Data)
					fmt.Println("Message Received  on Topic: <--", m.Subject, "\nMessage:", string(m.Data))
					fmt.Println("Current subscription list:\n", routes)
					fmt.Println(" ")
				})

			case "Publish":
				nc , _ := nats.Connect(nats.DefaultURL)
				nc.Publish(latest_topic,[]byte("Hello Subscribers"))
				fmt.Println("Message Published on Topic: -->", latest_topic)
			}
		}
	}()

	return messageChan
}


func main() {

	topics := make(map[topic]topicChannel)

	topics["Topic1"] = make(topicChannel)
	topics["Topic2"] = make(topicChannel)

	Chan := PerformSubscription(topics)

	//Subscribing to the topics
	Chan <- message{method: "Subscribe",topic:   topic("Topic1"),	user:  "User1",}
	Chan <- message{method: "Subscribe",topic:   topic("Topic2"),	user:  "User2",}
	Chan <- message{method: "Subscribe",topic:   topic("Topic2"),	user:  "User1",}
	Chan <- message{method: "Subscribe",topic:   topic("Topic1"),	user:  "User2",}

	<-time.After(5 * time.Second)

	//Publishing to the Topics
	Chan <- message{method: "Publish",	topic:   topic("Topic1"),		user:  "User1",}

	<-time.After(5 * time.Second)  //Wait for 5 seconds (Until it is published in Go routine)

	Chan <- message{method: "Publish",	topic:   topic("Topic2"),		user:  "User2",}

	<-time.After(5 * time.Second)

	Chan <- message{method: "Publish",	topic:   topic("Topic1"),		user:  "User2",}

	<-time.After(5 * time.Second)

	fmt.Println("==========================================================")
	fmt.Println("               Unsubscribing the list started.            ")
	fmt.Println("==========================================================")


	Chan <- message{method: "Unsubscribe",topic:   topic("Topic1"),	user:  "User1",	}
	Chan <- message{method: "Unsubscribe",topic:   topic("Topic1"),	user:  "User2",	}
	Chan <- message{method: "Unsubscribe",topic:   topic("Topic2"),	user:  "User1",	}
	Chan <- message{method: "Unsubscribe",topic:   topic("Topic2"),	user:  "User2",	}

	<-time.After(5 * time.Second)
}
