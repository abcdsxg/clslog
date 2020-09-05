package clslog

import (
	"fmt"
	"time"

	"github.com/abcdsxg/clslog/pb"
	"github.com/sirupsen/logrus"
)

type LogrusHook struct {
	cls       *ClsClient
	LogLevels []logrus.Level
}

func NewLogrusHook(cls *ClsClient, logLevels []logrus.Level) *LogrusHook {
	return &LogrusHook{
		cls:       cls,
		LogLevels: logLevels,
	}
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *LogrusHook) Fire(entry *logrus.Entry) error {
	log := entryToPbLog(entry)
	hook.cls.UploadStructuredLog(log)
	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *LogrusHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func entryToPbLog(entry *logrus.Entry) *pb.Log {
	contents := []*pb.Log_Content{}

	contents = append(contents, &pb.Log_Content{
		Key:   "message",
		Value: entry.Message,
	})

	contents = append(contents, &pb.Log_Content{
		Key:   "level",
		Value: entry.Level.String(),
	})

	for k, v := range entry.Data {
		contents = append(
			contents,
			&pb.Log_Content{
				Key:   k,
				Value: fmt.Sprintf("%v", v),
			},
		)
	}

	return &pb.Log{
		Time:     entry.Time.UnixNano() / int64(time.Millisecond),
		Contents: contents,
	}
}
