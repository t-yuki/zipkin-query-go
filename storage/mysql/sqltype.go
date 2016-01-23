package mysql

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"net"
	"strconv"

	"github.com/t-yuki/zipkin-go/models"
)

func hex2int64(v *string) (n int64) {
	if v != nil {
		n, _ = strconv.ParseInt(*v, 16, 64)
	}
	return n
}

func bool2int(v *bool) (n int) {
	if v != nil && *v {
		return 1
	}
	return 0
}

func base64str2sql(v *string) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(*v)
	return b, err
}

func ep2sql(ep *models.Endpoint) (ipv4 sql.NullInt64, port sql.NullInt64, service sql.NullString) {
	if ep == nil {
		return
	}
	return ipv4str2sql(ep.IPV4), int2sql(ep.Port), string2sql(ep.ServiceName)
}

func ipv4str2sql(v *string) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	ip := net.ParseIP(*v)
	if ip == nil {
		return sql.NullInt64{}
	}
	ip = ip.To4()
	if ip == nil {
		return sql.NullInt64{}
	}
	val := binary.BigEndian.Uint32([]byte(ip))
	return sql.NullInt64{Int64: int64(val), Valid: true}
}

func int2sql(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}

func string2sql(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}
