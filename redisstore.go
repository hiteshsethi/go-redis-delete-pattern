/*
    Written By: hiteshsethi
    Time      : 23/08/16 4:04 PM 
*/
package main

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"time"
)

type RedisDSN struct {
	Host     string
	Password string
	Port     string
	DbNum    int
}

var masterTravel RedisDSN

type RedisStore struct {
}

func (dCon *RedisStore) init() {
	masterTravel = ConfigRedis
}

var RedisNilErr error

func init()  {
	RedisNilErr = redis.ErrNil
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		MaxActive: 10, // max number of connections
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", masterTravel.Host + ":" + masterTravel.Port)
			if err != nil {
				return nil, err
			}
			if masterTravel.Password != "" {
				_, errAuth := c.Do("AUTH", masterTravel.Password)
				if errAuth != nil {
					panic(errAuth) //worst case
				}
			}
			_, err1 := c.Do("SELECT", masterTravel.DbNum) //here db num is sending
			if err1 != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}

var pool = newPool()

func (dCon *RedisStore) Run(command string, args...interface{}) (reply interface{}, err error) {
	dCon.init()
	c := pool.Get()
	r, err := c.Do(command, args...);
	//defer c.Close()
	c.Close()
	return r, err
}

func (dCon *RedisStore) ConvertToString(value interface{}, err error) (string, error) {
	return redis.String(value, err)
}

func (dCon *RedisStore) ConvertToInt(value interface{}, err error) (int, error) {
	return redis.Int(value, err)
}

func (dCon *RedisStore) ConvertToStringMap(value interface{}, err error) (map[string]string, error) {
	return redis.StringMap(value, err)
}

func (dCon *RedisStore) Values(value interface{}, err error) ([]interface{}, error) {
	return redis.Values(value, err)
}

func (dCon *RedisStore) Strings(value interface{}, err error) ([]string, error) {
	return redis.Strings(value, err)
}

func (dCon *RedisStore) RunTransaction(commands map[string][]string) (reply interface{}, err error) {

	dCon.init()
	c := pool.Get()
	c.Send("MULTI")
	for commandKey, commandVals := range commands {
		commandArgs := make([]interface{}, len(commandVals))
		for i, v := range commandVals {
			commandArgs[i] = v
		}
		c.Send(commandKey, commandArgs...)
	}
	r, err := c.Do("EXEC")
	c.Close()
	return r, err
}

func (dCon *RedisStore) DecodeJson(str string) (interface{}, error) {
	dCon.init()
	var f interface{}
	err := json.Unmarshal([]byte(str), &f)
	return f, err
}

