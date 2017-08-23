package mx

import (
	"fmt"
	"os"
	"testing"
)

func TestConn_Handle(t *testing.T) {
	SetCSTALog(os.Stdout, Lcolor)
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
	defer conn.Close()
	conn.MonitorStart("")
	defer conn.MonitorStopID(0)

	go func() {
		err = conn.Handle(func(resp *Response) error {
			fmt.Println("event:", resp.String())
			return nil
		}, "presence")
		if err != nil {
			t.Error(err)
		}
	}()

	<-conn.Done()
	fmt.Println("finish")
}
