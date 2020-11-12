package gopher_and_rabbit

type Configuration struct {
	AMQPConnectionURL string
}

type MySqlConfiguration struct {
	MysqlURL string
}

type AddTask struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Operation string `json:"operation"`
	RecordId  int64  `json:"recordId"`
	Response  string `json:"response"`
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
}

var MysqlConfig = MySqlConfiguration{
	MysqlURL: "go:go@tcp(localhost:3306)/mauqah",
}
