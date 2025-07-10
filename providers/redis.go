package providers

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type RedisProviderOptions struct {
	Addr     string
	Database int
	Prefix   string
}

type RedisProvider struct {
	pool   *redis.Pool
	db     int
	prefix string
}

func (p *RedisProvider) getConnection() redis.ConnWithTimeout {
	conn := p.pool.Get()
	_, _ = conn.Do("SELECT", p.db)
	return conn.(redis.ConnWithTimeout)
}

func (p *RedisProvider) GetListOfFeatureFlags(ctx context.Context) (map[string]bool, error) {
	conn := p.getConnection()
	defer conn.Close()

	featureFlags := make(map[string]bool)
	cursor := 0
	
	for {
		// Use SCAN to iterate through keys with the prefix
		reply, err := conn.DoWithTimeout(time.Minute, "SCAN", cursor, "MATCH", p.prefix+"*", "COUNT", 100)
		if err != nil {
			return nil, fmt.Errorf("error scanning redis keys: %v", err)
		}

		// Parse the reply: [cursor, [keys...]]
		replyArray, ok := reply.([]interface{})
		if !ok || len(replyArray) != 2 {
			return nil, fmt.Errorf("unexpected reply format from SCAN")
		}

		// Get the new cursor
		cursorBytes, ok := replyArray[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected cursor format")
		}
		cursor, err = redis.Int(cursorBytes, nil)
		if err != nil {
			return nil, fmt.Errorf("error parsing cursor: %v", err)
		}

		// Get the keys
		keys, ok := replyArray[1].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected keys format")
		}

		// Process each key
		for _, keyInterface := range keys {
			keyBytes, ok := keyInterface.([]byte)
			if !ok {
				continue
			}
			key := string(keyBytes)
			
			// Remove the prefix to get the feature flag name
			flagName := key[len(p.prefix):]
			
			// Get the value for this key
			valueReply, err := conn.DoWithTimeout(time.Minute, "GET", key)
			if err != nil {
				log.Printf("error getting value for key %s: %v", key, err)
				continue
			}

			value, err := redis.Bool(valueReply, nil)
			if err != nil {
				log.Printf("error parsing bool value for key %s: %v", key, err)
				continue
			}

			featureFlags[flagName] = value
		}

		// If cursor is 0, we've scanned all keys
		if cursor == 0 {
			break
		}
	}

	return featureFlags, nil
}

func (p *RedisProvider) SetFeatureFlagStatus(ctx context.Context, key string, value bool) error {
	cacheKey := fmt.Sprintf("%s%s", p.prefix, key)
	conn := p.getConnection()
	defer conn.Close()

	_, err := conn.DoWithTimeout(time.Minute, "SET", cacheKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (p *RedisProvider) GetFeatureFlagStatus(ctx context.Context, key string) bool {
	cacheKey := fmt.Sprintf("%s%s", p.prefix, key)
	conn := p.getConnection()
	defer conn.Close()

	reply, err := conn.DoWithTimeout(time.Minute, "GET", cacheKey)
	if err != nil {
		log.Printf("error getting feature flag status for key %s: %v", key, err)
		return false
	}

	val, err := redis.Bool(reply, err)
	if err != nil {
		log.Printf("error getting feature flag status for key %s: %v", key, err)
		return false
	}

	return val

}

func NewRedisProvider(options RedisProviderOptions) IFeatureFlagProvider {
	pool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: time.Second * 10,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(options.Addr, redis.DialConnectTimeout(time.Second*10))
			if err != nil {
				log.Fatal(err)
			}
			return conn, err
		},
	}

	conn, _ := pool.Get().(redis.ConnWithTimeout)
	_, err := conn.DoWithTimeout(time.Minute, "PING")
	if err != nil {
		log.Fatal("failed to connect to redis: " + err.Error())
	}

	return &RedisProvider{
		pool:   pool,
		db:     options.Database,
		prefix: options.Prefix,
	}
}
