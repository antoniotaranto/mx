package mx

import "encoding/xml"

// MonitorStart запускает монитор пользователя и возвращает его идентификатор.
// В качестве параметра указывается внутренний номер пользователя. Если номер
// не указан, то используется внутренний номер авторизованного пользователя.
// В последнем случае идентификатор монитора сохраняется в Conn.MonitorID.
func (c *Conn) MonitorStart(ext string) (int64, error) {
	if ext == "" {
		ext = c.Ext // номер авторизованного пользователя
	}
	// если монитор пользователя уже запущен, то возвращаем его идентификатор
	if ext == c.Ext && c.MonitorID != 0 {
		return c.MonitorID, nil
	}
	// отправляем команду на запуск монитора
	resp, err := c.SendWithResponse(&struct {
		XMLName xml.Name `xml:"MonitorStart"`
		Ext     string   `xml:"monitorObject>deviceObject"`
	}{
		Ext: ext,
	})
	if err != nil {
		return 0, err
	}
	// разбираем идентификатор монитора
	var monitor = new(struct {
		ID int64 `xml:"monitorCrossRefID"`
	})
	if err = resp.Decode(monitor); err != nil {
		return 0, err
	}
	// сохраняем идентификатор монитора пользователя, если это он
	if ext == c.Ext {
		c.MonitorID = monitor.ID
	}
	return monitor.ID, nil
}

// MonitorStop останавливает ранее запущенный монитор пользователя. Если
// идентификатор монитора не задан, то останавливается монитор авторизованного
// пользователя, если он был раньше запущен.
func (c *Conn) MonitorStop(id int64) error {
	if id == 0 {
		if c.MonitorID != 0 {
			id = c.MonitorID // идентификатор монитора авторизованного пользователя
		} else {
			return nil // монитор не был запущен
		}
	}
	// отправляем команду на остановку монитора
	_, err := c.SendWithResponse(&struct {
		XMLName xml.Name `xml:"MonitorStop"`
		ID      int64    `xml:"monitorCrossRefID"`
	}{
		ID: id,
	})
	return err
}
