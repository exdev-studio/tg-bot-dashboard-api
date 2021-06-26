package main

import (
	"flag"

	"github.com/exdev-studio/tg-bot-dashboard-api/internal/app/apiserver"
	"github.com/sirupsen/logrus"
)

var (
	bindAddr string
)

func init() {
	flag.StringVar(&bindAddr, "bind-addr", "0.0.0.0:8080", "will be used as an address for listening to requests; format 0.0.0.0:8080")
}

func main() {
	flag.Parse()

	c := &apiserver.Config{
		BindAddr: bindAddr,
		LogLevel: logrus.DebugLevel,
	}
	l := buildLogger(c.LogLevel)

	if err := apiserver.Start(c, l); err != nil {
		l.Fatal(err)
	}
}

func buildLogger(logLevel logrus.Level) *logrus.Logger {
	l := logrus.New()

	l.SetLevel(logLevel)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return l
}
