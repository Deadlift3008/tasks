package main

import (
	"fmt"
	"sync"
	"time"
)

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

type Message []byte

type IBroker interface {
	WriteToTopic(topic, message Message)
	Subscribe(topic string, ch chan Message)
}

type SyncableList[T any] struct {
	list []T
	mu   *sync.Mutex
}

type DeliveryData struct {
	consumers SyncableList[chan Message]
	messages  SyncableList[Message]
}

type DeliveryMap map[string]DeliveryData

type Broker struct {
	deliveryMap DeliveryMap
}

func (b *Broker) WriteToTopic(topic string, message Message) {
	deliveryData, ok := b.deliveryMap[topic]

	if !ok {
		b.createTopic(topic)
	}

	deliveryData.messages.mu.Lock()
	defer deliveryData.messages.mu.Unlock()

	deliveryData.messages.list = append(deliveryData.messages.list, message)
}

func (b *Broker) Subscribe(topic string, ch chan Message) {
	deliveryData, ok := b.deliveryMap[topic]

	if !ok {
		b.createTopic(topic)
	}

	currentConsumers := &deliveryData.consumers

	currentConsumers.mu.Lock()
	defer currentConsumers.mu.Unlock()

	currentConsumers.list = append(currentConsumers.list, ch)
}

func (b *Broker) createTopic(topic string) {
	b.deliveryMap[topic] = DeliveryData{}
}

func (b *Broker) sendMessages(topic string) error {
	deliveryData, ok := b.deliveryMap[topic]

	if !ok {
		return fmt.Errorf("No topic for sending messages: %s", topic)
	}

	var consumersListCopy []chan Message
	var messagesListCopy []Message

	copy(consumersListCopy, deliveryData.consumers.list)
	copy(messagesListCopy, deliveryData.messages.list)

	wg := sync.WaitGroup{}

	for i := 0; i < len(consumersListCopy); i++ {
		wg.Add(1)

		go func(index int) {
			for j := 0; j < len(messagesListCopy); j++ {
				consumersListCopy[index] <- messagesListCopy[j]
			}

			newDeliveryData, ok := b.deliveryMap[topic]

			if !ok {
				return
			}

			newDeliveryData.messages.mu.Lock()
			defer newDeliveryData.messages.mu.Unlock()

			newDeliveryData.messages.list = newDeliveryData.messages.list[len(messagesListCopy):]

			defer wg.Done()
		}(i)
	}

	wg.Wait()

	return nil
}

func (b *Broker) StartObserving(intervalMs int, workersCount int) {
	taskChan := make(chan string)

	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(intervalMs))

			for topic := range b.deliveryMap {
				taskChan <- topic
			}
		}
	}()

	go func() {
		for i := 0; i < workersCount; i++ {
			go func() {
				for {
					topic := <-taskChan
					b.sendMessages(topic)
				}
			}()
		}
	}()
}

func NewBroker() *Broker {
	return &Broker{}
}
