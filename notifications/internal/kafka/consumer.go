package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"route256/notifications/internal/model"
	"time"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
}

func (c *Consumer) Subscribe(topic string, handle func(text string) error, ch chan model.Info) error {

	// получаем все партиции топика
	partitionList, err := c.SingleConsumer.Partitions(topic)
	if err != nil {
		return err
	}

	/*
	   sarama.OffsetOldest - перечитываем каждый раз все
	   sarama.OffsetNewest - перечитываем только новые

	   Можем задавать отдельно на каждую партицию
	   Также можем сходить в отдельное хранилище и взять оттуда сохраненный offset
	*/
	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, err := c.SingleConsumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				txtMSG, err := bindToText(message)
				user, order, timeSt, err := bindToResp(message)
				if err != nil {
					log.Fatalln("binding: ", err)
				}

				err = handle(txtMSG)
				if err != nil {
					log.Fatalln("handler ", err)
				}

				fmt.Println("Read Topic: ", topic, " Partition: ", partition, " Offset: ", message.Offset)
				fmt.Println("Received Key: ", string(message.Key), " Value: ", string(message.Value))
				ch <- model.Info{
					UserID:    user,
					OrderID:   order,
					TimeStamp: timeSt,
				}
			}
		}(pc, partition)
	}

	return nil
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)

	if err != nil {
		return nil, err
	}

	/*
		consumer.Topics() - список топиков
		consumer.Partitions("test_topic") - партиции топика
		consumer.ConsumePartition("test_topic", 1, 12) - чтение конкретного топика с 12 сдвига в первой партиции
		consumer.Pause() - останавливаем чтение определенных топиков
		consumer.Resume() - восстанавливаем чтение определенных топиков
		consumer.PauseAll() - останавливаем чтение всех топиков
		consumer.ResumeAll() - восстанавливаем чтение всех топиков
	*/

	return &Consumer{
		brokers:        brokers,
		SingleConsumer: consumer,
	}, err
}

func bindToText(msg *sarama.ConsumerMessage) (string, error) {
	time := msg.Timestamp.Format("2006-01-02 15:04:05")

	type info struct {
		UserID  int64  `json:"UserID"`
		OrderID int64  `json:"OrderID"`
		Status  string `json:"Status"`
	}
	var msgInfo info

	err := json.Unmarshal(msg.Value, &msgInfo)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	return fmt.Sprintf(
		"%s: User's %d Order ID %d new status is %s",
		time, msgInfo.UserID, msgInfo.OrderID, msgInfo.Status), nil
}

func bindToResp(msg *sarama.ConsumerMessage) (int64, int64, time.Time, error) {
	timeStamp := msg.Timestamp

	type info struct {
		UserID  int64  `json:"UserID"`
		OrderID int64  `json:"OrderID"`
		Status  string `json:"Status"`
	}
	var msgInfo info

	err := json.Unmarshal(msg.Value, &msgInfo)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("unmarshal: %w", err)
	}

	return msgInfo.UserID, msgInfo.OrderID, timeStamp, nil
}
