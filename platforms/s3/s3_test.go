package s3

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var canonicalRequest string = strings.Join([]string{
	"GET",
	"/test.txt",
	"X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20130524%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20130524T000000Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host",
	"host:examplebucket.s3.amazonaws.com",
	"",
	"host",
	"UNSIGNED-PAYLOAD",
}, "\n")

var stringToSign string = strings.Join([]string{
	"AWS4-HMAC-SHA256",
	"20130524T000000Z",
	"20130524/us-east-1/s3/aws4_request",
	"3bfa292879f6447bbcda7001decf97f4a54dc650c8942174ae0a9121cf58ad04",
}, "\n")

func TestVerifyCanonicalRequest(t *testing.T) {
	expect := "3bfa292879f6447bbcda7001decf97f4a54dc650c8942174ae0a9121cf58ad04"
	got := fmt.Sprintf("%x", sha256.Sum256([]byte(canonicalRequest)))
	if strings.Compare(expect, got) != 0 {
		t.Errorf("got: %s; want: %s", got, expect)
	}
}

func TestFormatCanonicalQueryString(t *testing.T) {
	tests := []struct {
		query  [][]interface{}
		expect string
	}{
		{query: [][]interface{}{{"QueryParameter", "value"}}, expect: "QueryParameter=value"},
		{query: [][]interface{}{{"X-Amz-Credential", "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request"}}, expect: "X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20130524%2Fus-east-1%2Fs3%2Faws4_request"},
		{
			query: [][]interface{}{
				{"QueryParameter", "value"},
				{"QueryParameter2"},
				{"QueryParameter3", "value3", "value4"},
			}, expect: "QueryParameter=value",
		},
		{query: [][]interface{}{
			{"QueryParameter1", "value1"},
			{"QueryParameter2", "value2"},
		}, expect: "QueryParameter1=value1&QueryParameter2=value2"},
		{
			query: [][]interface{}{
				{"Query=Parameter", "key=value"},
			}, expect: "Query%3DParameter=key%3Dvalue"},
	}

	for i, tc := range tests {
		got := formatCanonicalQueryString(tc.query)
		if strings.Compare(tc.expect, got) != 0 {
			t.Errorf("test %d: expected: %v, got: %v", i, tc.expect, got)
		}
	}
}

func TestFormatHeaders(t *testing.T) {
	headers := []string{
		"Host:examplebucket.s3.amazonaws.com",
	}
	expect := "host:examplebucket.s3.amazonaws.com\n"
	got := formatHeaders(headers)
	if strings.Compare(expect, got) != 0 {
		t.Errorf("got: %s;\nwant: %s", got, expect)
	}
}

func TestFormatSignedHeaders(t *testing.T) {
	tests := []struct {
		headers []string
		expect  string
	}{
		{headers: []string{"Header"}, expect: "header"},
		{headers: []string{"Header1", "Header2"}, expect: "header1;header2"},
	}

	for i, tc := range tests {
		got := formatSignedHeaders(tc.headers)
		if strings.Compare(tc.expect, got) != 0 {
			t.Errorf("test %d: expected: %v, got: %v", i, tc.expect, got)
		}
	}
}

func TestGetCanonicalRequest(t *testing.T) {
	method := "GET"
	path := "/test.txt"
	query := [][]interface{}{
		{"X-Amz-Algorithm", "AWS4-HMAC-SHA256"},
		{"X-Amz-Credential", "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request"},
		{"X-Amz-Date", "20130524T000000Z"},
		{"X-Amz-Expires", 86400},
		{"X-Amz-SignedHeaders", "host"},
	}
	headers := []string{
		"host:examplebucket.s3.amazonaws.com",
	}
	signedHeaders := []string{
		"host",
	}

	expect := canonicalRequest
	got := getCanonicalRequest(method, path, query, headers, signedHeaders)
	if strings.Compare(expect, got) != 0 {
		t.Errorf("got: %s;\nwant: %s", expect, got)
	}
}

func TestSign(t *testing.T) {
	expect := "aeeed9bbccd4d02ee5c0109b86d86835f995330da4c265957d157751f604d404"

	secret := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	amzDate := "20130524"
	awsRegion := "us-east-1"
	awsService := "s3"

	got := sign(stringToSign, secret, amzDate, awsRegion, awsService)

	if strings.Compare(expect, got) != 0 {
		t.Errorf("expected: %v, got: %v", expect, got)
	}
}

func TestSignedURL(t *testing.T) {
	expect := strings.Join([]string{
		"https://examplebucket.s3.amazonaws.com/test.txt",
		"?X-Amz-Algorithm=AWS4-HMAC-SHA256",
		"&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20130524%2Fus-east-1%2Fs3%2Faws4_request",
		"&X-Amz-Date=20130524T000000Z",
		"&X-Amz-Expires=86400",
		"&X-Amz-SignedHeaders=host",
		"&X-Amz-Signature=aeeed9bbccd4d02ee5c0109b86d86835f995330da4c265957d157751f604d404",
	}, "")
	s := S3{
		AccessID:  "AKIAIOSFODNN7EXAMPLE",
		AWSRegion: "us-east-1",
		AMZDate:   "20130524T000000Z",
		Secret:    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		Host:      "examplebucket.s3.amazonaws.com",
		Expires:   86400,
	}
	method := "GET"
	objectName := "test.txt"
	got := s.Sign(method, objectName)
	if got != expect {
		t.Errorf("url = %s; want %s", got, expect)
	}
}

func TestSignedURL2(t *testing.T) {
	expect := strings.Join([]string{
		"https://s3.amazonaws.com/examplebucket/test.txt",
		"?X-Amz-Algorithm=AWS4-HMAC-SHA256",
		"&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20130524%2Fus-east-1%2Fs3%2Faws4_request",
		"&X-Amz-Date=20130524T000000Z",
		"&X-Amz-Expires=86400",
		"&X-Amz-SignedHeaders=host",
		"&X-Amz-Signature=733255ef022bec3f2a8701cd61d4b371f3f28c9f193a1f02279211d48d5193d7",
	}, "")
	s := S3{
		AccessID:   "AKIAIOSFODNN7EXAMPLE",
		BucketName: "examplebucket",
		AWSRegion:  "us-east-1",
		AMZDate:    "20130524T000000Z",
		Secret:     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		Host:       "s3.amazonaws.com",
		Expires:    86400,
	}
	method := "GET"
	path := "test.txt"
	got := s.Sign(method, path)
	if got != expect {
		t.Errorf("url = %s; want %s", got, expect)
	}
}

func TestSignedURL3(t *testing.T) {
	expect := "https://s3.amazonaws.com/examplebucket/test.txt"
	s := S3{
		AccessID:   "AKIAIOSFODNN7EXAMPLE",
		BucketName: "examplebucket",
		AWSRegion:  "us-east-1",
		Secret:     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		Host:       "s3.amazonaws.com",
		Expires:    86400,
	}
	method := "GET"
	path := "test.txt"
	got := s.Sign(method, path)
	if strings.HasPrefix(got, expect) == false {
		t.Errorf("url = %s; want %s", got, expect)
	}
}

func TestGetDateTimeFormat(t *testing.T) {
	got := getDateTime()
	expectedPattern := "^[0-9]{8}T[0-9]{6}Z$"
	match, _ := regexp.MatchString(expectedPattern, got)
	if !match {
		t.Errorf("expected to match: %v, got: %v", expectedPattern, got)
	}
}
