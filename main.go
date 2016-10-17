package main

import (
  "github.com/kristofmic/go_url/sources/redis"
  "github.com/kristofmic/go_url/web"
)

func main() {
  redis.Init()
  web.NewWeb()
}