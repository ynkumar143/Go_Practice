package main

import "fmt"

type User struct {
	FirstName, LastName string
}

func (u *User) Name() string{
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
type Namer interface{
	Name() string
}

func Greet(n Namer) string{
	return fmt.Sprintf("Dear %s", n.Name())
}

func main(){
	u := &User{"Nagendra", "Kumar"}
	fmt.Println(Greet(u))
}