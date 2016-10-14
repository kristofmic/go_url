package web

import (
  "net/http"
  "encoding/json"
)

func healthHandler(w http.ResponseWriter, req *http.Request) {
  status, err := json.Marshal(map[string]string{
    "status": "alive",
  })
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(status)
}