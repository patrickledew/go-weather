package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func get_api_into[T interface{}](result *T, urlFmt string, params ...any) (string, error) {
	url := fmt.Sprintf(urlFmt, params...);
	res, err := http.Get(url)
	if (err != nil) {
		fmt.Println("get_api_into: error getting", url)
		var body_str string
		res.Body.Read([]byte(body_str));
		return body_str, err
	}

	body, _ := io.ReadAll(res.Body)

	json.Unmarshal(body, result)

	return string(body), err
}