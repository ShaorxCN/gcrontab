package redis

import (
	"gcrontab/custom"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

// GetLock 在redis中注册锁
func GetLock(e, owner string, timeout int) error {
	conn := GetConn()
	if conn.Err() != nil {
		logrus.Errorf("get connection from redisPool failed:%v", conn.Err())
		return custom.ErrorSaveToDBFailed
	}

	defer conn.Close()

	var reply string
	var err error

	if timeout <= 0 {
		reply, err = redis.String(conn.Do("SET", e, owner, "NX"))
	} else {
		reply, err = redis.String(conn.Do("SET", e, owner, "EX", timeout, "NX"))
	}

	if err != nil {
		logrus.Errorf("save data[%s] to redis failed:%v", e, err)
		return custom.ErrorSaveToDBFailed
	}

	if reply != "OK" {
		logrus.Errorf("save object named %s to redis failed,reply is %s", e, reply)
		return custom.ErrorEntityLocked
	}

	return nil

}

// UnLock 在redis中注销entitier
func UnLock(e string) error {
	conn := GetConn()
	if conn.Err() != nil {
		logrus.Errorf("get connection from redisPool failed:%v", conn.Err())
		return custom.ErrorSaveToDBFailed
	}

	defer conn.Close()

	n, err := redis.Int64(conn.Do("DEL", e))
	if err != nil {
		logrus.Errorf("del key[%s] in redis failed:%v", e, err)
		return custom.ErrorSaveToDBFailed
	}

	if n != 1 {
		logrus.Errorf("del %d keys named [%s]", n, e)
	}

	return nil

}
