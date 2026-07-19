package main

import "sync"

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

// TODO: дописать sendMessages
// TODO: поправить синхронизацию для кейсов:
// когда пишут много и когда читателей много
// чтобы сохранялась независимость

type Message []byte

type IBroker interface {
	WriteToTopic(topic, message Message) error
	Subscribe(topic string, ch chan Message) error
}

type DeliveryData struct {
	consumers []chan Message
	messages  []Message
	mu        sync.Mutex
	cmu       sync.Mutex
}

type DeliveryMap map[string]DeliveryData

type Broker struct {
	deliveryMap DeliveryMap
}

func (b *Broker) WriteToTopic(topic string, message Message) error {
	deliveryData, ok := b.deliveryMap[topic]

	if !ok {
		b.createTopic(topic)
	}

	deliveryData.mu.Lock()
	defer deliveryData.mu.Unlock()

	deliveryData.messages = append(deliveryData.messages, message)

	go b.sendMessages(topic)

	return nil
}

func (b *Broker) Subscribe(topic string, ch chan Message) error {
	deliveryData, ok := b.deliveryMap[topic]

	if !ok {
		b.createTopic(topic)
	}

	deliveryData.cmu.Lock()
	defer deliveryData.cmu.Unlock()

	deliveryData.consumers = append(deliveryData.consumers, ch)
	return nil
}

func (b *Broker) createTopic(topic string) {
	b.deliveryMap[topic] = DeliveryData{}
}

func (b *Broker) sendMessages(topic string) error {
	deliveryData, ok := b.deliveryMap[topic]

}

func NewBroker() *Broker {
	return &Broker{}
}
