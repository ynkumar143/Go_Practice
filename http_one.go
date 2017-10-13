package main

import "context"

type message struct{
	ctx context.Context
}

func ProcessMessage(work <- chan message){
	for job := range work{
		select{
			case <- job.ctx.Done();
		}
	}
}
func main(){
	q := make(chan message)
	go ProcessMessage(q)
	ctx := context.Background()

}
