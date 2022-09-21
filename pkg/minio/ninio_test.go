package minio

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	client, err := New(WithEndpoint("192.168.5.110:9660"), WithAccess("umTp5YVVCWjeJG2V", "SaXvSRxhvcNw0WLSdXDPhAZJ6e9Irux4"), WithEncryption(""), WithUseSSL(false))
	if err != nil {
		t.Error(err)
	}
	filePath := "./111.rar"
	bucketName := "2022124"
	objectName := "111111111111WithEncryption.rar"
	//objectName1 := "123445566.zip"  time.Now().Add(time.Second*60),

	b, err := client.GetObjectByte(bucketName, objectName)
	fmt.Println(len(b), err)
	fmt.Println(client.FPutObject(bucketName, objectName, filePath, map[string]string{
		"file_name": "111.rar",
	}))

	fmt.Println(client.PutObject(bucketName, objectName, nil, 0, map[string]string{
		"file_name": "111.rar",
	}))

	b, err = client.GetObjectByte(bucketName, objectName)
	fmt.Println(len(b), err)

	fmt.Println(client.StatObject(bucketName, objectName))

}
