package main

import (
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

func initCache() *redis.Pool {
	log.Println("Creating Cache pool: redis")

	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ":6379")
			if err != nil {
				log.Printf("ERROR: fail init redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

type CacheAPI struct {
	pool *redis.Pool
}

func NewCacheAPI(pool *redis.Pool) (ca CacheAPI, err error) {
	conn := pool.Get()
	err = ca.ping(conn)
	if err != nil {
		return
	}
	defer conn.Close()
	ca.pool = pool

	return
}

func (ch *CacheAPI) ping(conn redis.Conn) (err error) {
	_, err = redis.String(conn.Do("PING"))

	return
}

func (ch *CacheAPI) Set(key string, val interface{}) error {
	conn := ch.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %v, error %v", key, val, err.Error())
		return err
	}

	return nil
}

func (ch *CacheAPI) Get(key string) (string, error) {
	conn := ch.pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %v", key, err.Error())
		return "", err
	}

	return s, nil
}

func (ch *CacheAPI) Del(key string) (err error) {
	conn := ch.pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	if err != nil {
		log.Printf("ERROR: fail delete key %s, error %v", key, err.Error())
	}

	return
}
