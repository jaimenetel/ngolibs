package urlcaller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	// URLBase is the base URL for the API
	URLIot  = "https://iot.liftel.es:8443/verif/executeurl"
	URLBeta = "https://beta.liftel.es:8443/verif/executeurl"
	URLLocal = "http://localhost:8703/executeurl"
)

type Prm struct {
	K string `json:"k"`
	V string `json:"v"`
}
type Auth []string
type ApiKey []string
type PrmList []Prm
type Body []string
type Method []string
type Object interface{}
type VeToken []bool

type UrlCall struct {
	UrlBase       string            `json:"urlbase"`
	Params        map[string]string `json:"params"`
	Authorization string            `json:"authorization"`
	ApiKey        string            `json:"apikey"`
	Body          string            `json:"body"`
	Method        string            `json:"method"`
	VToken        bool              `json:"vtoken"`
}
type UrlCaller struct {
	UrlExecuter string `json:"urlexecuter"`
}

var instance *UrlCaller
var once sync.Once

func GetUrlCaller(url ...string) *UrlCaller {
	urlcall := URLIot
	if len(url) > 0 {
		urlcall = url[0]
	}
	once.Do(func() {
		instance = &UrlCaller{UrlExecuter: urlcall}
	})
	return instance
}

func GetUrl(urlbase string, prmList ...interface{}) string {
	urlcall := UrlCall{UrlBase: urlbase, Params: map[string]string{}}

	// Iterate over the slice of Prm
	for _, param := range prmList {
		switch p := param.(type) {
		case Prm:
			prm := p
			urlcall.Params[prm.K] = prm.V

		case PrmList:
			prmList := p
			for _, prm := range prmList {
				urlcall.Params[prm.K] = prm.V
			}

		case Auth:
			urlcall.Authorization = p[0]
		case ApiKey:
			urlcall.ApiKey = p[0]
		case Body:
			urlcall.Body = p[0]
		case Method:
			urlcall.Method = p[0]
		case VeToken:
			urlcall.VToken = p[0]
		case Object:
			abody, _ := json.MarshalIndent(param, "", "  ")
			urlcall.Body = string(abody)

		}
	}
	ajson, err := json.MarshalIndent(urlcall, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	return string(ajson)
}

func (uc *UrlCaller) CallUrl(url string, params ...interface{}) string {
	abody := GetUrl(url, params...)
	fmt.Println(abody)
	salida, err := _FetchURLPost(context.Background(), uc.UrlExecuter, abody, "application/json")
	if err != nil {
		fmt.Println(err)
	}
	salida = ReformatJSON(salida)
	return salida

}

func _FetchURLPost(ctx context.Context, url string, body string, contentType string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud POST: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al hacer la solicitud POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta fallida con código de estado: %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
	}

	return string(responseBody), nil
}
func ReformatJSON(s string) string {
	var f interface{}
	err := json.Unmarshal([]byte(s), &f)
	if err != nil {
		fmt.Println(err)
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
