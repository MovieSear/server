package configs

const (
	HOST = "localhost"
	PORT = "9999"

	// 此处由于server服务是运行在docker容器内的，
	// 而mongodb服务则运行在宿主主机上，
	// 因此在容器内可以通过172.17.0.1来访问宿主主机
	MONGODB_URL  = "172.17.0.1:27017"
	MONGODB_NAME = "ABTicket"
)
