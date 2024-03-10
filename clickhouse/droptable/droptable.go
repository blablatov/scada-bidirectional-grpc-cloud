// Запрос drop для удаления таблицы СУБД ClickHouse ОП Yandex Cloud

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

// Файл сертификата. Структура запроса
var (
	crtFile   = filepath.Join("..", "certs", "YandexInternalRootCA.crt")
	dropTable = `DROP TABLE IF EXISTS grpcdb.grpc_log`
)

func main() {

	log.SetPrefix("Drop event: ")
	log.SetFlags(log.Lshortfile)

	// DSN для подключения к СУБД ClickHouse.
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "grpcdb"
	const DB_USER = "gogrpc"
	const DB_PASS = "rootroot"

	// Формирование метаданных структуры запроса. Struct of request
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

	// Форматирование запроса
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", dropTable)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	start := time.Now()
	// Выполнение запроса. Run request
	resp, err := conn.Do(req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Table drop status: ", resp.Status)
	}

	defer resp.Body.Close()

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of request\n", secs)
}
