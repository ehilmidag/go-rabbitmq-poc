package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "RabbitMQ'ya bağlanılamadı")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Channel açılamadı")
	defer ch.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go setupPublisher(ch, &wg)
	go setupConsumer(ch, &wg)

	wg.Wait()

	log.Printf(" [*] Mesajlar bekleniyor")
	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func setupPublisher(ch *amqp.Channel, wg *sync.WaitGroup) {
	defer wg.Done()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Que oluşturulamadı")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		body := "Merhaba Dünya!"
		err = ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			failOnError(err, "Mesaj yayınlanamadı")
			return
		}
		log.Printf(" [x] Gönderildi: %s\n", body)
	}
}
func setupConsumer(ch *amqp.Channel, wg *sync.WaitGroup) {
	defer wg.Done()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Que oluşturulamadı")
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Tüketici kaydedilemedi")
		return
	}

	go func() {
		for d := range msgs {
			log.Printf("Bir mesaj alındı: %s", d.Body)
		}
	}()
}
