package clslog

import (
	"testing"
	"time"
)

var cfg = ClsClientConfig{
	TopicID:         "TopicID",
	SecretID:        "SecretID",
	SecretKey:       "SecretKey",
	ClsEndPoint:     "ClsEndPoint",   //https://cloud.tencent.com/document/product/614/18940
	MaxBatchUploads: 2,               //Greater than the total number of logs,uploading
	DelayUploadTime: 5 * time.Second, //Greater than this time interval. ,uploading
}

func TestClsMaxBatchUploads(t *testing.T) {
	client := NewClsClient(cfg)

	for i := 0; i < 20; i++ {
		m := map[string]interface{}{
			"level": "debug",
			"name":  "hello world",
		}
		client.UploadMap(m)
	}

}

func TestClsDelayUploadTime(t *testing.T) {
	client := NewClsClient(cfg)

	for i := 0; i < 20; i++ {
		time.Sleep(2 * time.Second)
		m := map[string]interface{}{
			"level": "debug",
			"name":  "hello world",
		}
		client.UploadMap(m)
	}

}
