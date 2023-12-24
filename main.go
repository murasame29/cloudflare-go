package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/murasame29/cloudflare-go/client"
	"github.com/murasame29/cloudflare-go/r2"
)

const (
	Bucket = "r2-sample"
)

var (
	EndPoint        string = os.Getenv("PROVIDER")
	AccountID       string = os.Getenv("ACCOUNT_ID")
	AccessKeyID     string = os.Getenv("ACCESS_KEY_ID")
	accessKeySecret string = os.Getenv("ACCESS_KEY_SECRET")
)

func init() {
	flag.StringVar(&EndPoint, "endpoint", EndPoint, "endpoint")
	flag.StringVar(&AccountID, "account-id", AccountID, "account-id")
	flag.StringVar(&AccessKeyID, "access-key-id", AccessKeyID, "access-key-id")
	flag.StringVar(&accessKeySecret, "account-key-secret", accessKeySecret, "account-key-secret")
	flag.Parse()

	if EndPoint == "" || AccountID == "" || AccessKeyID == "" || accessKeySecret == "" {
		panic("missing required parameters")
	}
}

// 画像ファイルをバイト列に変換
func filToByte(file *os.File) ([]byte, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("file stat error :%v\n", err)
	}
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("file read error :%v\n", err)
	}

	return fileBytes, nil
}

func main() {
	client, err := client.New(
		AccountID,
		EndPoint,
		AccessKeyID,
		accessKeySecret,
	).Connect(context.TODO())
	if err != nil {
		log.Fatalf("r2 client conneciton error :%v\n", err)
	}

	r2 := r2.NewR2CRUD(Bucket, client, 60)

	// 画像ファイルを開いておく
	file, err := os.OpenFile("sample.png", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("file open error :%v\n", err)
	}

	filedata, err := filToByte(file)
	if err != nil {
		log.Fatalf("file transcate error :%v\n", err)
	}

	// 画像ファイルをアップロード
	if err := r2.UploadObject(context.Background(), filedata, "simple.png"); err != nil {
		log.Fatalf("r2 upload error :%v\n", err)
	}

	log.Println("upload success")

	// 全てのオブジェクトを取得
	objects, err := r2.ListObjects(context.Background(), "")
	if err != nil {
		log.Fatalf("r2 list objects error :%v\n", err)
	}

	for i, obj := range objects.Contents {
		log.Printf("%d: %s", i, *obj.Key) // オブジェクトのキー表示
	}

	// オブジェクトのURLを取得
	url, err := r2.PublishPresignedObjectURL(context.Background(), "simple.png")
	if err != nil {
		log.Fatalf("r2 publish presigned object url error :%v\n", err)
	}

	log.Println("sample.png Access URL :", url)

	// オブジェクトを削除
	if err := r2.DeleteObject(context.Background(), "simple.png"); err != nil {
		log.Fatalf("r2 delete object error :%v\n", err)
	}

	log.Println("delete success")
}
