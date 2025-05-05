package s3

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const amzAlgorithm string = "AWS4-HMAC-SHA256"
const awsService string = "s3"
const dateTimeLayout string = "20060102T150405Z"

// S3 is responsible for signing and returning a presigned URL.
// For more details see: https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html.
type S3 struct {
	AccessID   string
	AWSRegion  string
	Expires    int
	AMZDate    string
	Host       string
	Secret     string
	BucketName string
}

// Sign and return a URL for the configured environment and the requested object.
func (s S3) Sign(method, objectName string) string {
	var amzDate = s.AMZDate
	if len(amzDate) == 0 {
		amzDate = getDateTime()
	}
	timeStamp := strings.Split(amzDate, "T")[0]

	amzCredential := strings.Join([]string{
		s.AccessID,
		timeStamp,
		s.AWSRegion,
		awsService,
		"aws4_request",
	}, "/")

	signedHeaders := []string{
		"host",
	}

	var canonicalURI string
	if len(s.BucketName) == 0 {
		canonicalURI = fmt.Sprintf("https://%s/%s", s.Host, objectName)
	} else {
		canonicalURI = fmt.Sprintf("https://%s/%s/%s", s.Host, s.BucketName, objectName)
	}

	canonicalQueryString := [][]interface{}{
		{"X-Amz-Algorithm", amzAlgorithm},
		{"X-Amz-Credential", amzCredential},
		{"X-Amz-Date", amzDate},
		{"X-Amz-Expires", fmt.Sprintf("%d", s.Expires)},
		{"X-Amz-SignedHeaders", formatSignedHeaders(signedHeaders)},
	}

	canonicalHeaders := []string{
		fmt.Sprintf("host:%s", s.Host),
	}

	if len(s.BucketName) != 0 {
		objectName = fmt.Sprintf("%s/%s", s.BucketName, objectName)
	}
	canonicalRequest := getCanonicalRequest(method, fmt.Sprintf("/%s", objectName), canonicalQueryString, canonicalHeaders, signedHeaders)

	crHash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		amzAlgorithm,
		amzDate,
		strings.Join([]string{
			timeStamp,
			s.AWSRegion,
			awsService,
			"aws4_request",
		}, "/"),
		hex.EncodeToString([]byte(crHash[:])),
	}, "\n")

	signature := sign(stringToSign, s.Secret, timeStamp, s.AWSRegion, awsService)

	query := strings.Join([]string{
		fmt.Sprintf("X-Amz-Algorithm=%s", amzAlgorithm),
		fmt.Sprintf("X-Amz-Credential=%s", url.QueryEscape(amzCredential)),
		fmt.Sprintf("X-Amz-Date=%s", amzDate),
		fmt.Sprintf("X-Amz-Expires=%d", s.Expires),
		fmt.Sprintf("X-Amz-SignedHeaders=%s", formatSignedHeaders(signedHeaders)),
		fmt.Sprintf("%s=%s", "X-Amz-Signature", signature),
	}, "&")

	url := fmt.Sprintf("%s?%s", canonicalURI, query)

	return url
}

func getDateTime() string {
	n := time.Now().UTC()
	return n.Format(dateTimeLayout)
}

func hmacSHA256(key []byte, msg []byte) []byte {
	h := hmac.New(sha256.New, (key))
	h.Write(msg)
	sig := h.Sum(nil)
	return sig
}

func sign(stringToSign, secret, amzDate, awsRegion, awsService string) string {
	dateKey := hmacSHA256([]byte(fmt.Sprintf("AWS4%s", secret)), []byte(amzDate)) // HMAC-SHA256("AWS4" + "<YourSecretAccessKey>","20130524")
	dateRegionKey := hmacSHA256(dateKey, []byte(awsRegion))                       // HMAC-SHA256(dateKey, "us-east-1")
	dateRegionServiceKey := hmacSHA256(dateRegionKey, []byte(awsService))         // HMAC-SHA256(dateRegionKey,"s3")
	signingKey := hmacSHA256(dateRegionServiceKey, []byte("aws4_request"))        // HMAC-SHA256(dateRegionServiceKey, "aws4_request")
	signature := hmacSHA256(signingKey, []byte(stringToSign))                     // Hex(HMAC-SHA256(signingKey, stringToSign))
	return hex.EncodeToString(signature)
}

func formatCanonicalQueryString(query [][]interface{}) string {
	params := []string{}
	for _, v := range query {
		if len(v) != 2 {
			continue
		}
		key := fmt.Sprintf("%v", v[0])
		value := fmt.Sprintf("%v", v[1])
		params = append(params, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
	}
	return strings.Join(params, "&")
}

func formatHeaders(headers []string) string {
	buf := bytes.Buffer{}
	for _, h := range headers {
		buf.WriteString(fmt.Sprintf("%s\n", strings.ToLower(h)))
	}
	return buf.String()
}

func formatSignedHeaders(signedHeaders []string) string {
	return strings.ToLower(strings.Join(signedHeaders, ";"))
}

func getCanonicalRequest(method string, path string, query [][]interface{}, headers []string, signedHeaders []string) string {
	return strings.Join([]string{
		strings.ToUpper(method),
		path,
		formatCanonicalQueryString(query),
		formatHeaders(headers),
		formatSignedHeaders(signedHeaders),
		"UNSIGNED-PAYLOAD",
	}, "\n")
}
