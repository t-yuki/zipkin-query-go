package mysql_test

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/t-yuki/zipkin-go/models"
	"github.com/t-yuki/zipkin-go/storage/mysql"
)

// to run:  MYSQL_TCP_PORT=3306 MYSQL_HOST=localhost MYSQL_USER=zipkin MYSQL_PASS=zipkin go test -v

func TestStoreSpans(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	store, err := mysql.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	traceID := fmt.Sprintf("%016x", rand.Int63())
	id, name, duration := fmt.Sprintf("%016x", rand.Int63()), escape("TestPostSpans"), int64(time.Microsecond*1000)
	id2, name2, duration2 := fmt.Sprintf("%016x", rand.Int63()), escape("TestPostSpans2"), int64(time.Microsecond*1000)
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
	err = store.StoreSpans(models.ListOfSpans(req))
	if err != nil {
		t.Fatal(err)
	}
}

func Int64(n int64) *int64 {
	return &n
}

func shouldEscape(c byte) bool {
	if ('A' <= c && 'Z' >= c) || c == '%' {
		return true
	}
	return false
}

func escape(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		if shouldEscape(s[i]) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789abcdef"[c>>4]
			t[j+2] = "0123456789abcdef"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
