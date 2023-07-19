package kafka

import (
	"encoding/json"
	"fmt"
	"route256/loms/internal/model"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers      []string
	topic        string
	syncProducer sarama.SyncProducer
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	// случайная партиция
	// syncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	// по кругу
	// syncProducerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	// по ключу
	// syncProducerConfig.Producer.Partitioner = sarama.NewHashPartitioner
	/**
	Кейсы:
		- одинаковые ключи в одной партиции
		- при cleanup.policy = compact останется только последнее сообщение по этому ключу
	*/
	syncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	/*
	  Если хотим exactly once, то выставляем в true

	  У продюсера есть счетчик (count)
	  Каждое успешно отправленное сообщение учеличивает счетчик (count++)
	  Если продюсер не смог отправить сообщение, то счетчик не меняется и отправляется в таком виде в другом сообщение
	  Кафка это видит и начинает сравнивать (в том числе Key) сообщения с одниковыми счетчиками
	  Далее не дает отправить дубль, если Idempotent = true
	*/
	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	// Если хотим сжимать, то задаем нужный уровень кодировщику
	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault

	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	// И сам кодировщик
	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	return syncProducer, nil
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	producer := &Producer{
		brokers:      brokers,
		topic:        topic,
		syncProducer: syncProducer,
	}

	return producer, nil
}

func (k *Producer) SendMessage(message model.ProducerMessage) error {
	msgSarama, err := k.bindToSaramaMessage(message)
	if err != nil {
		return fmt.Errorf("bind to sarama message style: %w", err)
	}

	partition, offset, err := k.syncProducer.SendMessage(msgSarama)
	if err != nil {
		return fmt.Errorf("send message to: partition: %d offset: %d,  %w", partition, offset, err)
	}

	return nil
}

func (k *Producer) bindToSaramaMessage(in model.ProducerMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(in)

	if err != nil {
		fmt.Println("Send message marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     k.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(in.UserID, in.OrderID)),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	return nil
}
