package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ToJson(res http.ResponseWriter, status int, dto interface{}) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)

	if err := encoder.Encode(dto); err != nil {
		http.Error(res, "Fail on serialization process.", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Server", "Account/1.0")
	res.WriteHeader(status)
	res.Write(buf.Bytes())
}

func BindJson(source io.Reader, data interface{}) error {
	decoder := json.NewDecoder(source)
	return decoder.Decode(data)
}
