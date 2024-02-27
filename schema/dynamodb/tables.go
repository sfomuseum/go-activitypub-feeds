package dynamodb

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/SecondaryIndexes.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GSI.html

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var FEEDS_TABLE_NAME = "feeds_publication_logs"

var BILLING_MODE = aws.String("PAY_PER_REQUEST")

var DynamoDBTables = map[string]*dynamodb.CreateTableInput{
	FEEDS_TABLE_NAME: DynamoDBFeedsPublicationsLogsTable,
}
