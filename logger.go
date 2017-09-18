package mx

import (
	"fmt"

	"github.com/mdigger/log3"
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
	var msg = fmt.Sprintf("%s %s", LogINOUT[inFlag], data)
	if id > 0 && id < 9999 {
		c.logger.Debug(msg, "id", fmt.Sprintf("%04d", id))
	} else {
		c.logger.Debug(msg)
	}
	c.mul.RUnlock()
}

// SetLogger устанавливает лог.
func (c *Conn) SetLogger(l log.Logger) {
	c.mul.Lock()
	c.logger = l
	c.mul.Unlock()
}
