package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mediocregopher/radix/v4"
)

var errKeyCollision = errors.New("key collision")
var errKeyNotFound = fmt.Errorf("key %w", errNotFound)

type store interface {
	put(key string, data []byte, expire time.Duration) error
	get(key string) ([]byte, error)
}

type redisStore struct {
	client radixDoer
}

type radixDoer interface {
	Do(context.Context, radix.Action) error
	Close() error
}

func newRedisStore(connection string) (*redisStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	connectionParts := strings.SplitN(connection, ":", 2)
	rtype := connectionParts[0]

	var client radixDoer
	var err error
	if rtype == "simple" {
		client, err = (radix.PoolConfig{}).New(ctx, "tcp", connectionParts[1])
		if err != nil {
			return nil, nil
		}
	} else if rtype == "cluster" {
		clusterAddrs := strings.Split(connectionParts[1], ",")
		client, err = (radix.ClusterConfig{}).New(ctx, clusterAddrs)
	}
	if err != nil {
		return nil, nil
	}
	return &redisStore{client}, nil
}

func (s *redisStore) put(key string, data []byte, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	exists := 0
	err := s.client.Do(ctx, radix.Cmd(&exists, "EXISTS", key))
	if err != nil {
		return err
	}

	if exists != 0 {
		return fmt.Errorf("%w: %s", errKeyCollision, key)
	}

	expireS := int64(expire / time.Second)
	err = s.client.Do(ctx, radix.FlatCmd(nil, "SETEX", key, expireS, data))
	if err != nil {
		return err
	}
	return nil
}

func (s *redisStore) get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var data []byte
	err := s.client.Do(ctx, radix.Cmd(&data, "GET", key))
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("%w: %s", errKeyNotFound, key)
	}
	return data, nil
}
