package services

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/fmo/players-api/internal/database"
	"github.com/fmo/players-api/internal/models"
)

const tableName = "fmo-players"

type PlayersService struct {
	DB *database.Database
}

func NewPlayers(db *database.Database) PlayersService {
	return PlayersService{
		DB: db,
	}
}

func (ps PlayersService) CreateOrUpdate(p models.Player) (response *dynamodb.PutItemOutput, err error) {
	playerParsed, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      playerParsed,
		TableName: aws.String(tableName),
	}

	return ps.DB.Connection.PutItem(input)
}

func (ps PlayersService) FindPlayers(teamId int) (players []models.Player, err error) {
	filter := expression.Name("teamId").Equal(expression.Value(teamId))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return players, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	result, err := ps.DB.Connection.Scan(input)
	if err != nil {
		return players, err
	}

	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &players)
		if err != nil {
			return players, err
		} else {
			return players, nil
		}
	}

	return nil, errors.New("no result")
}

func (ps PlayersService) FindOnePlayer(playerId string) (player *models.Player, err error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(playerId),
			},
		},
	}

	result, err := ps.DB.Connection.GetItem(input)
	if err != nil {
		return player, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("no player found with playerId: %s", playerId)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &player)
	if err != nil {
		return player, err
	}

	return player, nil
}
