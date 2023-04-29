// redis_mem.go
package memory

import (
	"backend/utils"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type RedisMem struct {
	searchclient *redisearch.Client
	redispool    *redis.Pool
}

var memoryindex string = "auto-gpt"
var vecnum int = 0

//var memories []models.Memory = []models.Memory{}

func NewRedisMem() (*RedisMem, error) {

	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Infof(".env file not found, using OS ENV variables. Err: %s", err)
	}

	//get port from env
	redispassword, exists := os.LookupEnv("REDIS_PASSWORD")

	if !exists {
		utils.Logger.Error("REDIS_PASSWORD not defined in env")
	}

	//get the host from env
	redishost, exists := os.LookupEnv("REDIS_HOST")
	if !exists {
		utils.Logger.Info("REDIS_HOST not defined in env, defaulting to localhost")
		redishost = "localhost"
	}

	//get the port from the env
	redisport, exists := os.LookupEnv("REDIS_PORT")
	if !exists {
		utils.Logger.Info("REDIS_PORT not defined in env, defaulting to 6379")
		redisport = "6379"
	}

	//get the memory index from the env
	memoryindex, exists = os.LookupEnv("MEMORY_INDEX")
	if !exists {
		utils.Logger.Info("MEMORY_INDEX not defined in env, defaulting to auto-gpt")
		memoryindex = "auto-gpt"
	}

	//create redis connection pool
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			//return redis.Dial("tcp", redishost+":"+redisport, redis.DialPassword(redispassword))
			c, err := redis.Dial("tcp", redishost+":"+redisport)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", redispassword); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", 0); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	//create redis search client
	client := redisearch.NewClientFromPool(redisPool, memoryindex)

	// Create a schema
	//initialize empty redisearch.VectorFieldOptions struct
	var options redisearch.VectorFieldOptions = redisearch.VectorFieldOptions{
		Algorithm: "HNSW",
		Attributes: map[string]interface{}{
			"TYPE":            "FLOAT32",
			"DIM":             1536,
			"DISTANCE_METRIC": "COSINE",
		},
	}

	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("data")).
		AddField(redisearch.NewVectorFieldOptions("embedding", options))

	// Drop an existing index. If the index does not exist an error is returned
	if err := client.Drop(); err != nil {
		utils.Logger.Error(err)
	}

	def := redisearch.NewIndexDefinition()
	def.AddPrefix(memoryindex + ":")

	// Create the index with the given schema
	if err := client.CreateIndexWithIndexDefinition(sc, def); err != nil {
		utils.Logger.Fatal(err)
	}

	// get the value from Redis
	var vecNumStr *redisearch.Document = nil
	if vecNumStr, err = client.Get(memoryindex + "-vec_num"); err != nil {
		utils.Logger.Errorf("Error getting vecnum from Redis, setting to 0:  %s", err)
		vecnum = 0
	}
	// convert the string value to an int64
	if vecNumStr != nil {
		vecnum, err = strconv.Atoi(string(vecNumStr.Payload))
		if err != nil {
			utils.Logger.Errorf("Error converting vecnum from Redis to int:  %s", err)
			vecnum = 0
		}
	}

	/*
		// Create a document with an id and given score
		doc := redisearch.NewDocument("ExampleNewClientFromPool:doc2", 1.0)
		doc.Set("title", "Hello world").
			Set("body", "foo bar").
			Set("date", time.Now().Unix())

		// Index the document. The API accepts multiple documents at a time
		if err := client.Index([]redisearch.Document{doc}...); err != nil {
			utils.Logger.Fatal(err)
		}

		// Searching with limit and sorting
		docs, total, err := client.Search(redisearch.NewQuery("hello world").
			Limit(0, 2).
			SetReturnFields("title"))

		fmt.Println(docs[0].Id, docs[0].Properties["title"], total, err)
	*/
	///================= TEST IT WORKS =========================///

	return &RedisMem{searchclient: client, redispool: redisPool}, nil
}

func (r *RedisMem) AddMemory(text string) error {

	embedding := createAdaEmbeddings(text)
	if embedding == nil {
		return errors.New("couldn't create embedding")
	}

	// Get a vanilla connection and create 100 hashes
	vanillaConnection := r.redispool.Get()
	vanillaConnection.Do("HSET", fmt.Sprintf("%s:%d", memoryindex, vecnum), "data", text, "embedding", embedding)
	vecnum++
	vanillaConnection.Do("HSET", memoryindex+"-vec_num", vecnum)

	return nil
}

func (r *RedisMem) Clear() error {
	return nil
}

func (r *RedisMem) GetRelevantMemories(query string) []string {

	return []string{}
}

func (r *RedisMem) GetStats() int {

	return 0
}
