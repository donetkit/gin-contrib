package minio

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	client, err := New(WithEndpoint("192.168.5.110:9660"), WithAccess("zxiBHe4BQMMltM8N", "YdM4EDkz9oESbqZs6ayJWMzqpG5TqWvp"), WithEncryption(""), WithUseSSL(false))
	if err != nil {
		t.Error(err)
	}
	//filePath := "./2022-07-13_14-40-24.png"
	bucketName := "test"
	objectName := "123456.png"
	//objectName1 := "123445566.zip"  time.Now().Add(time.Second*60),

	//b, err := client.GetObjectByte(bucketName, objectName)
	//fmt.Println(len(b), err)
	//fmt.Println(client.FPutObject(bucketName, objectName, filePath, map[string]string{
	//	"File-Name": "Snipaste_2022-07-13_14-40-24.png",
	//}))

	//fmt.Println(client.PutObject(bucketName, objectName, nil, 0, map[string]string{
	//	"file_name": "111.rar",
	//}))

	//b, err = client.GetObjectByte(bucketName, objectName)
	//fmt.Println(len(b), err)
	info, _ := client.StatObject(bucketName, objectName)
	fmt.Println(info)

}
