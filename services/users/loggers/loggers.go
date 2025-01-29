package loggers

import "github.com/sirupsen/logrus"

var ServerLogger *logrus.Logger
var EventLogger *logrus.Logger

func init() {
  EventLogger = logrus.New()
  EventLogger.SetFormatter(&logrus.JSONFormatter{
    DisableTimestamp: false,
  })
  ServerLogger = logrus.New()
  ServerLogger.SetFormatter(&logrus.TextFormatter{
    ForceColors: true,
    DisableColors: false,
    EnvironmentOverrideColors: true,
    FullTimestamp: true,
  TimestampFormat: "2006/1/2 15:04:05",
  })
}
