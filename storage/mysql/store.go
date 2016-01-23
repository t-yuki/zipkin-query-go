package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/t-yuki/zipkin-go/models"
)

// see https://github.com/openzipkin/zipkin/blob/master/zipkin-query-service/config/query-mysql.scala
// https://github.com/openzipkin/zipkin/blob/master/zipkin-anormdb/src/main/resources/mysql.sql
// from https://github.com/openzipkin/zipkin/blob/master/zipkin-anormdb/src/main/scala/com/twitter/zipkin/storage/anormdb/AnormSpanStorage.scala

type Storage struct {
	db *sql.DB
}

func Open() (*Storage, error) {
	dbname := os.Getenv("MYSQL_DB")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_TCP_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	ssl := os.Getenv("MYSQL_USE_SSL")
	// MYSQL_MAX_CONNECTIONS is ignored
	if dbname == "" {
		dbname = "zipkin"
	}
	if host == "" {
		host = "localhost"
	}
	proto := "tcp"
	if port == "" {
		proto = "unix"
	}
	addr := ""
	if host != "" {
		if port != "" {
			port = ":" + port
		}
		addr = proto + "(" + host + port + ")"
	}
	params := url.Values{}
	if ssl != "" {
		params.Set("tls", ssl)
	}
	params.Set("loc", "Local")
	params.Set("parseTime", "true")
	params.Set("interpolateParams", "true")
	db, err := sql.Open("mysql", user+":"+pass+"@"+addr+"/"+dbname+"?"+params.Encode())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (stor *Storage) Close() error {
	return stor.db.Close()
}

func (stor *Storage) StoreSpans(spans models.ListOfSpans) error {
	for _, span := range spans {
		err := stor.storeSpan(span)
		if err != nil {
			return errors.New("storSpan: " + err.Error())
		}
	}
	return nil
}

func (stor *Storage) storeSpan(span *models.Span) error {
	// TODO: s.annotations.sorted
	// TODO: ApplyTimestampAndDuration

	if span.Timestamp == 0 { // fallback if we have no timestamp
		span.Timestamp = time.Now().UnixNano() / 1000
	}

	tx, err := stor.db.Begin()
	if err != nil {
		return errors.New("begin: " + err.Error())
	}
	defer tx.Rollback()

	if err := stor.insertSpan(tx, span); err != nil {
		return errors.New("insertSpan: " + err.Error())
	}
	if err := stor.insertAnnotations(tx, span); err != nil {
		return errors.New("insertAnnotations: " + err.Error())
	}
	if err := stor.insertBinaryAnnotations(tx, span); err != nil {
		return errors.New("insertBinaryAnnotations: " + err.Error())
	}

	if err = tx.Commit(); err != nil {
		return errors.New("commit: " + err.Error())
	}
	return nil
}

func (stor *Storage) insertSpan(tx *sql.Tx, span *models.Span) error {
	stmt, err := tx.Prepare(`INSERT INTO zipkin_spans VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE trace_id = ? AND id = ?`)
	if err != nil {
		return fmt.Errorf("prepare: %+v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		hex2int64(span.TraceID),
		hex2int64(span.ID),
		*span.Name,
		hex2int64(span.ParentID),
		span.Debug,
		span.Timestamp,
		span.Duration,
		hex2int64(span.TraceID),
		hex2int64(span.ID),
	)
	return err
}

func (stor *Storage) insertAnnotations(tx *sql.Tx, span *models.Span) error {
	stmt, err := tx.Prepare(`INSERT INTO zipkin_annotations VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare: %+v", err)
	}
	defer stmt.Close()
	for _, ann := range span.Annotations {
		ipv4, port, service := ep2sql(ann.Endpoint)
		_, err := stmt.Exec(
			hex2int64(span.TraceID),
			hex2int64(span.ID),
			ann.Value,
			nil,
			-1,
			ann.Timestamp,
			ipv4, port, service,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (stor *Storage) insertBinaryAnnotations(tx *sql.Tx, span *models.Span) error {
	stmt, err := tx.Prepare(`INSERT INTO zipkin_annotations VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare: %+v", err)
	}
	defer stmt.Close()

	for _, ann := range span.BinaryAnnotations {
		// BinaryAnnotation.Value is base64, see: https://github.com/openzipkin/zipkin/blob/master/zipkin-common/src/main/scala/com/twitter/zipkin/json/JsonBinaryAnnotation.scala
		value, err := base64str2sql(ann.Value)
		if err != nil {
			return fmt.Errorf("parse value: %+v", value)
		}

		ipv4, port, service := ep2sql(ann.Endpoint)
		_, err = stmt.Exec(
			hex2int64(span.TraceID),
			hex2int64(span.ID),
			ann.Key,
			value,
			ann.AnnotationType,
			span.Timestamp,
			ipv4, port, service,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
