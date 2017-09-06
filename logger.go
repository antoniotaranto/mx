package mx

import (
	"fmt"
)

// LogINOUT задает символы, используемые для вывода направления
// (true - входящие, false - исходящие)
var LogINOUT = map[bool]string{true: "→", false: "←"}

// csta форматирует вывод лога с командами CSTA.
func (c *Conn) csta(inFlag bool, id uint16, data []byte) {
	if c.Logger == nil {
		return
	}
	var ctxLog = c.Logger
	if id > 0 && id < 9999 {
		ctxLog = ctxLog.WithField("id", fmt.Sprintf("%04d", id))
	}
	ctxLog.Debugf("%s %s", LogINOUT[inFlag], data)
}
