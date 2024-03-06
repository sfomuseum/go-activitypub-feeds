package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBFeedsPublicationsLogsTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Id"), // partition key
			KeyType:       aws.String("HASH"),
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("FeedURL"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("Published"),
			AttributeType: aws.String("N"),
		},		
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("account_feed"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("FeedURL"),
					KeyType:       aws.String("RANGE"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("KEYS_ONLY"),
			},
		},
		{
			IndexName: aws.String("by_feed"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("FeedURL"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Published"),
					KeyType:       aws.String("RANGE"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("KEYS_ONLY"),
			},
		},		
		{
			IndexName: aws.String("by_published"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Published"),
					KeyType:       aws.String("HASH"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("KEYS_ONLY"),
			},
		},		
	},
	BillingMode: BILLING_MODE,
	TableName:   &FEEDS_TABLE_NAME,
}
