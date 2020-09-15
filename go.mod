module github.com/wikensmith/gohost

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/wikensmith/toLogCenter v0.1.4
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/wikensmith/toLogCenter v0.1.4 => github.com/wikensmith/toLogCenter v0.1.7
