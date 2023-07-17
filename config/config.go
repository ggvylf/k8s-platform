package config

var (
	// api
	ListenAddr = "0.0.0.0:10000"

	// websocket
	WsAddr = "0.0.0.0:10001"

	// k8s
	Kubeconfig      = "~/.kube/config"
	PodLogTailLline = 100

	// mysql
	DbUser = "root"
	DbPwd  = "123456"
	DbHost = "127.0.0.1"
	DbPort = 3306
	DbName = "k8splatform"

	DbMaxIdles    = 3
	DbMaxConns    = 50
	DbMaxLifetime = 30

	// jwt
	JwtSecret = "my-secret"

	// webshell

)
