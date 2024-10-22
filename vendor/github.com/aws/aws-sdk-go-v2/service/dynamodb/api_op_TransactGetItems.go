// Code generated by smithy-go-codegen DO NOT EDIT.

package dynamodb

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	internalEndpointDiscovery "github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// TransactGetItems is a synchronous operation that atomically retrieves multiple
// items from one or more tables (but not from indexes) in a single account and
// Region. A TransactGetItems call can contain up to 100 TransactGetItem objects,
// each of which contains a Get structure that specifies an item to retrieve from
// a table in the account and Region. A call to TransactGetItems cannot retrieve
// items from tables in more than one Amazon Web Services account or Region. The
// aggregate size of the items in the transaction cannot exceed 4 MB.
//
// DynamoDB rejects the entire TransactGetItems request if any of the following is
// true:
//
//   - A conflicting operation is in the process of updating an item to be read.
//
//   - There is insufficient provisioned capacity for the transaction to be
//     completed.
//
//   - There is a user error, such as an invalid data format.
//
//   - The aggregate size of the items in the transaction exceeded 4 MB.
func (c *Client) TransactGetItems(ctx context.Context, params *TransactGetItemsInput, optFns ...func(*Options)) (*TransactGetItemsOutput, error) {
	if params == nil {
		params = &TransactGetItemsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "TransactGetItems", params, optFns, c.addOperationTransactGetItemsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*TransactGetItemsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type TransactGetItemsInput struct {

	// An ordered array of up to 100 TransactGetItem objects, each of which contains a
	// Get structure.
	//
	// This member is required.
	TransactItems []types.TransactGetItem

	// A value of TOTAL causes consumed capacity information to be returned, and a
	// value of NONE prevents that information from being returned. No other value is
	// valid.
	ReturnConsumedCapacity types.ReturnConsumedCapacity

	noSmithyDocumentSerde
}

type TransactGetItemsOutput struct {

	// If the ReturnConsumedCapacity value was TOTAL , this is an array of
	// ConsumedCapacity objects, one for each table addressed by TransactGetItem
	// objects in the TransactItems parameter. These ConsumedCapacity objects report
	// the read-capacity units consumed by the TransactGetItems call in that table.
	ConsumedCapacity []types.ConsumedCapacity

	// An ordered array of up to 100 ItemResponse objects, each of which corresponds
	// to the TransactGetItem object in the same position in the TransactItems array.
	// Each ItemResponse object contains a Map of the name-value pairs that are the
	// projected attributes of the requested item.
	//
	// If a requested item could not be retrieved, the corresponding ItemResponse
	// object is Null, or if the requested item has no projected attributes, the
	// corresponding ItemResponse object is an empty Map.
	Responses []types.ItemResponse

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationTransactGetItemsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson10_serializeOpTransactGetItems{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson10_deserializeOpTransactGetItems{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "TransactGetItems"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpTransactGetItemsDiscoverEndpointMiddleware(stack, options, c); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpTransactGetItemsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opTransactGetItems(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addValidateResponseChecksum(stack, options); err != nil {
		return err
	}
	if err = addAcceptEncodingGzip(stack, options); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func addOpTransactGetItemsDiscoverEndpointMiddleware(stack *middleware.Stack, o Options, c *Client) error {
	return stack.Finalize.Insert(&internalEndpointDiscovery.DiscoverEndpoint{
		Options: []func(*internalEndpointDiscovery.DiscoverEndpointOptions){
			func(opt *internalEndpointDiscovery.DiscoverEndpointOptions) {
				opt.DisableHTTPS = o.EndpointOptions.DisableHTTPS
				opt.Logger = o.Logger
			},
		},
		DiscoverOperation:            c.fetchOpTransactGetItemsDiscoverEndpoint,
		EndpointDiscoveryEnableState: o.EndpointDiscovery.EnableEndpointDiscovery,
		EndpointDiscoveryRequired:    false,
		Region:                       o.Region,
	}, "ResolveEndpointV2", middleware.After)
}

func (c *Client) fetchOpTransactGetItemsDiscoverEndpoint(ctx context.Context, region string, optFns ...func(*internalEndpointDiscovery.DiscoverEndpointOptions)) (internalEndpointDiscovery.WeightedAddress, error) {
	input := getOperationInput(ctx)
	in, ok := input.(*TransactGetItemsInput)
	if !ok {
		return internalEndpointDiscovery.WeightedAddress{}, fmt.Errorf("unknown input type %T", input)
	}
	_ = in

	identifierMap := make(map[string]string, 0)
	identifierMap["sdk#Region"] = region

	key := fmt.Sprintf("DynamoDB.%v", identifierMap)

	if v, ok := c.endpointCache.Get(key); ok {
		return v, nil
	}

	discoveryOperationInput := &DescribeEndpointsInput{}

	opt := internalEndpointDiscovery.DiscoverEndpointOptions{}
	for _, fn := range optFns {
		fn(&opt)
	}

	go c.handleEndpointDiscoveryFromService(ctx, discoveryOperationInput, region, key, opt)
	return internalEndpointDiscovery.WeightedAddress{}, nil
}

func newServiceMetadataMiddleware_opTransactGetItems(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "TransactGetItems",
	}
}