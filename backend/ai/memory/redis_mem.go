// redis_mem.go
package memory

import (
	"backend/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	//TODO: consider https://github.com/rueian/rueidis as replacement for redis
	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type RedisMem struct {
	searchclient *redisearch.Client
	redispool    *redis.Pool
}

var memoryindex string = "auto-gpt"

//var vecnum int = 0

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
	//if err := client.Drop(); err != nil {
	//	utils.Logger.Error(err)
	//}

	def := redisearch.NewIndexDefinition()
	def.AddPrefix(memoryindex + ":")

	// Create the index with the given schema
	if err := client.CreateIndexWithIndexDefinition(sc, def); err != nil {
		utils.Logger.Info(err)
	}

	// get the value from Redis
	info, err := client.Info()
	if err != nil {
		utils.Logger.Errorf("Error getting info from Redis, probably index not created:  %s", err)
	}

	//dump usefull info to log
	utils.Logger.Infof("redis Index Name: %s", info.Name)
	utils.Logger.Infof("redis Index DocCount: %d", info.DocCount)
	utils.Logger.Infof("redis Index IsIndexing: %t", info.IsIndexing)
	utils.Logger.Infof("redis Index RecordCount: %d", info.RecordCount)
	utils.Logger.Infof("redis Index MaxDocID: %d", info.MaxDocID)
	utils.Logger.Infof("redis Index PercentIndexed: %f", info.PercentIndexed)
	utils.Logger.Infof("redis Index HashIndexingFailures: %d", info.HashIndexingFailures)
	//utils.Logger.Infof("redis Index info: %v", info)

	/*
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
	*/
	return &RedisMem{searchclient: client, redispool: redisPool}, nil
}

func (r *RedisMem) AddMemory(text string) error {

	embedding := createAdaEmbeddings(text)
	if embedding == nil {
		return errors.New("couldn't create embedding")
	}

	// get the value from Redis
	info, err := r.searchclient.Info()
	if err != nil {
		utils.Logger.Errorf("Error getting info from Redis, probably index not created:  %s", err)
	}

	//convert embeddings to byte array
	bytes := utils.Float32ToBytesFastSafe(embedding)

	// Create a document with an id and given score
	doc := redisearch.NewDocument(fmt.Sprintf("%s:%d", memoryindex, info.MaxDocID), 1.0)
	doc.Set("data", text).
		Set("embedding", bytes)

	//set index options
	var opt redisearch.IndexingOptions = redisearch.DefaultIndexingOptions
	opt.Replace = true
	// Index the document. The API accepts multiple documents at a time
	if err := r.searchclient.IndexOptions(opt, []redisearch.Document{doc}...); err != nil {
		return err
	}

	return nil
}

func (r *RedisMem) Clear() error {
	err := r.searchclient.Drop()
	if err != nil {
		utils.Logger.Errorf("Error dropping redis index: %s", err)
		return err
	}
	return nil
}

func (r *RedisMem) GetRelevantMemories(data string, max int) []string {

	//nothing to do
	if strings.TrimSpace(data) == "" {
		return []string{}
	}

	embedding := createAdaEmbeddings(data)
	if len(embedding) == 0 {
		return []string{}
	}

	//convert embeddings to byte array
	bytes := utils.Float32ToBytesFastSafe(embedding)
	param := map[string]interface{}{"vector": bytes}

	//build the query
	queryStr := fmt.Sprintf("*=>[KNN %d @embedding $vector AS vector_score]", max)
	query := redisearch.NewQuery(queryStr).
		SetReturnFields("data", "vector_score").
		SetDialect(2).
		SetSortBy("vector_score", false).
		SetParams(param)

	//search the index
	qResults, total, err := r.searchclient.Search(query)
	if err != nil {
		utils.Logger.Errorf("Error querying index: %s", err)
		return []string{}
	} else if total == 0 {
		return []string{}
	}

	//get the results
	var results []string = []string{}
	for i := 0; i < len(qResults); i++ {
		if data, ok := qResults[i].Properties["data"].(string); ok {
			results = append(results, data)
		} else {
			// handle the case where the value is not a string
		}
	}
	return results
}

func (r *RedisMem) GetStats() interface{} {

	info, err := r.searchclient.Info()
	if err != nil {
		utils.Logger.Errorf("Error getting index info from RedisMem: %s", err)
		return nil
	}
	return info
}
