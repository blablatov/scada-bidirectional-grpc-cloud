// Запрос create для создания таблицы СУБД ClickHouse ОП Yandex Cloud

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// Файл сертификата. Структура таблицы, с движком ReplacingMergeTree
var (
	crtFile  = filepath.Join("..", "certs", "YandexInternalRootCA.crt")
	createDB = `CREATE TABLE grpc_log
	(
    id UInt32,
    sensors String,
    description String,
    destination String,
    measurement Double,
    timestamp DateTime,
	)
	ENGINE = ReplacingMergeTree	
	PRIMARY KEY (id)`
)

func main() {
	// DSN для подключения к СУБД ClickHouse
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "grpcdb"
	const DB_USER = "gogrpc"
	const DB_PASS = "rootroot"

	// Формирование метаданных структуры запроса
	caCert, err := ioutil.ReadFile(crtFile)
	if err != nil {
		p := recover()
		log.Fatalln(err)
		panic(p)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	conn := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	// Формирование rest запроса, создание DBase
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	if err != nil {
		log.Fatalf("Error post request: %v", err)
	}
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", createDB)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	start := time.Now()
	// Выполнение запроса. Run request
	resp, err := conn.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Table created status: ", resp.Status)
	}

	defer resp.Body.Close()

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of request\n", secs)
}
