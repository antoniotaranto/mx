package mx

import (
	"encoding/xml"
	"testing"
)

func TestHandler(t *testing.T) {
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

	// запрос на получение лога звонков
	err = conn.Send(&struct {
		XMLName   xml.Name `xml:"iq"`
		Type      string   `xml:"type,attr"`
		ID        string   `xml:"id,attr"`
		Timestamp int64    `xml:"timestamp,attr"`
	}{
		Type:      "get",
		ID:        "calllog",
		Timestamp: -1,
	})
	if err != nil {
		t.Error(err)
	}

	var callLog []*CallInfo
	err = conn.HandleWait(func(resp *Response) error {
		var items = new(struct {
			LogItems []*CallInfo `xml:"callinfo"`
		})
		if err = resp.Decode(items); err != nil {
			return err
		}
		if callLog == nil {
			callLog = items.LogItems
		} else {
			callLog = append(callLog, items.LogItems...)
		}
		// fmt.Println("size:", len(items.LogItems))
		if len(items.LogItems) < 21 {
			// fmt.Println("-- ! --")
			return Stop // заканчиваем обработку
		}
		return nil // ждем следующего ответа
	}, ReadTimeout, "callloginfo")
	if err != nil {
		if err != ErrTimeout {
			t.Error(err)
		}
	}
	// fmt.Println("finish")
}
