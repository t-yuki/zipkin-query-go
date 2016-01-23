package zipkin_server_go_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/t-yuki/zipkin-go/models"
)

func TestPostSpans_empty(t *testing.T) {
	req := []*models.Span{}
	reqData, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:8081/api/v1/spans", "application/json", bytes.NewReader(reqData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Logf("%s", msg)
		t.Fatal(resp.StatusCode)
	}
}

func TestPostSpans(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	traceID := fmt.Sprintf("%016x", rand.Int63())
	id, name, duration := fmt.Sprintf("%016x", rand.Int63()), "test_post_spans", int64(time.Microsecond*1000)
	id2, name2, duration2 := fmt.Sprintf("%016x", rand.Int63()), "test_post_spans2", int64(time.Microsecond*1000)
	debug := true
	ipv4, port, serviceName1, serviceName2 := "127.0.0.1", int64(80), "store_test", "store_test2"
	ep1 := &models.Endpoint{IPV4: &ipv4, Port: &port, ServiceName: &serviceName1}
	ep2 := &models.Endpoint{IPV4: &ipv4, Port: &port, ServiceName: &serviceName2}
	ts := time.Now().UnixNano() / 1000
	annKey1, annValue1 := "key1", base64.StdEncoding.EncodeToString([]byte("value1"))
	req := []*models.Span{
		{
			TraceID:   &traceID,
			ID:        &id,
			Name:      &name,
			ParentID:  nil,
			Timestamp: ts,
			Duration:  &duration,
			Debug:     &debug,
			Annotations: []*models.Annotation{
				{ep1, Int64(ts + 100), models.AnnotationValue_SERVER_RECV.Addr()},
				{ep1, Int64(ts + 200), models.AnnotationValue_CLIENT_SEND.Addr()},
				{ep1, Int64(ts + 300), models.AnnotationValue_CLIENT_RECV.Addr()},
				{ep1, Int64(ts + 400), models.AnnotationValue_SERVER_SEND.Addr()},
			},
		},
		{
			TraceID:   &traceID,
			ID:        &id2,
			Name:      &name2,
			ParentID:  &id,
			Timestamp: ts + 200,
			Duration:  &duration2,
			Debug:     &debug,
			Annotations: []*models.Annotation{
				{ep2, Int64(ts + 210), models.AnnotationValue_SERVER_RECV.Addr()},
				{ep2, Int64(ts + 220), models.AnnotationValue_CLIENT_SEND.Addr()},
				{ep2, Int64(ts + 230), models.AnnotationValue_CLIENT_RECV.Addr()},
				{ep2, Int64(ts + 240), models.AnnotationValue_SERVER_SEND.Addr()},
			},
			BinaryAnnotations: []*models.BinaryAnnotation{
				{models.AnnotationType_STRING.Addr(), ep2, &annKey1, &annValue1},
			},
		}}

	reqData, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:8081/api/v1/spans", "application/json", bytes.NewReader(reqData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Logf("%s", msg)
		t.Fatal(resp.StatusCode)
	}
}

func Int64(n int64) *int64 {
	return &n
}
