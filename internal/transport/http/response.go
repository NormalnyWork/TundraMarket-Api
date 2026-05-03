package httptransport

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeProto(w http.ResponseWriter, statusCode int, payload proto.Message) {
	b, err := proto.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/protobuf")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func writeProtoError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
