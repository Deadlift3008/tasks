package main

import (
	"fmt"
	"sync"
	"time"
)

// Message – тип сообщения.
type Message []byte

// IBroker – интерфейс брокера.
type IBroker interface {
	WriteToTopic(topic string, msg Message)
	Subscribe(topic string, ch chan Message)
}

// ======================== Команды для диспетчера топика ======================

type cmdType int

const (
	cmdSubscribe   cmdType = iota // подписаться: ch – канал подписчика
	cmdUnsubscribe                // отписаться: ch – канал подписчика
	cmdWrite                      // записать сообщение: msg – данные
	cmdStop                       // остановить диспетчер
)

type command struct {
	typ cmdType
	ch  chan Message
	msg Message
}

// ======================== Диспетчер одного топика ============================

// Topic управляет одним топиком.
// Вся работа с данными топика ведётся в горутине, запущенной методом start().
type Topic struct {
	commandCh chan command // канал для приёма команд

	// Данные, изменяемые только из горутины-диспетчера
	subscribers map[chan Message]int // ключ – канал подписчика, значение – следующий неотправленный индекс
	messages    []Message            // общая очередь сообщений
}

// newTopic создаёт топик и запускает его горутину-диспетчер.
func newTopic() *Topic {
	t := &Topic{
		commandCh:   make(chan command, 100), // буферизованный канал, чтобы писатели не блокировались
		subscribers: make(map[chan Message]int),
	}
	go t.run()
	return t
}

// stop мягко останавливает горутину диспетчера.
func (t *Topic) stop() {
	t.commandCh <- command{typ: cmdStop}
}

// run – основной цикл диспетчера топика.
func (t *Topic) run() {
	// Тикер для периодической попытки доставки накопившихся сообщений медленным подписчикам
	retryTicker := time.NewTicker(100 * time.Millisecond)
	defer retryTicker.Stop()

	for {
		select {
		case cmd := <-t.commandCh:
			switch cmd.typ {
			case cmdSubscribe:
				// Новый подписчик начинает получать сообщения, которые будут добавлены после него.
				// Текущая длина очереди – его стартовый офсет.
				t.subscribers[cmd.ch] = len(t.messages)

			case cmdUnsubscribe:
				delete(t.subscribers, cmd.ch)

			case cmdWrite:
				t.messages = append(t.messages, cmd.msg)
				// Сразу пытаемся разослать сообщение всем подписчикам
				t.deliverPending()

			case cmdStop:
				// Мягкая остановка: пытаемся доставить всё, что ещё не доставлено
				t.deliverPending()
				return
			}

		case <-retryTicker.C:
			// Периодически подчищаем хвост очереди и пробуем доставить отставшим
			t.deliverPending()
		}
	}
}

// deliverPending пробует отправить каждому подписчику все сообщения,
// начиная с его сохранённого офсета. Доставка неблокирующая (select-default).
// После успешной отправки офсет сдвигается.
func (t *Topic) deliverPending() {
	// Определяем минимальный офсет среди всех подписчиков
	minOffset := len(t.messages)
	for _, offset := range t.subscribers {
		if offset < minOffset {
			minOffset = offset
		}
	}

	// Удаляем сообщения, которые уже получили все подписчики
	if minOffset > 0 {
		t.messages = t.messages[minOffset:]
		// Корректируем офсеты всех подписчиков
		for ch, off := range t.subscribers {
			t.subscribers[ch] = off - minOffset
		}
	}

	// Пытаемся разослать каждому подписчику его непрочитанные сообщения
	for ch, off := range t.subscribers {
		for off < len(t.messages) {
			msg := t.messages[off]
			// Неблокирующая отправка
			select {
			case ch <- msg:
				off++
				t.subscribers[ch] = off
			default:
				// Канал подписчика полон или закрыт – прерываемся, попробуем позже
				goto nextSubscriber
			}
		}
	nextSubscriber:
	}

	// Удаляем подписчиков, каналы которых закрыты (запись в закрытый канал вызывает панику,
	// но select-default защищает от этого. Однако если канал закрыт, запись в него
	// в select-default приводит к мгновенному выполнению default без паники.
	// Поэтому дополнительная проверка не требуется – если канал закрыт,
	// попытка отправить всегда будет попадать в default, и офсет не сдвинется.
	// Но если канал закрыт, подписчик «мёртв», так что можно удалить его
	// после нескольких неудачных попыток или сразу. Здесь для простоты не удаляем,
	// чтобы не усложнять код; можно добавить счётчики ошибок.
}

