package verifloginconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	ut "github.com/jaimenetel/ngolibs/urltools"
)

const (
	URLIot   = "https://iot.liftel.es:8443/verif"
	URLLocal = "http://localhost:8703"
	URLBeta  = "https://beta.liftel.es:8443/verif"
)

type VerificacionesUser struct {
	Security    string `json:"security"`
	Usuario     string `json:"usuario"`
	IDEmpresa   int    `json:"idempresa"`
	Rol         int    `json:"rol"`
	TipoCliente string `json:"tipocliente"`
	Idioma      string `json:"idioma"`
	Token       string `json:"token"`
	Ce          int    `json:"ce"`
	AppToken    string `json:"apptoken"`
}
type TokenItem struct {
	Token       string `json:"token"`
	ValidoHasta string `json:"validohasta"`
	User        string `json:"user"`
	//Password    string `json:"password"`
	Vusuario VerificacionesUser `json:"vusuario"`
}
type VrfloginConnector struct {
	ServerAddress string
	URLGetToken   string
}

var (
	vlcinstance  *VrfloginConnector
	vlconce      sync.Once
	defaultURL   = "http://172.17.0.60:8703"
	defaultPaths = map[string]string{
		"URLGetToken": "/verificaciones/gettoken?user=%s&password=%s",
	}
)

func makeVerifConnect(serverAddress string) VrfloginConnector {

	return VrfloginConnector{
		ServerAddress: serverAddress,
		URLGetToken:   serverAddress + defaultPaths["URLGetToken"],
	}
}

func GetVerifConnect(serveraddress ...string) *VrfloginConnector {
	if serveraddress == nil {
		serveraddress = []string{URLIot}
	}
	vlconce.Do(func() {
		unaurl := defaultURL
		if serveraddress[0] != "" {
			unaurl = serveraddress[0]

		}
		apalo := makeVerifConnect(unaurl)
		vlcinstance = &apalo

	})
	return vlcinstance
}

func (fc *VrfloginConnector) GetToken(user, password string) (TokenItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := fc._GetToken(ctx, user, password)
	if err != nil {
		fmt.Println("Error al obtener el token1:", err)
		return TokenItem{}, err
	}

	return token, err
}

func (fc *VrfloginConnector) _GetToken(ctx context.Context, user, password string) (TokenItem, error) {
	URL := fmt.Sprintf(fc.URLGetToken, user, password)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Timeout de 30 segundos
	defer cancel()

	result, err := ut.FetchURLCtx(ctx, URL)
	if err != nil {
		fmt.Println("Error al obtener el token2:", err)
		return TokenItem{}, err
	}

	var token TokenItem
	err = json.Unmarshal([]byte(result), &token)
	if err != nil {
		return TokenItem{}, err
	}
	return token, nil
}
