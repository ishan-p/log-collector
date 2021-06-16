package logcollector

type ServerConnection struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	ServerWaitTimeSec int    `json:"wait_time"`
}

func defaultServerConnection() ServerConnection {
	var connection ServerConnection
	connection.Host = "127.0.0.1"
	connection.Port = 8888
	connection.ServerWaitTimeSec = 3
	return connection
}
