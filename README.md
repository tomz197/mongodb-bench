# MongoDB Benchmark Tool

A simple benchmarking tool for MongoDB queries written in Go.

## Features

- Benchmark multiple MongoDB queries from a JSON file
- Configurable number of iterations for each query
- Detailed statistics for each query (min, max, avg time)
- Support for different MongoDB connection URIs and database names

## Installation

```bash
# Clone the repository
git clone https://github.com/tomz197/mongodb-bench.git
cd mongodb-bench

# Build the tool
go build -o mongodb-bench ./cmd/mongodb-bench
```

## Usage

```bash
./mongodb-bench -queries example-queries.json -uri mongodb://localhost:27017 -db test -iterations 100
```

### Command Line Arguments

- `-uri`: MongoDB connection URI (default: "mongodb://localhost:27017")
- `-db`: MongoDB database name (default: "test")
- `-iterations`: Number of times to run each query (default: 10)
- `-queries`: Path to JSON file containing queries (required)

## Query JSON Format

The queries JSON file should contain an array of objects with the following structure:

```json
[
  {
    "name": "Query Name",
    "description": "Query Description",
    "collection": "collection_name",
    "query": {
      // MongoDB query object or aggregation pipeline
    }
  }
]
```

See `example-queries.json` for a complete example.

## Example Output

```
Loaded 4 queries from example-queries.json
Connected to MongoDB at mongodb://localhost:27017
Using database: test
Running each query 10 times
Starting benchmark...

MongoDB Benchmark Results
=======================

- Query: Find All Documents
  Description: Simple query to retrieve all documents in a collection
  Collection: users
  Iterations: 10
  Total Time: 14.7ms
  Average Time: 1.47ms
  Min Time: 1.2ms
  Max Time: 2.5ms
  P50 (Median): 1.4ms
  P95: 2.1ms
  P99: 2.4ms
  Std Dev: 0.31ms
  Successful: 10
  Errors: 0

- Query: Find Documents by Age
  Description: Query to find users above age 30
  Collection: users
  Iterations: 10
  Total Time: 12.1ms
  Average Time: 1.21ms
  Min Time: 1.0ms
  Max Time: 1.8ms
  P50 (Median): 1.15ms
  P95: 1.6ms
  P99: 1.75ms
  Std Dev: 0.25ms
  Successful: 10
  Errors: 0
```

## License

MIT 