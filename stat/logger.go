package main

import "github.com/sirupsen/logrus"

func initLogger() {
	logrus.SetLevel(logrus.DebugLevel)
}
