Pre-Requirement

go get github.com/rabbitmq/amqp091-go

https://www.docker.com/products/docker-desktop/ download docker desktop

docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

Run

go run main.go

For Check Details about rabbitmq

http://localhost:15672 if you start with default configs


