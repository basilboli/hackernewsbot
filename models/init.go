package models

import (
	"log"

	"github.com/mediocregopher/radix.v2/pool"
)

var db *pool.Pool

func init() {
	var err error
	db, err = pool.New("tcp", "redis:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}
