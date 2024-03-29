// Запрос insert для вставки данных в таблицу СУБД ClickHouse ОП Yandex Cloud

package insertlog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type InsData struct {
	Mu          sync.Mutex
	Id          string
	Sensor      string
	Sensors     string
	Description string
	Destination string
	Measurement float64
}

// Файл сертификата
var (
	crtFile = filepath.Join("..", "certs", "YandexInternalRootCA.crt")
	req     *http.Request
	conn    *http.Client
)

// Для тестирования без отправки, логика запроса хранится в отдельной функции,
// в неэкспортируемой переменной уровня пакета. For later the test
var respDo = func(req *http.Request) error {
	resp, err := conn.Do(req)
	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Data inserted status: ", resp.Status)
	}
	defer resp.Body.Close()
	return nil
}

func (d *InsData) InsertLog(wg sync.WaitGroup, chs chan string, chm chan float64, chi chan bool) error {

	fm := strconv.FormatFloat(<-chm, 'g', -1, 64)
	fs := strings.ReplaceAll(fm, ".", ",")
	insLog := `INSERT INTO grpcdb.grpc_log (id, sensor, description, destination, measurement, timestamp) VALUES (` + <-chs + ` , 'Dallas", "Texas Instruments', 'sensor#99', 'Surgut, City', ` + fs + `)`

	log.SetPrefix("Insert event: ")
	log.SetFlags(log.Lshortfile)

	defer wg.Done()

	// DSN для подключения к СУБД ClickHouse.
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "grpcdb"
	const DB_USER = "gogrpc"
	const DB_PASS = "rootroot"

	// Формирование метаданных структуры запроса. Struct of request
	caCert, err := ioutil.ReadFile("YandexInternalRootCA.crt")
	if err != nil {
		log.Println(err)
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	conn = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// Форматирование запроса
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", insLog)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	start := time.Now()
	// Выполнение запроса. Run request
	respDo(req)

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of request\n", secs)

	return nil
}
