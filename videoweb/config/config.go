package config

//[mysql]
const (
	SqlUserName = "root"
	SqlPassword = "123456"
	DataBase    = "videoWebsite"
	MysqlIP     = "192.168.239.126" //"192.168.1.107.3306" //ipconfig->WLAN->ipv4
)

//[RabbitMQ]
const (
	RabbitmqUserName = "admin"
	RabbitmqPassword = "123456"
	RabbitmqIP       = "192.168.136.128:5672"
	ExchangeName     = "direct"
)

//[redis]
const (
	RedisAddr     = "172.23.21.149:6379" //"redis:6379" //"192.168.1.100:6379"
	RedisPassword = "root"
	RedisDB       = 0
)

//[elasticsearch]
const (
	EsAddr                 = "https://localhost:9200"
	EsName                 = "elastic"
	EsPassword             = "123456"
	CertificateFingerprint = "81bdbeefa2b378da4fd29dfd4d8e82e96cd0d92a50b24db5ec30401260fad917"
)

//[MinIo]
const (
	EndPoint        = "172.23.21.149:9000"
	AccessKeyID     = "minioadmin"
	SecretAccessKey = "minioadmin"
	SSL             = false
	BucketName      = "videoweb"
)
