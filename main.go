package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type TestWay interface {
	Ping() error
	Get(key string) error
	Set(key string) error
}

type RedisClient struct {
	cluster_client  *redis.ClusterClient
	sentinel_client *redis.Client
}

type Architecture string

const (
	Cluster  Architecture = "cluster"
	Sentinel Architecture = "sentinel"
)

func NewRedisClient(address []string, arc Architecture) RedisClient {
	if arc == Cluster {
		return RedisClient{
			cluster_client: redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: address,
			}),
		}
	}
	if arc == Sentinel {
		return RedisClient{
			sentinel_client: redis.NewFailoverClient(&redis.FailoverOptions{
				SentinelAddrs: address,
			}),
		}
	}
	return RedisClient{}
}

func (r *RedisClient) Ping() error {
	fmt.Println("Staring Ping Test")
	if r.cluster_client != nil {
		if _, err := r.cluster_client.Ping(context.TODO()).Result(); err != nil {
			return err
		}
	}
	if r.sentinel_client != nil {
		if _, err := r.sentinel_client.Ping(context.TODO()).Result(); err != nil {
			return err
		}
	}
	fmt.Println("Passed PING Test")
	return nil
}

func (r *RedisClient) Get() error {
	fmt.Println("Staring Get Test")
	if r.cluster_client != nil {
		for i := 0; i < 100; i++ {
			if _, err := r.cluster_client.Get(context.TODO(), strconv.Itoa(i)).Result(); err != nil {
				if err != redis.Nil {
					return err
				}
			}
		}

	}
	if r.sentinel_client != nil {
		for i := 0; i < 100; i++ {
			if _, err := r.sentinel_client.Get(context.TODO(), strconv.Itoa(i)).Result(); err != nil {
				return err
			}
		}

	}
	fmt.Println("Passed Get Test")
	return nil
}

func (r *RedisClient) Set() error {
	fmt.Println("Staring Set Test")
	if r.cluster_client != nil {
		for i := 0; i < 100; i++ {
			if _, err := r.cluster_client.Set(context.TODO(), strconv.Itoa(i), strconv.Itoa(i), 10*time.Second).Result(); err != nil {
				return err
			}
		}

	}
	if r.sentinel_client != nil {
		for i := 0; i < 100; i++ {
			if _, err := r.sentinel_client.Set(context.TODO(), strconv.Itoa(i), strconv.Itoa(i), 10*time.Second).Result(); err != nil {
				return err
			}
		}

	}
	fmt.Println("Passed Set Test")
	return nil
}

func main() {
	var arch = flag.String("arch", "cluster", "support cluster or sentinel")
	var loop = flag.Bool("l", false, "loop test")
	flag.Parse()
	address := flag.Args()

	fmt.Printf("Starting Test arch is: %s , address: %v \n", *arch, address)
	if len(address) == 0 || *arch == "" {
		fmt.Println("need arg:\nexample: \n   redis-test --arch cluster 127.0.0.1:6379")
		return
	}
	client := NewRedisClient(address, Architecture(*arch))
	if err := client.Ping(); err != nil {
		fmt.Println(err.Error())
	}
	if err := client.Get(); err != nil {
		fmt.Println(err.Error())
	}
	if err := client.Set(); err != nil {
		client.Get()
	}
	for *loop {
		if err := client.Ping(); err != nil {
			fmt.Println(err.Error())
		}
		if err := client.Get(); err != nil {
			fmt.Println(err.Error())
		}
		if err := client.Set(); err != nil {
			client.Get()
		}

	}

}
