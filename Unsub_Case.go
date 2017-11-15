package main

import (
	"github.com/nats-io/go-nats"
	"log"
	"fmt"
	"runtime"
	"time"
	"strings"
	"reflect"
)

const (
	Nats_URL    = "nats://localhost:4222"
)

type natConn struct{
	nats_connct *nats.Conn
	unsub_req bool
}

type userinfo struct{
	FirstName string
	LastName  string
	InfoApi   string
}

func main(){
	//Connecting to NATS Package and storing nats connection information in nc.
	nc , err := nats.Connect(Nats_URL)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, Nats_URL)
	}

	//savednatConn is holding the connection information of NATS
	savednatConn := &natConn{
		nats_connct:nc,
	}

	//I am having user whose details are as below. This user needed nats published message into Api given in InfoApi
	userOneInfo := &userinfo{
		FirstName: "Tree",
		LastName: "Branch",
		InfoApi: "http://localhost:8080/FirstApi",
	}

	SubscribeInfo(savednatConn, userOneInfo)

	fmt.Println("SID Of First Subscription:", reflect.Indirect(reflect.ValueOf(savednatConn.nats_connct).Elem().Field(12)))
	// Here only difference is Api url is given different
	userTwoInfo := &userinfo{
		FirstName: "Tree",
		LastName: "Branch",
		InfoApi: "http://localhost:8001/SecondApi",
	}

	<-time.After(1 * time.Second)

	SubscribeInfo(savednatConn, userTwoInfo)

	fmt.Println("SID Of Second Subscription:", reflect.Indirect(reflect.ValueOf(savednatConn.nats_connct).Elem().Field(12)))

	<-time.After(2 * time.Second)

	//I will send a message to the topic FirstName + LastName
	PublishMessage(savednatConn)

	//Message Publishing is happening properly. Order is not important it should that message to both API's.

	//Here I am not having limit of messages received for each subscription so I can't use AUTOUNSUBSCRIBE
	//I need this messages continuously and i cant give unsubscribe once the message is delivered.

	<-time.After(2 * time.Second)

	fmt.Println("Subscription Information Before Unsubscribing", reflect.Indirect(reflect.ValueOf(savednatConn.nats_connct).Elem().Field(14)))

	//Here I want to Unsubscribe to the Topic(i.e First Subscription of two subscriptions with same topic).

	UnsubscribeTopic(savednatConn, userOneInfo)

	<-time.After(2 * time.Second)

	fmt.Println("Subscription Information after Unsubscribing", reflect.Indirect(reflect.ValueOf(savednatConn.nats_connct).Elem().Field(14)))

	runtime.Goexit()
}


func SubscribeInfo(recvdnatConn *natConn, recvduserInfo *userinfo){
	//I want the topic to be as FirstName + Last Name
	latest_topic := recvduserInfo.FirstName + "." + recvduserInfo.LastName

	handle := func(m *nats.Msg) {
			if recvdnatConn.unsub_req{
				if strings.EqualFold(m.Reply, "Unsubscribe" + recvduserInfo.InfoApi){
					m.Sub.Unsubscribe()
					recvdnatConn.unsub_req = false
				}
			return
		}
		fmt.Println("Message Received  on Topic: <--", m.Subject, "\nMessage:", string(m.Data))

		fmt.Println("To URL:", recvduserInfo.InfoApi)
		<-time.After(1 * time.Second)
		// Code for the message received will be sent to API.
		// As this handle hits two times on publishing to that topic the message will send to two API's
	}

	fmt.Println("Subscribing for Topic:", latest_topic)
	_ , nats_sub_err := recvdnatConn.nats_connct.Subscribe(latest_topic, handle)
	if nats_sub_err != nil{
		recvdnatConn.nats_connct.Close()
		log.Fatal(nats_sub_err)
	}
}

func PublishMessage(recvdnatConn *natConn){
	published_topic := "Tree.Branch"
	fmt.Println("Message Publishing on Topic:", published_topic)
	//Building Message Structure

	mes := nats.Msg{}
	mes.Subject = published_topic
	mes.Data = []byte("Information")

	// Publish message on subject to Unsubscribe
	if err := recvdnatConn.nats_connct.PublishMsg(&mes); err != nil {
		log.Println(err)
	}
}

func UnsubscribeTopic(recvdnatConn *natConn, recvduserinfo *userinfo){
	//Here to access Subscription details I have to message that is handling the subscription.

	//So I am doing a publish message with reply message as unsubscribe.

	unsubscribe_topic := "Tree.Branch"

	fmt.Println("Unsubscribing from Topic", unsubscribe_topic)
	//Building Message Structure

	mes := nats.Msg{}
	mes.Subject = unsubscribe_topic
	mes.Reply = "Unsubscribe" + recvduserinfo.InfoApi
	recvdnatConn.unsub_req = true

	// Publish message on subject to Unsubscribe
	// If I am doing this the both subscription information is getting erased.

	if err := recvdnatConn.nats_connct.PublishMsg(&mes); err != nil {
		log.Println(err)
	}
}
