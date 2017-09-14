package mx

import (
	"fmt"

	"github.com/mdigger/log"
)

// LogINOUT задает символы, используемые для вывода направления
// (true - входящие, false - исходящие)
var LogINOUT = map[bool]string{true: "→", false: "←"}

// csta форматирует вывод лога с командами CSTA.
func (c *Conn) csta(inFlag bool, id uint16, data []byte) {
	c.mul.RLock()
	if c.logger == nil {
		c.mul.RUnlock()
		return
	}
	var ctxLog = c.logger
	if id > 0 && id < 9999 {
		ctxLog = ctxLog.WithField("id", fmt.Sprintf("%04d", id))
	}
	ctxLog.Debugf("%s %s", LogINOUT[inFlag], data)
	c.mul.RUnlock()
}

// SetLogger устанавливает лог.
func (c *Conn) SetLogger(l *log.Context) {
	c.mul.Lock()
	c.logger = l
	c.mul.Unlock()
}
