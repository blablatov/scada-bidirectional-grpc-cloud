// Вставка данных лога через SQL запрос
// go test -v .

package insertlog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Структура тестового запроса. For test
var (
	sid       = fmt.Sprint(103)
	fm        = strconv.FormatFloat(49.7, 'g', -1, 64)
	fs        = strings.ReplaceAll(fm, ".", ",")
	insertLog = `INSERT INTO grpcdb.grpc_log (id, sensor, description, destination, measurement, timestamp) VALUES (` + sid + `, 'Dallas", "Texas Instruments', 'sensor#99', 'Surgut, City', ` + fs + `)`
)

func TestInsertLog(t *testing.T) {

	log.SetPrefix("Insert-test event: ")
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
	query.Add("query", insertLog)

	log.Println(insertLog)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	start := time.Now()
	// Выполнение запроса. Run request
	resp, err := conn.Do(req)
	if err != nil {
		p := recover()
		log.Fatalln(err)
		panic(p)
	}

	// Отложеное выполнение закрытия запроса, до получения ответа
	defer resp.Body.Close()

	// Чтение данных сервера, обработка ошибок
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Test insert is successfully")

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of request\n", secs)
}
