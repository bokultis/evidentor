package apiutil

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"sync"

// 	"golang.org/x/net/context"

// 	"github.com/go-redis/redis"

// 	"git.horisen.biz/bg/horisenlib/common/util/etcdconf"
// )

// var (
// 	// ErrNoRedisConnection is returned on GetDB when no database is connected so far
// 	ErrNoRedisConnection = errors.New("No redis connection")
// )

// // RedisDynConf keeps all the data for dynamic configuration
// type RedisDynConf struct {
// 	// Key used for dynamic configuration from etcd
// 	key string

// 	// Options keeps redis connection parameters
// 	options         *redis.Options
// 	failoverOptions *redis.FailoverOptions

// 	// Client is current connection to redis client
// 	client *redis.Client

// 	// etcd client
// 	etcdClient client.Client

// 	// Keys api
// 	kapi client.KeysAPI

// 	// Lock
// 	lock sync.Mutex

// 	// isFailover
// 	isFailover bool
// }

// func NewRedisDynConf(key string) *RedisDynConf {
// 	return &RedisDynConf{
// 		key:        key,
// 		isFailover: false,
// 	}
// }

// func NewRedisFailoverConf(key string) *RedisDynConf {
// 	return &RedisDynConf{
// 		key:        key,
// 		isFailover: true,
// 	}
// }

// // Initialize initializes mysql connector and connection
// func (c *RedisDynConf) Initialize() error {
// 	kapi, etcdClient, err := etcdconf.CreateEtcdClient()
// 	if err != nil {
// 		return err
// 	}

// 	c.kapi = kapi
// 	c.etcdClient = etcdClient

// 	err = c.reloadConnection(nil)
// 	if err != nil {
// 		return err
// 	}

// 	go c.watchChanges()

// 	return nil
// }

// func (c *RedisDynConf) GetClient() (*redis.Client, error) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()

// 	if c.client == nil {
// 		return nil, ErrNoRedisConnection
// 	}
// 	return c.client, nil
// }

// func (c *RedisDynConf) watchChanges() {
// 	kapi := client.NewKeysAPI(c.etcdClient)

// 	watcher := kapi.Watcher(c.key,
// 		&client.WatcherOptions{AfterIndex: 0, Recursive: false})
// 	for {
// 		log.Printf("Get next change...\n")
// 		rsp, err := c.getNextChange(watcher)
// 		if err != nil {
// 			log.Printf("Error watching key %s\n", err)
// 			continue
// 		}
// 		if rsp != nil {
// 			log.Printf("Got key for change\n")
// 			c.reloadConnection(rsp)
// 		}
// 	}
// }

// func (c *RedisDynConf) getNextChange(watcher client.Watcher) (*client.Response, error) {
// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()

// 	ctx := context.Background()
// 	return watcher.Next(ctx)
// }

// func (c *RedisDynConf) getConf(rsp *client.Response) (*redis.Options, error) {
// 	if rsp == nil {
// 		resp, err := c.kapi.Get(context.Background(), c.key, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		rsp = resp
// 	}
// 	if rsp.Node != nil {
// 		var options redis.Options
// 		b := []byte(rsp.Node.Value)
// 		err := json.Unmarshal(b, &options)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &options, nil
// 	}
// 	return nil, fmt.Errorf("Key %s is empty", c.key)
// }

// func (c *RedisDynConf) getFailoverConf(rsp *client.Response) (*redis.FailoverOptions, error) {
// 	if rsp == nil {
// 		resp, err := c.kapi.Get(context.Background(), c.key, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		rsp = resp
// 	}
// 	if rsp.Node != nil {
// 		var failoverOptions redis.FailoverOptions
// 		b := []byte(rsp.Node.Value)
// 		err := json.Unmarshal(b, &failoverOptions)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &failoverOptions, nil
// 	}
// 	return nil, fmt.Errorf("Key %s is empty", c.key)
// }

// func (c *RedisDynConf) reloadConnection(rsp *client.Response) error {
// 	if c.isFailover {
// 		failoverOptions, err := c.getFailoverConf(rsp)
// 		if err != nil {
// 			return err
// 		}

// 		log.Printf("Reloading connection...\n")

// 		client := redis.NewFailoverClient(failoverOptions)

// 		// Execute Ping to test connection
// 		pong, err := client.Ping().Result()
// 		if err != nil {
// 			return err
// 		}
// 		log.Printf("PONG=%#v\n", pong)
// 		c.replaceFailoverConnection(failoverOptions, client)
// 	} else {
// 		options, err := c.getConf(rsp)
// 		if err != nil {
// 			return err
// 		}

// 		log.Printf("Reloading connection...\n")

// 		client := redis.NewClient(options)

// 		// Execute Ping to test connection
// 		pong, err := client.Ping().Result()
// 		if err != nil {
// 			return err
// 		}
// 		log.Printf("PONG=%#v\n", pong)
// 		c.replaceConnection(options, client)
// 	}
// 	return nil
// }

// func (c *RedisDynConf) replaceConnection(options *redis.Options, client *redis.Client) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()

// 	c.options = options
// 	c.client = client
// }

// func (c *RedisDynConf) replaceFailoverConnection(failoverOptions *redis.FailoverOptions, client *redis.Client) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()

// 	c.failoverOptions = failoverOptions
// 	c.client = client
// }
