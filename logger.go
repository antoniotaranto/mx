package mx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// LogINOUT задает символы, используемые для вывода направления
// (true - входящие, false - исходящие)
var LogINOUT = map[bool]string{true: "→", false: "←"}

// Lcolor является флагом принудительного вывода в лог в цвете.
const Lcolor = 1 << 7

var (
	// флаг вывода в консоль в цвете
	// устанавливается автоматически при задании SetLogOutput, если
	// поддерживается ASCII
	logTTY = false
	// лог, используемый для вывода команд CSTA
	cstaLogger = log.New(ioutil.Discard, "", log.LstdFlags)
	output     = false
)

// SetCSTALog задает куда выводить лог с командами CSTA и в каком виде.
func SetCSTALog(w io.Writer, flag int) {
	cstaLogger.SetOutput(w)
	output = (w != ioutil.Discard)
	if out, ok := w.(*os.File); ok {
		if fi, err := out.Stat(); err == nil {
			logTTY = fi.Mode()&(os.ModeDevice|os.ModeCharDevice) != 0
		}
	}
	cstaLogger.SetFlags(flag)
	// взводим вывод лога в цвете
	if Lcolor&flag != 0 {
		logTTY = true
	}
}

// csta форматирует вывод лога с командами CSTA.
func (c *Conn) csta(inFlag bool, id uint16, data []byte) {
	if !output {
		return
	}
	var fmtID = "%04d"
	if logTTY {
		// добавляем цветовое выделение к идентификатору команды или ответа
		switch id {
		case 0:
			fmtID = "\033[37m" + "%04d" + "\033[0m"
		case 9999:
			fmtID = "\033[34m" + "%04d" + "\033[0m"
		default:
			fmtID = "\033[33m" + "%04d" + "\033[0m"
		}
		// выделяем цветом название команды или ответа
		indx := bytes.IndexAny(data, ">/ ")
		data = []byte(fmt.Sprintf("<\033[32m%s\033[0m%s",
			data[1:indx], data[indx:]))
	}
	cstaLogger.Printf("%s %s %s %s",
		fmt.Sprintf("[%5s:%4s]", c.SN, c.Ext),
		LogINOUT[inFlag],
		fmt.Sprintf(fmtID, id),
		data)
}
