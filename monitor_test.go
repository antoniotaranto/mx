package mx

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	KeepAliveDuration = time.Second * 5
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

	id, err := conn.MonitorStart("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("monitor:", id)

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

	// разбор ответов сервера
	var callLog []*CallInfo
	if err := conn.Handle(func(resp *Response) error {
		switch resp.Name {
		case "callloginfo":
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
		case "FaxBoxUpToDate":
			var faxBox = new(struct {
				Box string `xml:"box,attr"`
			})
			if err = resp.Decode(faxBox); err != nil {
				return err
			}
			if faxBox.Box == "Sentbox" {
				return Stop
			}
		}
		return nil
	}, "callloginfo", "FaxBoxUpToDate"); err != nil {
		t.Error(err)
	}

	if err := conn.MonitorStop(0); err != nil {
		t.Fatal(err)
	}

	JSON(callLog)
}

// CallInfo описывает информацию о записи в логе звонков.
type CallInfo struct {
	Missed                bool   `xml:"missed,attr" json:"missed,omitempty"`
	Direction             string `xml:"direction,attr" json:"direction"`
	RecordID              uint32 `xml:"record_id" json:"recordId"`
	GCID                  string `xml:"gcid" json:"gcid"`
	ConnectTimestamp      int64  `xml:"connectTimestamp" json:"connect,omitempty"`
	DisconnectTimestamp   int64  `xml:"disconnectTimestamp" json:"disconnect,omitempty"`
	CallingPartyNo        string `xml:"callingPartyNo" json:"callingPartyNo"`
	OriginalCalledPartyNo string `xml:"originalCalledPartyNo" json:"originalCalledPartyNo"`
	FirstName             string `xml:"firstName" json:"firstName,omitempty"`
	LastName              string `xml:"lastName" json:"lastName,omitempty"`
	Extension             string `xml:"extension" json:"ext,omitempty"`
	ServiceName           string `xml:"serviceName" json:"serviceName,omitempty"`
	ServiceExtension      string `xml:"serviceExtension" json:"serviceExtension,omitempty"`
	CallType              uint32 `xml:"callType" json:"callType,omitempty"`
	LegType               uint32 `xml:"legType" json:"legType,omitempty"`
	SelfLegType           uint32 `xml:"selfLegType" json:"selfLegType,omitempty"`
	MonitorType           uint32 `xml:"monitorType" json:"monitorType,omitempty"`
}
