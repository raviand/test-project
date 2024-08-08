package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/raviand/test-project/pkg"
)

type RepositoryDynamo interface {
	Store(ctx context.Context, model *pkg.User) error
	GetOne(ctx context.Context, id string) (*pkg.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, model *pkg.User) error
}

func itemToUser(input map[string]*dynamodb.AttributeValue) (*pkg.User, error) {
	var item pkg.User
	err := dynamodbattribute.UnmarshalMap(input, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

type dynamoRepository struct {
	dynamo *dynamodb.DynamoDB
	table  string
}

func NewDynamoRepository(dynamo *dynamodb.DynamoDB, table string) RepositoryDynamo {
	return &dynamoRepository{
		dynamo: dynamo,
		table:  table,
	}
}

func (receiver *dynamoRepository) Store(ctx context.Context, model *pkg.User) error {
	av, err := dynamodbattribute.MarshalMap(model)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(receiver.table),
	}

	_, err = receiver.dynamo.PutItemWithContext(ctx, input)

	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}

	return nil
}

func (receiver *dynamoRepository) GetOne(ctx context.Context, id string) (*pkg.User, error) {
	result, err := receiver.dynamo.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(receiver.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, pkg.GetError(pkg.InternalError)
	}

	if result.Item == nil {
		return nil, nil
	}
	return itemToUser(result.Item)
}

func (receiver *dynamoRepository) Delete(ctx context.Context, id string) error {
	out, err := receiver.dynamo.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(receiver.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	if out == nil {
		return pkg.GetError(pkg.InternalError)
	}
	if out.Attributes == nil {
		return pkg.GetError(pkg.NotFound)
	}
	return err
}

func (receiver *dynamoRepository) Update(ctx context.Context, model *pkg.User) error {
	av, err := dynamodbattribute.MarshalMap(model)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(receiver.table),
	}

	_, err = receiver.dynamo.PutItemWithContext(ctx, input)

	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}

	return nil
}
