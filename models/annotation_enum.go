package models

// see: https://github.com/openzipkin/zipkin/blob/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCore.thrift#L214

type AnnotationType int64

const (
	AnnotationType_BOOL   AnnotationType = 0
	AnnotationType_BYTES  AnnotationType = 1
	AnnotationType_I16    AnnotationType = 2
	AnnotationType_I32    AnnotationType = 3
	AnnotationType_I64    AnnotationType = 4
	AnnotationType_DOUBLE AnnotationType = 5
	AnnotationType_STRING AnnotationType = 6
)

func (v AnnotationType) Addr() *int64 {
	i := int64(v)
	return &i
}

// handicraft enum definitions of Annotation.value
// see: https://github.com/openzipkin/zipkin/blob/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCore.thrift#L18

type AnnotationValue string

const (
	AnnotationValue_CLIENT_SEND          AnnotationValue = "cs"
	AnnotationValue_CLIENT_RECV          AnnotationValue = "cr"
	AnnotationValue_SERVER_SEND          AnnotationValue = "ss"
	AnnotationValue_SERVER_RECV          AnnotationValue = "sr"
	AnnotationValue_WIRE_SEND            AnnotationValue = "ws"
	AnnotationValue_WIRE_RECV            AnnotationValue = "wr"
	AnnotationValue_CLIENT_SEND_FRAGMENT AnnotationValue = "csf"
	AnnotationValue_CLIENT_RECV_FRAGMENT AnnotationValue = "crf"
	AnnotationValue_SERVER_SEND_FRAGMENT AnnotationValue = "ssf"
	AnnotationValue_SERVER_RECV_FRAGMENT AnnotationValue = "srf"
)

func (v AnnotationValue) Addr() *string {
	s := string(v)
	return &s
}

// handicraft enum definitions of BinaryAnnotation.key
// see: https://github.com/openzipkin/zipkin/blob/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCore.thrift#L110

type BinaryAnnotationKey string

const (
	BinaryAnnotationKey_LOCAL_COMPONENT BinaryAnnotationKey = "lc"
	BinaryAnnotationKey_CLIENT_ADDR     BinaryAnnotationKey = "ca"
	BinaryAnnotationKey_SERVER_ADDR     BinaryAnnotationKey = "sa"
)

func (v BinaryAnnotationKey) Addr() *string {
	s := string(v)
	return &s
}
