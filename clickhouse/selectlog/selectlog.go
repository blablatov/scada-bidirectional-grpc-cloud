// Запрос select для выборки данных из таблицы СУБД ClickHouse ОП Yandex Cloud

package selectlog

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//"path/filepath"
	"sync"
	"time"
)

// Структура user_id
type CheckId struct {
	Mu          sync.Mutex
	Id          string
	Measurement float64
}

// var (
// 	crtFile   = filepath.Join("..", "certs", "YandexInternalRootCA.crt")
// )

// Метод подключения, аутентификации, выполнения запроса к БД grpcdb
func (c *CheckId) SelectLog(wg sync.WaitGroup, chs chan string, chm chan float64) error {

	log.SetPrefix("Select event: ")
	log.SetFlags(log.Lshortfile)

	defer wg.Done()

	// DSN для подключения к СУБД ClickHouse.
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "grpcdb"
	const DB_USER = "gogrpc"
	const DB_PASS = "rootroot"

	// Проверка существования Id в ClickHouse
	idGet := `SELECT id FROM grpcdb.grpc_log FINAL WHERE id = ` + <-chs //SELECT * FROM grpcdb.grpc_log FINAL WHERE id = 105;

	// Формирование метаданных запроса. Struct of request
	caCert, err := ioutil.ReadFile("YandexInternalRootCA.crt")
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

	// Формирование rest-запроса, его заголовков и тела
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", idGet)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	start := time.Now()

	// Выполнение запроса. Run request
	resp, err := conn.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	defer resp.Body.Close()

	fmt.Printf("Status = %v\n", resp.Status) // Статус ответа сервера

	// Чтение данных сервера, обработка ошибок
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Selected Id =", string(data))

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of request\n", secs)

	if string(data) != "" && string(data) != "0" {
		chs <- string(data)
	}

	// if strings.Contains(string(data), "Exception") {
	// 	log.Println("\nData is selected:", false)
	// 	chi <- false
	// } else {
	// 	log.Println("\nData is selected:", true)
	// 	chi <- true
	// }
	return nil
}

func (d *CheckId) listSql(db sql.DB, id string) sql.Result {
	res, err := db.Exec(`SELECT * FROM grpcdb.grpc_log FINAL WHERE id = ?`, id)
	if err != nil {
		log.Fatalln(err)
	}
	return res
}
