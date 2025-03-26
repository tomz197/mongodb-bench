package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Query represents a MongoDB query to benchmark
type Query struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Query       interface{} `json:"query"`
	Collection  string      `json:"collection"`
}

// BenchmarkResult stores the results of a single query benchmark
type BenchmarkResult struct {
	Name        string
	Description string
	Collection  string
	Runs        int
	TotalTime   time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
	AvgTime     time.Duration
}

// Benchmark represents the main benchmark runner
type Benchmark struct {
	Client     *mongo.Client
	Database   string
	Iterations int
}

// NewBenchmark creates a new benchmark runner
func NewBenchmark(connectionURI, database string, iterations int) (*Benchmark, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return &Benchmark{
		Client:     client,
		Database:   database,
		Iterations: iterations,
	}, nil
}

// Close closes the MongoDB connection
func (b *Benchmark) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return b.Client.Disconnect(ctx)
}

// LoadQueries loads benchmark queries from a JSON file
func LoadQueries(filename string) ([]Query, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var queries []Query
	err = json.Unmarshal(data, &queries)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return queries, nil
}

// RunBenchmark runs the benchmark for a single query
func (b *Benchmark) RunBenchmark(query Query) BenchmarkResult {
	collection := b.Client.Database(b.Database).Collection(query.Collection)

	result := BenchmarkResult{
		Name:        query.Name,
		Description: query.Description,
		Collection:  query.Collection,
		Runs:        b.Iterations,
		MinTime:     time.Duration(1<<63 - 1), // Max duration value
	}

	var totalTime time.Duration

	for i := 0; i < b.Iterations; i++ {
		ctx := context.Background()
		start := time.Now()

		// Check if the query is an aggregation pipeline (array) or a find query (object)
		isAggregation := false
		if arr, ok := query.Query.([]interface{}); ok && len(arr) > 0 {
			isAggregation = true
		}

		var cursor *mongo.Cursor
		var err error

		if isAggregation {
			// Handle aggregation pipeline
			pipeline, err := convertToPipeline(query.Query)
			if err != nil {
				log.Printf("Error converting aggregation pipeline : %v", err)
				continue
			}

			cursor, err = collection.Aggregate(ctx, pipeline)
		} else {
			// Handle find query
			queryBSON, err := convertToBSON(query.Query)
			if err != nil {
				log.Printf("Error converting query to BSON: %v", err)
				continue
			}

			cursor, err = collection.Find(ctx, queryBSON)
		}

		if err != nil {
			log.Printf("Error executing query: %v", err)
			continue
		}

		// Consume the cursor to ensure the query completes
		for cursor.Next(ctx) {
			// Just iterate without doing anything
		}
		cursor.Close(ctx)

		elapsed := time.Since(start)
		totalTime += elapsed

		if elapsed < result.MinTime {
			result.MinTime = elapsed
		}
		if elapsed > result.MaxTime {
			result.MaxTime = elapsed
		}
	}

	result.TotalTime = totalTime
	result.AvgTime = totalTime / time.Duration(b.Iterations)

	return result
}

// Helper function to convert interface{} to BSON document
func convertToBSON(query interface{}) (bson.D, error) {
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	var bsonDoc bson.D
	err = bson.UnmarshalExtJSON(data, true, &bsonDoc)
	if err != nil {
		return nil, err
	}

	return bsonDoc, nil
}

// Helper function to convert interface{} to aggregation pipeline
func convertToPipeline(query interface{}) (mongo.Pipeline, error) {
	// First serialize the entire query to JSON
	data, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pipeline: %v", err)
	}

	// Then directly unmarshal into a mongo.Pipeline
	var pipeline mongo.Pipeline
	err = bson.UnmarshalExtJSON(data, true, &pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal pipeline: %v", err)
	}

	return pipeline, nil
}

// RunBenchmarks runs benchmarks for multiple queries
func (b *Benchmark) RunBenchmarks(queries []Query) []BenchmarkResult {
	results := make([]BenchmarkResult, 0, len(queries))

	for _, query := range queries {
		log.Printf("Running query: %s", query.Name)
		result := b.RunBenchmark(query)
		results = append(results, result)
	}

	return results
}

// PrintResults prints benchmark results to the console
func PrintResults(results []BenchmarkResult) {
	fmt.Println("\nMongoDB Benchmark Results")
	fmt.Println("=======================")

	for _, result := range results {
		fmt.Printf("\n- Query: %s\n", result.Name)
		fmt.Printf("  Description: %s\n", result.Description)
		fmt.Printf("  Collection: %s\n", result.Collection)
		fmt.Printf("  Iterations: %d\n", result.Runs)
		fmt.Printf("  Total Time: %v\n", result.TotalTime)
		fmt.Printf("  Average Time: %v\n", result.AvgTime)
		fmt.Printf("  Min Time: %v\n", result.MinTime)
		fmt.Printf("  Max Time: %v\n", result.MaxTime)
	}
}
