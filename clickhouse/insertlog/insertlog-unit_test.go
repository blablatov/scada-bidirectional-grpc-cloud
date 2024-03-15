// Тестирование insert, без выполнения запроса
// go test -v .

package insertlog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var (
	insLog = "SELECT version()"
)

func TestInsertLogUnit(t *testing.T) {

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

	if conn != nil {
		t.Log("Cert get failed")
	}

	// Форматирование запроса
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", insLog)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	// Имитация запроса. Imitation of request
	// Сохранение и восстановление исходного значения respDo, с последующим восстановлением
	saved := respDo
	defer func() { respDo = saved }()

	var reqs *http.Request
	var method, url string
	Do := func(reqs *http.Request) error {
		method = "POST"
		url = "http://any"

		if method == "" && url == "" {
			t.Fatal("request is empty")
		}
		return nil
	}

	err = Do(reqs)
	if err != nil {
		t.Log(err)
	}
}
