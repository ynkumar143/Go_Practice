package main

import (
	"fmt"
	"time"
	"math/rand"
)


/*
 This example is explaining about usage of Channel Data type and Go routine functionality.

 */
func main(){
	my_channel := make(chan string)       //Create a channel

	//Go Routine function that will call the procedure parallely.

	go SampleRoutine("Testing the Channel Message", my_channel)

	for i := 0; i < 5; i++ {
		fmt.Printf("Message received for me: %q \n", <- my_channel)
	}

	fmt.Println("Finished processing the messaage")
}

func SampleRoutine(Message string, c chan string){
	for i := 0 ; ; i++{
		// Channel receives formatted message.
		c <- fmt.Sprintf("%s %d", Message , i) //Expression to be sent can be any suitable value

		//This command will wait until the message is delivered properly in the channel.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

/*  OUTPUT
Message received for me: "Testing the Channel Message 0"
Message received for me: "Testing the Channel Message 1"
Message received for me: "Testing the Channel Message 2"
Message received for me: "Testing the Channel Message 3"
Message received for me: "Testing the Channel Message 4"
Finishe processing the messaage


 */