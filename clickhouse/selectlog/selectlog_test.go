// Тестирование select с выборкой Id из таблицы grpc_log СУБД ClickHouse ОП Yandex Cloud

package selectlog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestSelectLog(t *testing.T) {

	log.SetPrefix("Select-test event: ")
	log.SetFlags(log.Lshortfile)

	var tests = []struct {
		dbhost string
		dbname string
		dbuser string
		dbpass string
	}{
		{"rc1a-u620db3mp7svl13i.mdb.yandexcloud.net", "grpcdb", "gogrpc", "rootroot"},
		{"rc1a-u620db3mp7svl13i.mdb.yandexcloud.net", "grpcdb", "gogrpc", "rootroot"},
		{"rc1a-u620db3mp7svl13i.mdb.yandexcloud.net", "grpcdb", "gogrpc", "rootroot"},
	}

	var prev_dbhost string
	for _, test := range tests {
		if test.dbhost != prev_dbhost {
			log.Printf("%s\n", test.dbhost)
			prev_dbhost = test.dbhost
		}
	}

	var prev_dbname string
	for _, test := range tests {
		if test.dbname != prev_dbname {
			log.Printf("%s\n", test.dbname)
			prev_dbname = test.dbname
		}
	}
	var prev_dbuser string
	for _, test := range tests {
		if test.dbuser != prev_dbuser {
			log.Printf("%s\n", test.dbuser)
			prev_dbuser = test.dbuser
		}
	}
	var prev_dbpass string
	for _, test := range tests {
		if test.dbpass != prev_dbpass {
			log.Printf("%s\n", test.dbpass)
			prev_dbpass = test.dbpass
		}
	}

	chs := make(chan string, 1)
	chs <- "102"

	// Функциональный SQL запрос для получения Id
	idGet := `SELECT * FROM grpcdb.grpc_log WHERE id = ` + <-chs

	// Чтение сертификата из файла. Формирование метаданных запроса
	caCert, err := ioutil.ReadFile("YandexInternalRootCA.crt")
	if err != nil {
		// Восстанавливается для анализа, после вывода err, завершается
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
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:8443/", prev_dbhost), nil)
	query := req.URL.Query()
	query.Add("database", prev_dbname)
	query.Add("query", idGet)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", prev_dbuser)
	req.Header.Add("X-ClickHouse-Key", prev_dbpass)

	// Выполнение запроса
	resp, err := conn.Do(req)
	if err != nil {
		// Восстанавливается для анализа, после вывода err, завершается
		p := recover()
		log.Fatalln(err)
		panic(p)
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	defer resp.Body.Close()

	log.Printf("Status = %v ", resp.Status) // Статус ответа сервера

	// Чтение данных сервера, обработка ошибок
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error response:", err)
	}
	log.Println("Id successfully =", string(data))
}
