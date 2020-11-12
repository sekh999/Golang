module consumer

replace rabbitmq => ../

go 1.15

require (
	rabbitmq v0.0.0-00010101000000-000000000000
	github.com/streadway/amqp v1.0.0
)
