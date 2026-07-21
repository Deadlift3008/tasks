/*
Нужно реализовать брокер сообщений, где сообщения — это []byte.
Брокер поддерживает несколько топиков.
В топик можно писать независимо от читателей, и писать могут несколько потоков.
Сообщения должны доставляться только один раз (каждому читателю).
Из топика могут читать несколько потоков параллельно.
Для синхронизации использовать мьютексы и каналы.
Для хранения сообщений — очередь (slice).
Для каждого топика — отдельная структура с очередью сообщений и списком читателей.

*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// Message — тип сообщения
type Message []byte

// IBroker — интерфейс брокера (исправлен с учётом реальных сигнатур)
type IBroker interface {
	WriteToTopic(topic string, message Message)
	Subscribe(topic string, ch chan Message)
}

// Topic хранит состояние одного топика
type Topic struct {
	mu         sync.Mutex
	messages   []Message
	consumers  []chan Message
	processing bool // признак того, что топик уже обрабатывается
}

// Broker реализует IBroker
type Broker struct {
	topics map[string]*Topic
	mu     sync.RWMutex // защищает map топиков
}

// NewBroker создаёт новый брокер с инициализированной map
func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

// getOrCreateTopic возвращает топик, создавая его при необходимости.
// Держит блокировку записи b.mu только на время работы с map.
func (b *Broker) getOrCreateTopic(topic string) *Topic {
	b.mu.RLock()
	t, ok := b.topics[topic]
	b.mu.RUnlock()
	if ok {
		return t
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	// повторная проверка, т.к. за время разблокировки могли создать
	if t, ok := b.topics[topic]; ok {
		return t
	}
	t = &Topic{}
	b.topics[topic] = t
	return t
}

// WriteToTopic добавляет сообщение в очередь топика
func (b *Broker) WriteToTopic(topic string, message Message) {
	t := b.getOrCreateTopic(topic)

	t.mu.Lock()
	t.messages = append(t.messages, message)
	t.mu.Unlock()
}

// Subscribe добавляет подписчика (канал) в топик
func (b *Broker) Subscribe(topic string, ch chan Message) {
	t := b.getOrCreateTopic(topic)

	t.mu.Lock()
	t.consumers = append(t.consumers, ch)
	t.mu.Unlock()
}

// sendMessages пытается разослать накопившиеся сообщения подписчикам.
// Если топик уже обрабатывается, выходит сразу (предотвращает параллельную обработку).
func (b *Broker) sendMessages(topicName string) {
	// получаем топик под read-блокировкой
	b.mu.RLock()
	t, ok := b.topics[topicName]
	b.mu.RUnlock()
	if !ok {
		return
	}

	t.mu.Lock()
	if t.processing {
		t.mu.Unlock()
		return
	}
	t.processing = true

	// делаем снимки под мьютексом топика
	// правильное копирование слайсов: создаём slice нужной длины и копируем
	messagesSnapshot := make([]Message, len(t.messages))
	copy(messagesSnapshot, t.messages)

	consumersSnapshot := make([]chan Message, len(t.consumers))
	copy(consumersSnapshot, t.consumers)

	t.mu.Unlock()

	// рассылаем сообщения всем подписчикам
	// для обработки медленных/отвалившихся потребителей используем неблокирующую отправку с таймаутом
	const sendTimeout = 100 * time.Millisecond

	var wg sync.WaitGroup
	for _, ch := range consumersSnapshot {
		wg.Add(1)
		go func(ch chan Message) {
			defer wg.Done()
			for _, msg := range messagesSnapshot {
				select {
				case ch <- msg:
					// успешно отправлено
				case <-time.After(sendTimeout):
					// потребитель не успел принять — считаем его мёртвым и удаляем
					b.removeConsumer(t, ch)
					return
				}
			}
		}(ch)
	}
	wg.Wait()

	// после рассылки удаляем отправленные сообщения
	t.mu.Lock()
	// если очередь не была модифицирована параллельно, удаляем snapshot
	if len(t.messages) >= len(messagesSnapshot) {
		t.messages = t.messages[len(messagesSnapshot):]
	} else {
		// на всякий случай очищаем всё (не должно происходить при корректном использовании)
		t.messages = nil
	}
	t.processing = false
	t.mu.Unlock()
}

// removeConsumer удаляет канал подписчика из топика (вызывается, если отправка не удалась)
func (b *Broker) removeConsumer(t *Topic, ch chan Message) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, c := range t.consumers {
		if c == ch {
			// удаление без сохранения порядка
			t.consumers[i] = t.consumers[len(t.consumers)-1]
			t.consumers = t.consumers[:len(t.consumers)-1]
			break
		}
	}
}

// StartObserving запускает фоновую рассылку: раз в intervalMs миллисекунд
// все топики отправляются в пул воркеров (workersCount штук).
func (b *Broker) StartObserving(intervalMs int, workersCount int) {
	taskChan := make(chan string, 100) // буферизированный канал задач

	// горутина-генератор задач: периодически кидает все известные топики в канал
	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(intervalMs))
		defer ticker.Stop()
		for range ticker.C {
			b.mu.RLock()
			for topicName := range b.topics {
				select {
				case taskChan <- topicName:
				default:
					// если канал переполнен — пропускаем, чтобы не блокироваться
				}
			}
			b.mu.RUnlock()
		}
	}()

	// пул воркеров
	for i := 0; i < workersCount; i++ {
		go func() {
			for topicName := range taskChan {
				b.sendMessages(topicName)
			}
		}()
	}
}

// для примера использования
func main() {
	broker := NewBroker()

	ch1 := make(chan Message, 10)
	ch2 := make(chan Message, 10)

	broker.Subscribe("orders", ch1)
	broker.Subscribe("orders", ch2)

	broker.WriteToTopic("orders", []byte("order1"))
	broker.WriteToTopic("orders", []byte("order2"))

	broker.StartObserving(500, 2)

	time.Sleep(1 * time.Second)

	// читаем
	fmt.Println("Consumer 1:", string(<-ch1))
	fmt.Println("Consumer 1:", string(<-ch1))
	fmt.Println("Consumer 2:", string(<-ch2))
	fmt.Println("Consumer 2:", string(<-ch2))
}
