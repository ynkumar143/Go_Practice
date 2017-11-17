package main

import (
	"github.com/nats-io/go-nats"
	"log"
	"fmt"
	"runtime"
	"time"
	"reflect"
)

type natConnect struct{
	nats_connct *nats.Conn
	unsub_req bool
}

type usersubinfo struct{
	FirstName string
	LastName  string
	InfoApi   string
	subs      *nats.Subscription
}

func main(){
	//Connecting to NATS Package and storing nats connection information in nc.
	nc , err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, nats.DefaultURL)
	}

	//savednatConn is holding the connection information of NATS
	savednatConnect := &natConnect{
		nats_connct:nc,
	}

	//I am having user whose details are as below. This user needed nats published message into Api given in InfoApi
	userOneInfo := &usersubinfo{
		FirstName: "Tree",
		LastName: "Branch",
		InfoApi: "http://localhost:8080/FirstApi",
		subs:nil,
	}

	SubscribeInfoAndSave(savednatConnect, userOneInfo)

	fmt.Println("SID Of First Subscription:", reflect.Indirect(reflect.ValueOf(savednatConnect.nats_connct).Elem().Field(12)))
	// Here only difference is Api url is given different
	userTwoInfo := &usersubinfo{
		FirstName: "Tree",
		LastName: "Branch",
		InfoApi: "http://localhost:8001/SecondApi",
		subs:nil,
	}

	<-time.After(1 * time.Second)

	SubscribeInfoAndSave(savednatConnect, userTwoInfo)

	fmt.Println("SID Of Second Subscription:", reflect.Indirect(reflect.ValueOf(savednatConnect.nats_connct).Elem().Field(12)))
	<-time.After(2 * time.Second)

	//I will send a message to the topic FirstName + LastName
	PublishMessages(savednatConnect)

	//Message Publishing is happening properly. Order is not important it should that message to both API's.

	//Here I am not having limit of messages received for each subscription so I can't use AUTOUNSUBSCRIBE
	//I need this messages continuously and i cant give unsubscribe once the message is delivered.

	<-time.After(2 * time.Second)

	fmt.Println("Subscription Information Before Unsubscribing", reflect.Indirect(reflect.ValueOf(savednatConnect.nats_connct).Elem().Field(14)))

	//Now Subscribe for 2nd User

	userTwoInfo.subs.Unsubscribe()

	<-time.After(2 * time.Second)

	fmt.Println("After Unsubscribing 2nd user", reflect.Indirect(reflect.ValueOf(savednatConnect.nats_connct).Elem().Field(14)))

	userOneInfo.subs.Unsubscribe()

	<-time.After(2 * time.Second)

	fmt.Println("After Unsubscribing 1st user", reflect.Indirect(reflect.ValueOf(savednatConnect.nats_connct).Elem().Field(14)))

	runtime.Goexit()
}


func SubscribeInfoAndSave(recvdnatConn *natConnect, recvduserInfo *usersubinfo){
	//I want the topic to be as FirstName + Last Name
	latest_topic := recvduserInfo.FirstName + "." + recvduserInfo.LastName

	handle := func(m *nats.Msg) {
		fmt.Println("Message Received  on Topic: <--", m.Subject, "\nMessage:", string(m.Data))

		fmt.Println("To URL:", recvduserInfo.InfoApi)
		<-time.After(1 * time.Second)
		// Code for the message received will be sent to API.
		// As this handle hits two times on publishing to that topic the message will send to two API's
	}

	fmt.Println("Subscribing for Topic:", latest_topic)
	nc_info , nats_sub_err := recvdnatConn.nats_connct.Subscribe(latest_topic, handle)
	recvduserInfo.subs = nc_info;
	if nats_sub_err != nil{
		recvdnatConn.nats_connct.Close()
		log.Fatal(nats_sub_err)
	}
}

func PublishMessages(recvdnatConn *natConnect){
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