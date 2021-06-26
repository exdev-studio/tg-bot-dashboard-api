package apiserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func Start(config *Config, logger *logrus.Logger) error {
	logger.WithFields(logrus.Fields{
		"bind-addr": config.BindAddr,
		"log-level": config.LogLevel,
	}).Info("server starting")

	srv := newServer(logger)
	return http.ListenAndServe(config.BindAddr, srv)
}
