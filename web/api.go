package web

import (
  "net/http"
  "encoding/hex"
  "encoding/json"
  "crypto/md5"
  "github.com/gorilla/mux"
  "github.com/kristofmic/go_url/sources/redis"
)

type resolveURLApiRes struct {
  OriginalURL string `json:"original_url"`
}

type shortenURLApiReq struct {
  URL string `json:"url"`
}

type shortenedURL struct {
  Hash string `json:"hash"`
  OriginalURL string `json:"original_url"`
}

func resolveURLHandler(res http.ResponseWriter, req *http.Request) {
  hash := mux.Vars(req)["hash"]

  client := redis.GetClient()
  originalURL, getErr := client.HGet(hash, "OriginalURL").Result()

  if getErr != nil || originalURL == "" {
    http.Error(res, "Short URL not found", http.StatusNotFound)
    return
  }


  resData, marshalJSONErr := json.Marshal(resolveURLApiRes{originalURL})
  if marshalJSONErr != nil {
    http.Error(res, marshalJSONErr.Error(), http.StatusInternalServerError)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(http.StatusOK)
  res.Write(resData)
}

func shortenURLHandler(res http.ResponseWriter, req *http.Request) {
  if req.Body == nil {
    http.Error(res, "Please send a valid URL", http.StatusBadRequest)
    return
  }

  var reqData shortenURLApiReq
  decodeJSONErr := json.NewDecoder(req.Body).Decode(&reqData)
  if decodeJSONErr != nil {
    http.Error(res, decodeJSONErr.Error(), http.StatusInternalServerError)
    return
  }

  hash := md5.Sum([]byte(reqData.URL))
  encodedShortHash := hex.EncodeToString(hash[:3])
  model := shortenedURL{
    Hash: encodedShortHash,
    OriginalURL: reqData.URL,
  }

  client := redis.GetClient()
  setErr := client.HMSet(encodedShortHash, map[string]string{
    "Hash": model.Hash,
    "OriginalURL": model.OriginalURL,
  }).Err()
  if setErr != nil {
    http.Error(res, setErr.Error(), http.StatusInternalServerError)
    return
  }

  resData, marshalJSONErr := json.Marshal(model)
  if marshalJSONErr != nil {
    http.Error(res, marshalJSONErr.Error(), http.StatusInternalServerError)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(http.StatusOK)
  res.Write(resData)
}

