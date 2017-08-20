package mx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mdigger/log"
)

var JSON func(v interface{}) error

func init() {
	SetCSTALog(os.Stdout, 0)
	log.SetLevel(log.DebugLevel)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	JSON = enc.Encode
}

func TestAddressbook(t *testing.T) {
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
	JSON(conn.Info)

	contacts, err := conn.Addressbook()
	if err != nil {
		t.Fatal(err)
	}

	JSON(contacts)

}
