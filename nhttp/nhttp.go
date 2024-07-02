package nhtttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	jwttools "github.com/jaimenetel/ngolibs/jwttools"
)

// inmersion
type TinTemas map[string]struct{}

type Roles []string
type AntiRoles []string
type NoRoles []bool
type Methods []string
type QParamsAll []string
type QParamsAny []string
type ApiKey []bool
type LogUse []bool
type StringFunc func(string)
type CheckFunc func(string) bool

type MyHandlerFunc struct {
	Funcname string           `json:"funcname"`
	Func     http.HandlerFunc `json:"-"`
}

type SaveLogFunc struct {
	Funcname string     `json:"funcname"`
	Func     StringFunc `json:"-"`
}

type CheckApiFunc struct {
	Funcname string    `json:"funcname"`
	Func     CheckFunc `json:"-"`
}
type Container struct {
	Elementos []interface{} `json:"elementos"`
}

type Endpoint struct {
	Name         string        `json:"endpointname"`
	MyHandler    MyHandlerFunc `json:"handler"`
	Controller   string        `json:"controller"`
	Roles        Roles         `json:"roles"`
	AntiRoles    AntiRoles     `json:"antiroles"`
	QParamsAll   QParamsAll    `json:"qparamsall"`
	QParamsAny   QParamsAny    `json:"qparamsany"`
	Methods      Methods       `json:"methods"`
	InMethods    TinTemas      `json:"inmethods"`
	InRoles      TinTemas      `json:"inroles"`
	InAntiRoles  TinTemas      `json:"inantiroles"`
	InQParamsAll TinTemas      `json:"inqparamsall"`
	InQParamsAny TinTemas      `json:"inqparamsany"`
	SaveLog      SaveLogFunc   `json:"savelog"`
	CheckApi     CheckApiFunc  `json:"checkapi"`
	SinRoles     bool          `json:"sinroles"`
	ApiKey       bool          `json:"apikey"`
	LogUse       bool          `json:"loguse"`
}

func TinTemasToString(tin TinTemas) string {
	var s string
	for tema := range tin {
		s += tema + ","
	}
	return strings.TrimSuffix(s, ",")
}

type Nthttp struct {
	Port      string
	Endpoints []Endpoint
}
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Generar la respuesta en formato JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

var httpinstance *Nthttp
var oncelt sync.Once

