package redis

import (
	"time"

	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/hhy5861/logrus"
	"github.com/youtube/vitess/go/pools"
	"golang.org/x/net/context"
)

type (
	ResourceConn struct {
		redis.Conn
	}

	redisConnect struct {
		address  string
		port     int64
		password string
		db       int
	}
)

var (
	conn *pools.ResourcePool
)

func (r ResourceConn) Close() {
	r.Conn.Close()
}

func NewRedis(addres string, port int64, db int, password string) {
	config := &redisConnect{
		address:  addres,
		port:     port,
		password: password,
		db:       db,
	}

	config.ConnectRedis()
}

func (r *redisConnect) ConnectRedis() {
	conn = pools.NewResourcePool(func() (pools.Resource, error) {
		addres := fmt.Sprintf("%s:%d", r.address, r.port)
		c, err := redis.Dial("tcp", addres)
		if r.password != "" {
			if _, err := c.Do("AUTH", r.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		if _, err := c.Do("SELECT", r.db); err != nil {
			c.Close()
			return nil, err
		}

		return ResourceConn{c}, err
	}, 1, 20, 3*time.Microsecond)
}

func GetConnect() ResourceConn {
	var resource ResourceConn
	if conn != nil {
		ctx := context.TODO()
		res, err := conn.Get(ctx)
		if err != nil {
			ps := logrus.Params{
				"err": err,
			}
			logrus.Fatal(ps, err, "redis get pool error")
		}
		defer conn.Put(res)

		resource = res.(ResourceConn)
	}

	return resource
}
