package err

import "net/http"

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

var (
	RespDBDataInsertFailure = []byte(`{"error":"Failed to insert data"}`)
	RespDBDataAccessFailure = []byte(`{"error":"Failed to access data"}`)
	RespDBDataUpdateFailure = []byte(`{"error":"Failed to update data"}`)
	RespDBDataRemoveFailure = []byte(`{"error":"Failed to remove data"}`)

	RespJSONEncodeFailure = []byte(`{"error":"json encode failure"}`)
	RespJSONDecodeFailure = []byte(`{"error":"json decode failure"}`)

	RespInvalidURLParamID = []byte(`{"error":"invalid url param-id"}`)
)

func ServerError(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

func BadRequest(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}

func ValidationError(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(resp)
}
