package redis

import (
  "time"
  "gopkg.in/redis.v5"
)

var client *redis.Client

func Init() *redis.Client {
  client = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
    MaxRetries: 3,
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 5 * time.Second,
  })

  return client
}

func GetClient() *redis.Client {
  if client == nil {
    Init()
  }

  return client
}
