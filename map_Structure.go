package main

import "fmt"

func main(){


	subscriptions := make(map[string]string)
	subscriptions["Topic1"] = "User1"
	subscriptions["Topic1"] = "User2"
	subscriptions["Topic1"] = "User3"
	subscriptions["Topic2"] = "User1"
	subscriptions["Topic3"] = "User2"
	subscriptions["Topic2"] = "User3"

	fmt.Println(subscriptions)


}