package clslog

import (
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"

	"github.com/sirupsen/logrus"
)

func TestLogrusHook(t *testing.T) {
	client := NewClsClient(cfg)

	hook := NewLogrusHook(client, logrus.AllLevels)
	logrus.AddHook(hook)

	for i := 0; i < 10; i++ {
		logrus.WithError(errors.New("there's no error")).WithField("hello", "world").Debug("test error")
		logrus.WithError(errors.New("there's no error")).WithField("hello", "world").Info("test error")
		logrus.WithError(errors.New("there's no error")).WithField("hello", "world").Warn("test error")
		logrus.WithError(errors.New("there's no error")).WithField("hello", "world").Error("test error")
	}

	time.Sleep(10 * time.Second)
}

func TestZerologHook(t *testing.T) {
	client := NewClsClient(cfg)

	hook := NewZerologHook(client, []zerolog.Level{
		zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.PanicLevel,
	})

	log.Logger = log.Hook(hook)

	for i := 0; i < 10; i++ {
		log.Debug().Err(errors.New("there's no error")).Str("hello", "world").Msg("message no error")
		log.Info().Err(errors.New("there's no error")).Str("hello", "world").Msg("message no error")
		log.Warn().Err(errors.New("there's no error")).Str("hello", "world").Msg("message no error")
		log.Error().Err(errors.New("there's no error")).Str("hello", "world").Msg("message no error")
	}

	time.Sleep(10 * time.Second)
}
