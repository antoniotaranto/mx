package mx

import (
	"fmt"
	"testing"

	"github.com/mdigger/log3"
)

func TestConn_Handle(t *testing.T) {
	log.SetLevel(log.DEBUG)
	conn, err := Connect("89.185.246.134:7778")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetLogger(log.New(""))
	if _, err = conn.Login(Login{
		UserName: "peterh",
		Password: "981211",
		Type:     "User",
		Platform: "iPhone",
		Version:  "1.0",
	}); err != nil {
		t.Fatal(err)
	}
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
