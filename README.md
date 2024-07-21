# Players API and Consumer

This is a Go-based API for managing football player data. It stores the player profile image on AWS S3 for storage.
Uses Kafka for messaging, getting the player data and writing it to the AWS DynamoDB 

## Getting Started

Just get the binary for either consumer or the api or for the both and basically run it.

### Prerequisites

- AWS Client environment for DynamoDB
- Kafka credentials

### Set The Env Vars

Check .env.dist template and set the variables accordingly

### How To Run

```
./playersApi
```

```
./playerConsumer
```

### Generate Open API 

oapi-codegen --config=config.yaml ./api.yaml
