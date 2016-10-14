package web

import (
  "net/http"
  "encoding/hex"
  "encoding/json"
  "crypto/md5"
)

type apiShortenURL struct {
  URL string
}

type shortenedURL struct {
  ID uint64
  Hash string
  OriginalURL string
}

var urls = make(map[uint64]*shortenedURL)
var ids uint64

func hasURL(collection *map[uint64]*shortenedURL, URL string) (has bool, ID uint64) {
  for key, val := range *collection {
    if val.OriginalURL == URL {
      has = true
      ID = key
      break;
    }
  }

  return has, ID
}

func shortenURLHandler(res http.ResponseWriter, req *http.Request) {
  if req.Body == nil {
    http.Error(res, "Please send a valid URL", http.StatusBadRequest)
    return
  }

  var reqData apiShortenURL
  err := json.NewDecoder(req.Body).Decode(&reqData)
  if err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }

  ok, existingID := hasURL(&urls, reqData.URL)
  var shortenedStruct shortenedURL

  if (!ok) {
    ids++
    urlID := ids
    hash := md5.Sum([]byte(reqData.URL))

    shortenedStruct = shortenedURL{urlID, hex.EncodeToString(hash[:3]), reqData.URL}
    urls[urlID] = &shortenedStruct
  } else {
    shortenedStruct = *urls[existingID]
  }

  resData, err := json.Marshal(shortenedStruct)
  if err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(http.StatusOK)
  res.Write(resData)
}