func Ntinstance() *Nthttp {
	oncelt.Do(func() {
		httpinstance = &Nthttp{}

	})
	return httpinstance
}
func ProcessParametros(params []interface{}) Container {

	UnContainer := Container{}
	for _, param := range params {
		switch p := param.(type) {
		case CheckApiFunc:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case SaveLogFunc:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case LogUse:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case NoRoles:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case ApiKey:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case Roles:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case AntiRoles:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case QParamsAll:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case QParamsAny:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case Methods:
			UnContainer.Elementos = append(UnContainer.Elementos, p)
		case Container:
			AnotherContainer := ProcessParametros(p.Elementos)
			UnContainer.Elementos = append(UnContainer.Elementos, AnotherContainer.Elementos...)
		default:
			fmt.Println("No sé qué es")
			typeOfParam := reflect.TypeOf(param)
			fmt.Printf("Type: %v\n", typeOfParam)
		}
	}

	return UnContainer
}
func (lt *Nthttp) AddEndpoint(name string, handler http.HandlerFunc, params ...interface{}) Endpoint {

	epRoles := Roles{}
	epAntiRoles := AntiRoles{}
	epQParamsAll := QParamsAll{}
	epQParamsAny := QParamsAny{}
	epMethods := Methods{}
	var tinMethods TinTemas
	var tinRoles TinTemas
	var tinAntiRoles TinTemas
	var tinQParamsAll TinTemas
	var tinQParamsAny TinTemas
	var epSinRoles = false
	var epApiKey = false
	var epLogUse = false
	var epSaveLogFunc = SaveLogFunc{Funcname: "No activo", Func: nil}
	var epCheckApiFunc = CheckApiFunc{Funcname: "No activo", Func: nil}
	var epHandlerFunc = MyHandlerFunc{Funcname: runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), Func: handler}
	AContainer := ProcessParametros(params)
	for _, param := range AContainer.Elementos {
		switch p := param.(type) {
		case CheckApiFunc:
			fmt.Println("Es un CheckApiFunc")
			//epCheckApiFunc = p
			afunc := p.Func
			aname := runtime.FuncForPC(reflect.ValueOf(p.Func).Pointer()).Name()
			fmt.Println("Nombre de la función: ", aname)
			epCheckApiFunc = CheckApiFunc{Funcname: aname, Func: afunc}
		case SaveLogFunc:
			fmt.Println("Es un SaveLogFunc")
			afunc := p.Func
			aname := runtime.FuncForPC(reflect.ValueOf(p.Func).Pointer()).Name()
			epSaveLogFunc = SaveLogFunc{Funcname: aname, Func: afunc}
		case LogUse:
			if len(param.(LogUse)) > 0 {
				epLogUse = p[0]
			}
			fmt.Println("Activando LogUse", epLogUse)

		case NoRoles:
			if len(param.(NoRoles)) > 0 {
				epSinRoles = p[0]
			}
			fmt.Println("Activando NoRoles", epSinRoles)
		case ApiKey:
			if len(param.(ApiKey)) > 0 {
				epApiKey = p[0]
			}

			fmt.Println("Activando ApiKey", epApiKey)

		case Roles:
			epRoles = append(epRoles, p...)
			tinRoles = make(map[string]struct{})
			tinRoles["ROLE_ALL"] = struct{}{}
			tinRoles["ROLE_ADMIN"] = struct{}{}
			for _, role := range epRoles {
				tinRoles[role] = struct{}{}
			}

			fmt.Println("Activando Roles", epRoles)
		case AntiRoles:
			epAntiRoles = append(epAntiRoles, p...)
			tinAntiRoles = make(map[string]struct{})
			for _, role := range epAntiRoles {
				tinAntiRoles[role] = struct{}{}
			}
			fmt.Println("Activando AntiRoles", epAntiRoles)
		case QParamsAll:
			epQParamsAll = append(epQParamsAll, p...)
			tinQParamsAll = make(map[string]struct{})
			for _, role := range epQParamsAll {
				tinQParamsAll[role] = struct{}{}
			}
			fmt.Println("Activando QParamsAll", epQParamsAll)
		case QParamsAny:
			epQParamsAny = append(epQParamsAny, p...)
			tinQParamsAny = make(map[string]struct{})
			for _, role := range epQParamsAny {
				tinQParamsAny[role] = struct{}{}
			}
			fmt.Println("Activando QParamsAny", epQParamsAny)
		case Methods:
			epMethods = append(epMethods, p...)
			tinMethods = make(map[string]struct{})
			for _, method := range epMethods {
				tinMethods[method] = struct{}{}
			}
			fmt.Println("Es un tipo Methods", epMethods)
		default:
			fmt.Println("No sé qué es")
			typeOfParam := reflect.TypeOf(p)
			fmt.Printf("Type: %v\n", typeOfParam)
		}
	}

	endpoint := Endpoint{
		Name:         name,
		MyHandler:    epHandlerFunc,
		Roles:        epRoles,
		AntiRoles:    epAntiRoles,
		QParamsAll:   epQParamsAll,
		QParamsAny:   epQParamsAny,
		Methods:      epMethods,
		InMethods:    tinMethods,
		InRoles:      tinRoles,
		InAntiRoles:  tinAntiRoles,
		InQParamsAll: tinQParamsAll,
		InQParamsAny: tinQParamsAny,
		SinRoles:     epSinRoles,
		ApiKey:       epApiKey,
		LogUse:       epLogUse,
		SaveLog:      epSaveLogFunc,
		CheckApi:     epCheckApiFunc,
	}

	lt.Endpoints = append(lt.Endpoints, endpoint)
	return endpoint
}
func (lt *Nthttp) Start() {
	for _, endpoint := range lt.Endpoints {
		fmt.Println("Startr del endpoint: ", endpoint)

		// Envuelve el handler original con los middlewares de auth y log, y luego con el CORS middleware
		handlerWithMiddleware := corsMiddleware(authMiddlewareRoleLog(endpoint.MyHandler.Func, endpoint))

		handlerWithMiddleware = ConfigMethodType(handlerWithMiddleware, endpoint.InMethods)

		http.Handle(endpoint.Name, handlerWithMiddleware)
	}
}
func ConfigMethodType(next http.Handler, methods TinTemas) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, found := methods[r.Method]
		if !found {

			pp := ""
			for method := range methods {
				pp += method + ", "
			}
			pp = strings.TrimSuffix(pp, ", ")
			RespondWithError(w, http.StatusMethodNotAllowed, "Only "+pp+" is supported")
			return
		}

		next.ServeHTTP(w, r)
	})
}
func corsMiddleware(next http.Handler) http.Handler {
	fmt.Println("CORS Middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
func CheckStringSliceInTemas(slice []string, temas TinTemas) (bool, string) {
	for _, s := range slice {
		_, found := temas[s]
		if found {
			return true, s
		}
	}
	return false, ""
}
func (ep *Endpoint) CheckAPIKey(r *http.Request) bool {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey != "" {
		if ep.CheckApi.Func != nil {
			return ep.CheckApi.Func("API Key: " + apiKey)
		}

		return true
	}
	return false
}
func DoLogUse(next http.Handler, ep Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ep.LogUse {
			datos := imprimirDatosSolicitud(r)
			if ep.SaveLog.Func != nil {
				ep.SaveLog.Func("Guardando log: " + datos)
			}

		}
		startTime := time.Now()
		next.ServeHTTP(w, r)
		fmt.Println("Después de llamar al siguiente")
		duration := time.Since(startTime)
		fmt.Println("Tiempo de ejecución:", duration)

	})
}
func authMiddlewareRoleLog(next http.Handler, ep Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearer_string := "Bearer"
		//imprimirDatosSolicitud(w, r)
		fmt.Println("tras solicitud", TinTemasToString(ep.InRoles))
		if ep.ApiKey {
			if !ep.CheckAPIKey(r) {
				RespondWithError(w, http.StatusUnauthorized, "API Key no proporcionado o no autorizado")
				return
			}
		}
		tokenString := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), bearer_string, "", -1))
		if ep.SinRoles {
			fmt.Println("Sin roles")
			next.ServeHTTP(w, r)
		} else {

			if tokenString == "" {

				RespondWithError(w, http.StatusUnauthorized, "Token JWT no proporcionado")
				return
			}

			fmt.Println("con token:", tokenString)

			fmt.Println("token valido")
			// Verifica el rol del usuario
			myClaims, err := jwttools.DecodificarJWT2(tokenString)
			if err != nil {
				// http.Error(w, "Token JWT no válido", http.StatusUnauthorized)
				RespondWithError(w, http.StatusUnauthorized, "Token JWT no válido")
				return
			}

			fmt.Println("tras claims")
			//role := logrequest.Claims["role"].(string)
			role := myClaims["role"].(string)
			roles := strings.Split(role, ",")
			fmt.Println("Roles: ", role)
			pongoRoles, cualpone := CheckStringSliceInTemas(roles, ep.InRoles)
			quitoRoles, cualquita := CheckStringSliceInTemas(roles, ep.InAntiRoles)
			if quitoRoles {
				RespondWithError(w, http.StatusForbidden, "Acceso no autorizado por "+cualquita)
				return
			}
			if !pongoRoles {
				RespondWithError(w, http.StatusForbidden, "Acceso no autorizado")
				return
			}

			fmt.Println("tenemos roles", cualpone)

			if len(ep.InQParamsAll) > 0 {
				fmt.Println("QParamsAll")
				// Verificar que todos los parámetros estén presentes
				for param := range ep.InQParamsAll {
					apar := r.URL.Query().Get(param)
					if apar == "" {
						RespondWithError(w, http.StatusBadRequest, "Falta el parámetro "+param)
						return
					}
				}
			}
			tenemosparams := false
			if len(ep.InQParamsAny) > 0 {
				fmt.Println("QParamsAny")
				// Verificar que al menos uno de los parámetros esté presente
				for param := range ep.InQParamsAny {
					if r.FormValue(param) != "" {
						tenemosparams = true
						break
					}
				}
				if !tenemosparams {
					RespondWithError(w, http.StatusBadRequest, "Falta al menos uno de los parámetros")
					return
				}
			}
			DoLogUse(next, ep).ServeHTTP(w, r)
			//next.ServeHTTP(w, r)
		}
	})
}
func imprimirDatosSolicitud(r *http.Request) string {
	stbuilder := strings.Builder{}
	stbuilder.WriteString("Método de solicitud: " + r.Method + "\n")

	token := r.Header.Get("Authorization")
	stbuilder.WriteString("Token de autorización: " + token + "\n")

	user := jwttools.GetUserFromBearerToken(token)
	stbuilder.WriteString("Usuario: " + user + "\n")

	stbuilder.WriteString("URL solicitada: " + r.URL.String() + "\n")

	roles := jwttools.GetRolesFromBearerToken(token)
	rolestring := strings.Join(roles, ",")
	stbuilder.WriteString("Roles: " + rolestring + "\n")

	// Imprime los encabezados de la solicitud
	stbuilder.WriteString("Encabezados de la solicitud:\n")

	for nombre, valores := range r.Header {
		for _, valor := range valores {
			stbuilder.WriteString(nombre + ": " + valor + "\n")

		}
	}

	queryParams := r.URL.Query()
	stbuilder.WriteString("Parámetros de consulta:\n")

	for nombre, valores := range queryParams {
		for _, valor := range valores {
			stbuilder.WriteString(nombre + ": " + valor + "\n")

		}
	}
	return stbuilder.String()
}
