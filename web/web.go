package web

import (
  "net/http"
  "github.com/gorilla/mux"
)

func NewWeb() {
  muxRouter := mux.NewRouter()

  muxRouter.HandleFunc("/health", healthHandler)
  muxRouter.HandleFunc("/api/shorten", shortenURLHandler).
    Headers("Content-Type", "application/json").
    Methods("POST")
  muxRouter.HandleFunc("/api/resolve/{hash}", resolveURLHandler).
    Methods("GET")

  http.ListenAndServe(":8000", muxRouter)
}
