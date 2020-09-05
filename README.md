## Tencent CLS Library

### Example

- If you only want to use the upload function
```go
package main
import (
    "time"

	"github.com/abcdsxg/clslog"
    "github.com/abcdsxg/clslog/pb"
)

func main() {
    var cfg = clslog.ClsClientConfig{
        TopicID:         "TopicID",
        SecretID:        "SecretID",
        SecretKey:       "SecretKey",
        ClsEndPoint:     "ClsEndPoint",   //https://cloud.tencent.com/document/product/614/18940
        MaxBatchUploads: 2,               //Greater than the total number of logs, uploading
        DelayUploadTime: 5 * time.Second, //Greater than this time interval, uploading
    }

    client := clslog.NewClsClient(cfg)
    
    for i := 0; i < 20; i++ {
        m := map[string]interface{}{
            "level": "debug",
            "name":  "hello world",
        }

        //uploading by map
        client.UploadMap(m)

        //or use UploadStructuredLog to custom log
        log:=&pb.Log{}
        client.UploadStructuredLog(log)
    }
}
```

- If you want to use hooks to upload automatically
```go
package main
import (
    "time"

	"github.com/abcdsxg/clslog"
    "github.com/abcdsxg/clslog/pb"
    "github.com/sirupsen/logrus"
    "github.com/rs/zerolog/log"
    "github.com/rs/zerolog"
)

func main() {
    var cfg = clslog.ClsClientConfig{
        TopicID:         "TopicID",
        SecretID:        "SecretID",
        SecretKey:       "SecretKey",
        ClsEndPoint:     "ClsEndPoint",   //https://cloud.tencent.com/document/product/614/18940
        MaxBatchUploads: 2,               //Greater than the total number of logs, uploading
        DelayUploadTime: 5 * time.Second, //Greater than this time interval, uploading
    }
	client := clslog.NewClsClient(cfg)

	logrusHook := clslog.NewLogrusHook(client, logrus.AllLevels)
	logrus.AddHook(logrusHook) //add logrus hook

	zerologHook := clslog.NewZerologHook(client, []zerolog.Level{
		zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.PanicLevel,
	})

	log.Logger = log.Hook(zerologHook)//add zerolog hook
}

```

## License

[MIT](http://opensource.org/licenses/MIT)

Copyright (c) 2020-present, Shingle