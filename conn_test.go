package mx

import (
	"fmt"
	"testing"
	"time"
)

func TestDone(t *testing.T) {
	conn, err := Connect("89.185.246.134:7778", Login{
		UserName: "peterh",
		Password: "981211",
		Type:     "User",
		Platform: "iPhone",
		Version:  "1.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Close error:", err.Error())
		}
	}()

	go func() {
		<-conn.Done()
		fmt.Println("Connection done:", conn.Error())
		conn.Close()
		fmt.Println("Connection done 2:", conn.Error())
	}()

	for {
		if err = conn.Send("<keepalive/>"); err != nil {
			fmt.Println("Send error:", err.Error())
			break
		}
		time.Sleep(time.Second * 20)
	}
	fmt.Println("All finished:", conn.Error())
}
