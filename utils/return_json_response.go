package utils

import (
	"encoding/json"
	"net/http"
)

func ReturnJSONResponse(
	writer http.ResponseWriter,
	code uint16,
	payload any,
) {
	writer.Header().Set("Content-Type", "application/json")

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)

		if err != nil {
			writer.WriteHeader(500)

			return
		}

		defer writer.Write(jsonPayload)
	}

	writer.WriteHeader(int(code))
}
