package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/fmo/players-api/internal/application/core/domain"
)

type Adapter struct {
	Connection *dynamodb.DynamoDB
	TableName  string
}

func NewAdapter(tableName string) (*Adapter, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Adapter{
		Connection: dynamodb.New(sess),
		TableName:  tableName,
	}, nil
}

func (a Adapter) FindPlayersByTeamId(ctx context.Context, teamId int) (players []domain.Player, err error) {
	filter := expression.Name("teamId").Equal(expression.Value(teamId))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return players, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(a.TableName),
	}

	result, err := a.Connection.Scan(input)
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

func (a Adapter) FindPlayersById(ctx context.Context, playerId string) (player domain.Player, err error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(a.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(playerId),
			},
		},
	}

	result, err := a.Connection.GetItem(input)
	if err != nil {
		return player, err
	}

	if result.Item == nil {
		return domain.Player{}, fmt.Errorf("no player found with playerId: %s", playerId)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &player)
	if err != nil {
		return player, err
	}

	return player, nil
}
