module publisher

go 1.15

replace rabbitmq => ../

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/streadway/amqp v1.0.0
	rabbitmq v0.0.0-00010101000000-000000000000
)
