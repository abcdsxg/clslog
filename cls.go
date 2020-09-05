package clslog

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/abcdsxg/clslog/pb"
	"google.golang.org/protobuf/proto"
)

type ClsClientConfig struct {
	TopicID         string
	SecretID        string
	SecretKey       string
	ClsEndPoint     string
	MaxBatchUploads int
	DelayUploadTime time.Duration
}

type ClsClient struct {
	cfg      ClsClientConfig
	logsChan chan *pb.Log
}

var (
	HttpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func NewClsClient(cfg ClsClientConfig) *ClsClient {
	client := &ClsClient{
		cfg:      cfg,
		logsChan: make(chan *pb.Log, cfg.MaxBatchUploads*2),
	}

	go client.watchLogs()
	return client
}

func (c *ClsClient) UploadMap(m map[string]interface{}) {
	if m == nil || len(m) == 0 {
		return
	}

	c.UploadStructuredLog(c.mapToLogGroupList(m))
}

func (c *ClsClient) UploadStructuredLog(log *pb.Log) {
	if log == nil {
		return
	}

	c.logsChan <- log
}

func (c *ClsClient) fire(logs []*pb.Log) {
	logGroupList := &pb.LogGroupList{
		LogGroupList: []*pb.LogGroup{
			{
				Logs: logs,
			},
		},
	}
	pbData, err := proto.Marshal(logGroupList)
	if err != nil {
		log.Println("UploadStructuredLog Marshal err:", err)
		return
	}

	c.upload(pbData)
}

func (c *ClsClient) upload(pbData []byte) {
	u := fmt.Sprintf("%s/structuredlog?topic_id=%s", c.cfg.ClsEndPoint, c.cfg.TopicID)
	req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(pbData))
	req.Header.Add("Host", strings.TrimLeft(c.cfg.ClsEndPoint, "https://"))
	req.Header.Add("Content-Type", "application/x-protobuf")

	authorization := c.genAuthorization(req)
	req.Header.Add("Authorization", authorization)

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Println("request err:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, _ := ioutil.ReadAll(resp.Body)
		log.Println("response err:", string(respData))
	}
}

func (c *ClsClient) watchLogs() {
	t := time.NewTimer(c.cfg.DelayUploadTime)
	sendLogs := []*pb.Log{}
	for {
		select {
		case <-t.C:
			if len(sendLogs) > 0 {
				go c.fire(sendLogs)
				sendLogs = []*pb.Log{}
			}
			t.Reset(c.cfg.DelayUploadTime)
		case l := <-c.logsChan:
			sendLogs = append(sendLogs, l)
			if len(sendLogs) >= c.cfg.MaxBatchUploads {
				go c.fire(sendLogs)
				sendLogs = []*pb.Log{}
			}
		}
	}
}

func (c *ClsClient) mapToLogGroupList(m map[string]interface{}) *pb.Log {
	l := &pb.Log{
		Time: time.Now().UnixNano() / int64(time.Millisecond),
	}

	for k, v := range m {
		// key cannot start with _
		if strings.HasPrefix(k, "_") {
			k = strings.TrimLeft(k, "_")
		}

		content := &pb.Log_Content{
			Key:   k,
			Value: fmt.Sprintf("%#v", v),
		}
		l.Contents = append(l.Contents, content)
	}
	return l
}

func (c *ClsClient) genAuthorization(req *http.Request) string {
	algorithm := "sha1"
	now := time.Now().Unix()
	expireTime := now + 60
	signTime := fmt.Sprintf("%d;%d", now, expireTime)

	//If you don't sign the FormatedParameters, you need to splice the \n\n.
	httpRequestInfo := strings.ToLower(req.Method) + fmt.Sprintln() + req.URL.Path + fmt.Sprintln() + fmt.Sprintln()

	var headerKeys []string // Lowercase the key here.
	var keys []string
	for key := range req.Header {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	var tmpStr string
	for _, key := range keys {
		headerKeys = append(headerKeys, strings.ToLower(key))
		tmpStr += strings.ToLower(key) + "=" + url.QueryEscape(req.Header.Get(key)) + "&"
	}
	tmpStr = strings.TrimRight(tmpStr, "&")
	httpRequestInfo += tmpStr + fmt.Sprintln()

	stringToSign := algorithm + fmt.Sprintln() + signTime + fmt.Sprintln() + Sha1(httpRequestInfo) + fmt.Sprintln()
	signKey := HmacSha1(signTime, c.cfg.SecretKey)
	signature := HmacSha1(stringToSign, signKey)

	authContent := fmt.Sprintf("q-sign-algorithm=%s&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=%s&q-url-param-list=%s&q-signature=%s",
		algorithm, c.cfg.SecretID, signTime, signTime,
		strings.Join(headerKeys, ";"), "", signature)

	return authContent
}

func Sha1(content string) string {
	s := sha1.New()
	s.Write([]byte(content))
	return hex.EncodeToString(s.Sum(nil))
}

func HmacSha1(content string, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(content))
	return hex.EncodeToString(mac.Sum(nil))
}

func Md5(content []byte) string {
	m := md5.New()
	m.Write(content)
	return hex.EncodeToString(m.Sum(nil))
}
