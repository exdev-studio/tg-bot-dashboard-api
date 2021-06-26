package apiserver

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	BindAddr string
	LogLevel logrus.Level
}
