package httptransport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
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

func readProto(r *http.Request, payload proto.Message) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return fmt.Errorf("empty request body")
	}
	return proto.Unmarshal(b, payload)
}

func readProtoAllowEmpty(r *http.Request, payload proto.Message) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return proto.Unmarshal(b, payload)
}

func writeAuto(w http.ResponseWriter, r *http.Request, status int, payload proto.Message) {
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		b, _ := protojson.Marshal(payload)
		w.Write(b)
		return
	}
	writeProto(w, status, payload)
}
