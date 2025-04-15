package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/anhhuy1010/DATN-cms-customer/helpers/util"
	"github.com/gin-gonic/gin"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

type JSONLog struct {
	Service  string      `json:"service"`
	Url      string      `json:"url"`
	Method   string      `json:"method"`
	Headers  interface{} `json:"headers"`
	Request  interface{} `json:"Request"`
	Response interface{} `json:"Response"`
	Level    string      `json:"level"`
	Message  string      `json:"message"`
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogRequest(c *gin.Context, request []byte, response []byte, message string) {
	go func() {
		var requestJSON map[string]interface{}
		_ = json.Unmarshal(request, &requestJSON)

		var responseJSON map[string]interface{}
		_ = json.Unmarshal(response, &responseJSON)

		var inInterface map[string]interface{}
		inrec, err := json.Marshal(JSONLog{
			Service:  "system-service",
			Url:      c.Request.URL.String(),
			Method:   c.Request.Method,
			Headers:  c.Request.Header,
			Request:  requestJSON,
			Response: responseJSON,
			Level:    "information",
			Message:  message,
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		inrec = append(inrec, '\n')
		err = json.Unmarshal(inrec, &inInterface)
		if err != nil {
			fmt.Println(err)
			return
		}
		util.LogPrint(inInterface)
	}()
}

func readBody(reader io.Reader) []byte {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)

	s := buf.Bytes()
	return s
}
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = rdr2

		c.Next()

		LogRequest(c, readBody(rdr1), blw.Body.Bytes(), "Request - Response")
	}
}
