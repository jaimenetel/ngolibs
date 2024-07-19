package urltools

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func FetchURLCtx(ctx context.Context, url string) (string, error) {
	return _FetchURL(ctx, url)
}
func _FetchURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud GET: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al hacer la solicitud GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta fallida con código de estado: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
	}

	return string(body), nil
}
func FetchURL(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return _FetchURL(ctx, url)
}
