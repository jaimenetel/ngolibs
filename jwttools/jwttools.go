package jwttools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var SUPER_SECRET_KEY = ""

func GetSecretKey() string {

	value := os.Getenv("VARTOKEN")
	if value == "" {
		fmt.Println("T K NOT FOUND")
	}
	fmt.Println("T K FOUND", value)
	return value

}
func getUserName(tokenString string) (string, error) {
	if tokenString == "" {
		return "", nil
	}

	// Limpiar el token para remover 'Bearer' si está presente
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// Definir la clave secreta
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecretKey()), nil
	}

	// Parsear y validar el token
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extraer el usuario del campo 'sub'
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", nil
}

func DecodificarJWT2(tokenString string) (jwt.MapClaims, error) {
	//secretKeySeed := SUPER_SECRET_KEY
	//secretKey := hmac.New(sha256.New, []byte(secretKeySeed))
	// Decodificar el token
	atoken := strings.Replace(tokenString, "Bearer", "", -1)
	atoken = strings.TrimSpace(atoken)
	token, err := jwt.Parse(atoken, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetSecretKey()), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al decodificar el token: %v", err)
	}

	// Validar el token
	if !token.Valid {
		return nil, fmt.Errorf("token no válido")
	}

	// Obtener el contenido del token (payload)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error al obtener las claims del token")
	}
	fmt.Println(claims)
	return claims, nil
}
func GetExpirationTime(token string) int64 {
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	if parts, ok := myClaims["exp"].(float64); ok {
		return int64(parts)
	}
	return 0
}
func GetExpirationTimeAsString(token string) string {
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	if parts, ok := myClaims["exp"].(float64); ok {
		return fmt.Sprintf("%v", parts)
	}
	return ""
}
func GetUserFromBearerToken(token string) string {
	// Split the token by the space character
	//name, _ := getUserName(token)
	//fmt.Println("name:", name)
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	if parts, ok := myClaims["sub"].(string); ok && parts != "" {
		return parts
	}
	return ""

}
func GetVerifUserFromBearerToken(token string) string {

	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	if parts, ok := myClaims["cliuser"].(string); ok && parts != "" {
		return parts
	} else if parts, ok := myClaims["sub"].(string); ok && parts != "" {
		return parts
	}
	return ""
}
func GetRolesFromBearerToken(token string) []string {
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
		return nil
	}

	parts, ok := myClaims["cliuser"].(string)
	if !ok || parts == "" {
		parts, ok = myClaims["sub"].(string)
		if !ok {
			fmt.Println("Error: 'sub' claim is missing or not a string")
			return nil
		}
	}

	role, ok := myClaims["role"].(string)
	if !ok {
		fmt.Println("Error: 'role' claim is missing or not a string")
		return nil
	}

	roles := strings.Split(role, ",")
	return roles
}
func GetRolesFromBearerTokenString(token string) string {

	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error al decodificar el token:", err)
	}
	parts := myClaims["cliuser"].(string)
	if parts == "" {
		parts = myClaims["sub"].(string)
	}

	role := "," + myClaims["role"].(string) + ","

	return role
}

// CompareRolesWithToken comprueba si hay intersecciones entre roles de una cadena y una lista de cadenas
func CompareRolesWithToken(roleString, token string) bool {
	roles := GetRolesFromBearerToken(token)
	if roles == nil {
		return false
	}
	// Divide la cadena de roles en una lista de cadenas
	roleList := strings.Split(roleString, ",")
	// Comprobar intersección de lista
	for _, role := range roleList {
		for _, r := range roles {
			if strings.TrimSpace(role) == strings.TrimSpace(r) {
				return true
			}
		}
	}
	return false
}

// Funcion recibe (token, user(payload string), roles...) y devuelve bool comprobando si tokenDecodificado.payload == user OOO token.roles == cualquier rol que acabamos de pasar --> True
// CheckUserWithTokenOrRoles checks if the user payload or roles match with the token's claims
// Ejemplo: user = `{"exp":1.112123417e+09,"hostname":"","iat":1.112123417e+09,"iss":"Liftel.es","role":"ROLE_ADMIN","sub":"raul"}`
// var token = "eyJhbGci2iJIUzUxMiJ9.eyJzdWIiOiJSQVVMQUQiLCJyb2xlIjoiUk9MRV9BRE1JTixST0xFX0dFU1RJT04sUk9MRV9NQVRJQyxST0xFX1JFR05FVEVMIiwiaG9zdG5hbWUiOiJQLU1BVElDIiwiaXNzIjoiTGlmdGVsLmVzIiwiaWF0IjoxNzE5NDk3Nzc3LCJleHAiOjE3MTk1ODQxNzd9.5rkX99zpjSWLXVjqOfs4fM-3cAoUSiq2xZ9vj6AUDmpHEu4TFZ8R5PUZxlsnMpkvxkI6CD60HqkIfNvczpWECA"
func CheckUserWithTokenOrRoles(token string, user string, roles ...string) (bool, error) {
	myClaims, err := DecodificarJWT2(token)
	if err != nil {
		fmt.Println("Error decoding the token:", err)
		return false, err
	}

	fmt.Println("Decoded Claims:", myClaims)

	// Convert user (string) to jwt.MapClaims
	var userClaims jwt.MapClaims
	err = json.Unmarshal([]byte(user), &userClaims)
	if err != nil {
		fmt.Println("Error deserializing the user (FORMAT JSON):", err)
		// return false, err
	} else {
		fmt.Println("User Claims:", userClaims)

		// Check if user claims match the token claims
		match := true
		for key, value := range userClaims {
			if myClaims[key] != value {
				match = false
				break
			}
		}
		if match {
			return true, nil
		}
	}

	// Check if any of the roles match
	tokenRoles := strings.Split(myClaims["role"].(string), ",")
	for _, role := range roles {
		if contains(tokenRoles, role) {
			return true, nil
		}
	}

	return false, nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func init() {

}