// ======================== Брокер ==============================================

// Broker реализует IBroker.
type Broker struct {
	topics map[string]*Topic
	mu     sync.RWMutex // защищает карту топиков
	closed bool
}

// NewBroker создаёт новый брокер.
func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

// getOrCreateTopic возвращает существующий топик или создаёт новый (потокобезопасно).
func (b *Broker) getOrCreateTopic(topic string) (*Topic, error) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return nil, fmt.Errorf("broker is closed")
	}
	t, ok := b.topics[topic]
	b.mu.RUnlock()
	if ok {
		return t, nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil, fmt.Errorf("broker is closed")
	}
	// Повторная проверка, т.к. за время разблокировки могли создать
	if t, ok = b.topics[topic]; ok {
		return t, nil
	}
	t = newTopic()
	b.topics[topic] = t
	return t, nil
}

// WriteToTopic отправляет сообщение в топик. Потокобезопасен.
func (b *Broker) WriteToTopic(topic string, msg Message) {
	t, err := b.getOrCreateTopic(topic)
	if err != nil {
		return // брокер закрыт, сообщение отбрасывается
	}
	t.commandCh <- command{typ: cmdWrite, msg: msg}
}

// Subscribe добавляет подписчика. Потокобезопасен.
func (b *Broker) Subscribe(topic string, ch chan Message) {
	t, err := b.getOrCreateTopic(topic)
	if err != nil {
		return
	}
	t.commandCh <- command{typ: cmdSubscribe, ch: ch}
}

// Unsubscribe удаляет подписчика.
func (b *Broker) Unsubscribe(topic string, ch chan Message) {
	b.mu.RLock()
	t, ok := b.topics[topic]
	b.mu.RUnlock()
	if !ok {
		return
	}
	t.commandCh <- command{typ: cmdUnsubscribe, ch: ch}
}

// Close останавливает брокер: все горутины-диспетчеры завершаются.
// Новые операции после закрытия игнорируются.
func (b *Broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, t := range b.topics {
		t.stop()
	}
}

// ======================== Демонстрация ========================================

func main() {
	broker := NewBroker()

	ch1 := make(chan Message, 10)
	ch2 := make(chan Message, 10)

	broker.Subscribe("orders", ch1)
	broker.Subscribe("orders", ch2)

	broker.WriteToTopic("orders", []byte("order1"))
	broker.WriteToTopic("orders", []byte("order2"))
	broker.WriteToTopic("orders", []byte("order3"))

	// Дадим время на доставку
	time.Sleep(300 * time.Millisecond)

	// Читаем из каналов
	fmt.Println("Consumer 1:")
	closeChAfterTime(ch1, 50*time.Millisecond, func(m Message) { fmt.Println("  ", string(m)) })
	fmt.Println("Consumer 2:")
	closeChAfterTime(ch2, 50*time.Millisecond, func(m Message) { fmt.Println("  ", string(m)) })

	// Отпишем одного и проверим, что сообщения больше не приходят
	broker.Unsubscribe("orders", ch1)
	broker.WriteToTopic("orders", []byte("order4"))
	time.Sleep(300 * time.Millisecond)

	fmt.Println("Consumer 2 after unsub1:")
	closeChAfterTime(ch2, 50*time.Millisecond, func(m Message) { fmt.Println("  ", string(m)) })

	broker.Close()
	fmt.Println("Broker closed.")
}

// Вспомогательная функция – читает из канала в течение заданного времени и выводит сообщения.
func closeChAfterTime(ch chan Message, d time.Duration, fn func(Message)) {
	timeout := time.After(d)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fn(msg)
		case <-timeout:
			return
		}
	}
}
