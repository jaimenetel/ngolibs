package finderconnect

import (
	"encoding/json"
	"fmt"
	"sync"

	ut "github.com/jaimenetel/ngolibs/urltools"
)

type FinderConnect struct {
	ServerAddress string
	URLgetIp      string
	URLgetLtm     string
	URLgetDisp    string
	URLLineas     string
}

var (
	fcinstance   *FinderConnect
	fconce       sync.Once
	defaultURL   = "http://172.17.0.56:8701"
	defaultPaths = map[string]string{
		"URLgetIp":   "/findip?find=%s",
		"URLgetLtm":  "/findltm?find=%s",
		"URLgetDisp": "/finddispositivo?find=%s",
		"URLLineas":  "/lineas?imei=%s&userowner=%s",
	}
)

type LineaVerificaciones struct {
	Telefono  *string `json:"telefono"`
	ICCID     *string `json:"iccid"`
	IMEI      *string `json:"imei"`
	IDCliente *string `json:"idcliente"`
	Cliente   *string `json:"cliente"`
	IP        *string `json:"ip"`
}

type LineasResult struct {
	Result      string                `json:"result"`
	Explication string                `json:"explication"`
	Lineas      []LineaVerificaciones `json:"lineas"`
}

func makeFinderConnect(serverAddress string) FinderConnect {

	return FinderConnect{
		ServerAddress: serverAddress,
		URLgetIp:      serverAddress + defaultPaths["URLgetIp"],
		URLgetLtm:     serverAddress + defaultPaths["URLgetLtm"],
		URLgetDisp:    serverAddress + defaultPaths["URLgetDisp"],
		URLLineas:     serverAddress + defaultPaths["URLLineas"],
	}
}
func GetFinderConnect(serveraddress string) *FinderConnect {
	fconce.Do(func() {
		unaurl := defaultURL
		if serveraddress != "" {
			unaurl = serveraddress

		}
		apalo := makeFinderConnect(unaurl)
		fcinstance = &apalo

	})
	return fcinstance
}

// var aURLgetIp string = "http://172.17.0.56:8701/findip?find=%s"
// var aURLgetLtm string = "http://172.17.0.56:8701/findltm?find=%s"
// var aURLgetDisp string = "http://172.17.0.56:8701/finddispositivo?find=%s"

type Findiccid struct {
	Iccid  string `gorm:"column:iccid"`
	Tofind string `gorm:"column:tofind"`
	Tipo   string `gorm:"column:tipo"`
}

type _Cliente struct {
	Iccid   string `gorm:"column:iccid"`
	Cliente string `gorm:"column:cliente"`
	Codigo  string `gorm:"column:codigo"`
	Name    string `gorm:"column:name"`
}

type Dispositivo struct {
	ICCID   string   `json:"iccid"`
	IP      string   `json:"ip"`
	IMEI    string   `json:"imei"`
	Phone   string   `json:"phone"`
	LTM     string   `gorm:"column:ltm"`
	LTC     string   `gorm:"column:ltc"`
	Cliente _Cliente `json:"cliente,omitempty"`
}

// func FetchURL(url string) (string, error) {

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", fmt.Errorf("error al hacer la solicitud GET: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("respuesta fallida con código de estado: %d", resp.StatusCode)
// 	}

// 	// Lee el cuerpo de la respuesta
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
// 	}

// 	return string(body), nil
// }

func (fc *FinderConnect) GetIp(find string) (string, error) {
	URL := fmt.Sprintf(fc.URLgetIp, find)
	result, err := ut.FetchURL(URL)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (fc *FinderConnect) GetLTM(find string) ([]Findiccid, error) {
	URL := fmt.Sprintf(fc.URLgetLtm, find)
	result, err := ut.FetchURL(URL)
	if err != nil {
		return []Findiccid{}, err
	}
	fmt.Println("Result:", result)
	var resti []Findiccid

	err = json.Unmarshal([]byte(result), &resti)
	if err != nil {
		return []Findiccid{}, err
	}
	return resti, nil
}
func (fc *FinderConnect) GetDisp(find string) (Dispositivo, error) {
	URL := fmt.Sprintf(fc.URLgetDisp, find)
	result, err := ut.FetchURL(URL)
	if err != nil {
		return Dispositivo{}, err
	}
	fmt.Println("Result:", result)
	var dispositivo Dispositivo

	err = json.Unmarshal([]byte(result), &dispositivo)
	if err != nil {
		return Dispositivo{}, err
	}
	return dispositivo, nil
}

func (fc *FinderConnect) GetLineas(find string, userowner string) (LineasResult, error) {
	URL := fmt.Sprintf(fc.URLLineas, find, userowner)
	result, err := ut.FetchURL(URL)
	if err != nil {
		return LineasResult{}, err
	}
	fmt.Println("Result:", result)
	var lineasresult LineasResult

	err = json.Unmarshal([]byte(result), &lineasresult)
	if err != nil {
		return LineasResult{}, err
	}
	return lineasresult, nil
}
