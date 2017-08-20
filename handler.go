package mx

import (
	"time"
)

// Handler описывает функцию для обработки событий. Если функция возвращает
// ошибку, то дальнейшая обработка событий прекращается.
type Handler = func(*Response) error

// Stop для остановки обработки событий в Handle.
var Stop error = new(emptyError)

// HandleWait вызывает переданную функцию handler для обработки все событий с
// названиями из списка events. timeout задает максимальное время ожидания
// ответа от сервера. По истечение времени ожидания возвращается ошибка
// ErrTimeout. Если timeout установлен в 0 или отрицательный, то время ожидания
// ответа не ограничено. Для планового завершения обработки можно в качестве
// ошибки вернуть mx.Stop: выполнение прервется, но в ответе ошибкой будет nil.
func (c *Conn) HandleWait(handler Handler, timeout time.Duration,
	events ...string) (err error) {
	if len(events) == 0 {
		return nil // нет событий для отслеживания
	}
	// создаем канал для получения ответов от сервера и регистрируем его для
	// всех заданных имен событий
	var eventChan = make(chan *Response, 1)
	for _, event := range events {
		list, ok := c.eventHandlers.Load(event)
		if !ok {
			list = mapOfHandlerChan{eventChan: struct{}{}}
		} else {
			list.(mapOfHandlerChan)[eventChan] = struct{}{}
		}
		c.eventHandlers.Store(event, list) // сохраняем обновленный список
	}

	// взводим таймер ожидания ответа
	var timeoutTimer = time.NewTimer(timeout)
	if timeout <= 0 {
		<-timeoutTimer.C // сбрасываем таймер
	}
processing: // ждем ответа или
	select {
	case resp := <-eventChan: // получили событие от сервера
		// пустой ответ приходит только в случае закрытия соединения
		if resp == nil {
			if !timeoutTimer.Stop() {
				<-timeoutTimer.C
			}
			err = c.err // возвращаем ошибку соединения
			break
		}
		// запускаем обработчик события
		err = handler(resp)
		switch err {
		case nil:
			if timeout > 0 { // сдвигаем таймер, если задано время ожидания
				if !timeoutTimer.Stop() {
					<-timeoutTimer.C
				}
				timeoutTimer.Reset(timeout)
			}
			goto processing // ждем следующего ответа для обработки
		case Stop:
			err = nil // сбрасываем ошибку
		}
	case <-timeoutTimer.C:
		err = ErrTimeout // ошибка времени ожидания
	}
	// удаляем канал обработки для всех событий
	for _, event := range events {
		// получаем обработчики событий
		list, ok := c.eventHandlers.Load(event)
		if !ok {
			continue // обработчики не зарегистрированы
		}
		// проверяем, что наш обработчик есть в списке
		var handlers = list.(mapOfHandlerChan)
		if _, ok := handlers[eventChan]; ok {
			if len(handlers) < 2 {
				// кроме нашего обработчика других нет
				c.eventHandlers.Delete(event)
			} else {
				// в списке есть другие обработчики
				delete(handlers, eventChan)
				c.eventHandlers.Store(event, handlers)
			}
		}
	}
	close(eventChan) // закрываем наш канал
	return err
}

// Handle просто вызывает HandleWait с установленным временем ожидания 0.
func (c *Conn) Handle(handler Handler, events ...string) error {
	return c.HandleWait(handler, 0, events...)
}

// mapOfHandlerChan используется в качестве синонима для описания списка
// каналов для обработки событий.
type mapOfHandlerChan = map[chan<- *Response]struct{}
