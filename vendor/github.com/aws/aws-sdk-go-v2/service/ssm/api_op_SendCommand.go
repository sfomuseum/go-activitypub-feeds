// Code generated by smithy-go-codegen DO NOT EDIT.

package ssm

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Runs commands on one or more managed nodes.
func (c *Client) SendCommand(ctx context.Context, params *SendCommandInput, optFns ...func(*Options)) (*SendCommandOutput, error) {
	if params == nil {
		params = &SendCommandInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "SendCommand", params, optFns, c.addOperationSendCommandMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*SendCommandOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type SendCommandInput struct {

	// The name of the Amazon Web Services Systems Manager document (SSM document) to
	// run. This can be a public document or a custom document. To run a shared
	// document belonging to another account, specify the document Amazon Resource Name
	// (ARN). For more information about how to use shared documents, see Using shared
	// SSM documents (https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-using-shared.html)
	// in the Amazon Web Services Systems Manager User Guide. If you specify a document
	// name or ARN that hasn't been shared with your account, you receive an
	// InvalidDocument error.
	//
	// This member is required.
	DocumentName *string

	// The CloudWatch alarm you want to apply to your command.
	AlarmConfiguration *types.AlarmConfiguration

	// Enables Amazon Web Services Systems Manager to send Run Command output to
	// Amazon CloudWatch Logs. Run Command is a capability of Amazon Web Services
	// Systems Manager.
	CloudWatchOutputConfig *types.CloudWatchOutputConfig

	// User-specified information about the command, such as a brief description of
	// what the command should do.
	Comment *string

	// The Sha256 or Sha1 hash created by the system when the document was created.
	// Sha1 hashes have been deprecated.
	DocumentHash *string

	// Sha256 or Sha1. Sha1 hashes have been deprecated.
	DocumentHashType types.DocumentHashType

	// The SSM document version to use in the request. You can specify $DEFAULT,
	// $LATEST, or a specific version number. If you run commands by using the Command
	// Line Interface (Amazon Web Services CLI), then you must escape the first two
	// options by using a backslash. If you specify a version number, then you don't
	// need to use the backslash. For example: --document-version "\$DEFAULT"
	// --document-version "\$LATEST" --document-version "3"
	DocumentVersion *string

	// The IDs of the managed nodes where the command should run. Specifying managed
	// node IDs is most useful when you are targeting a limited number of managed
	// nodes, though you can specify up to 50 IDs. To target a larger number of managed
	// nodes, or if you prefer not to list individual node IDs, we recommend using the
	// Targets option instead. Using Targets , which accepts tag key-value pairs to
	// identify the managed nodes to send commands to, you can a send command to tens,
	// hundreds, or thousands of nodes at once. For more information about how to use
	// targets, see Using targets and rate controls to send commands to a fleet (https://docs.aws.amazon.com/systems-manager/latest/userguide/send-commands-multiple.html)
	// in the Amazon Web Services Systems Manager User Guide.
	InstanceIds []string

	// (Optional) The maximum number of managed nodes that are allowed to run the
	// command at the same time. You can specify a number such as 10 or a percentage
	// such as 10%. The default value is 50 . For more information about how to use
	// MaxConcurrency , see Using concurrency controls (https://docs.aws.amazon.com/systems-manager/latest/userguide/send-commands-multiple.html#send-commands-velocity)
	// in the Amazon Web Services Systems Manager User Guide.
	MaxConcurrency *string

	// The maximum number of errors allowed without the command failing. When the
	// command fails one more time beyond the value of MaxErrors , the systems stops
	// sending the command to additional targets. You can specify a number like 10 or a
	// percentage like 10%. The default value is 0 . For more information about how to
	// use MaxErrors , see Using error controls (https://docs.aws.amazon.com/systems-manager/latest/userguide/send-commands-multiple.html#send-commands-maxerrors)
	// in the Amazon Web Services Systems Manager User Guide.
	MaxErrors *string

	// Configurations for sending notifications.
	NotificationConfig *types.NotificationConfig

	// The name of the S3 bucket where command execution responses should be stored.
	OutputS3BucketName *string

	// The directory structure within the S3 bucket where the responses should be
	// stored.
	OutputS3KeyPrefix *string

	// (Deprecated) You can no longer specify this parameter. The system ignores it.
	// Instead, Systems Manager automatically determines the Amazon Web Services Region
	// of the S3 bucket.
	OutputS3Region *string

	// The required and optional parameters specified in the document being run.
	Parameters map[string][]string

	// The ARN of the Identity and Access Management (IAM) service role to use to
	// publish Amazon Simple Notification Service (Amazon SNS) notifications for Run
	// Command commands. This role must provide the sns:Publish permission for your
	// notification topic. For information about creating and using this service role,
	// see Monitoring Systems Manager status changes using Amazon SNS notifications (https://docs.aws.amazon.com/systems-manager/latest/userguide/monitoring-sns-notifications.html)
	// in the Amazon Web Services Systems Manager User Guide.
	ServiceRoleArn *string

	// An array of search criteria that targets managed nodes using a Key,Value
	// combination that you specify. Specifying targets is most useful when you want to
	// send a command to a large number of managed nodes at once. Using Targets , which
	// accepts tag key-value pairs to identify managed nodes, you can send a command to
	// tens, hundreds, or thousands of nodes at once. To send a command to a smaller
	// number of managed nodes, you can use the InstanceIds option instead. For more
	// information about how to use targets, see Sending commands to a fleet (https://docs.aws.amazon.com/systems-manager/latest/userguide/send-commands-multiple.html)
	// in the Amazon Web Services Systems Manager User Guide.
	Targets []types.Target

	// If this time is reached and the command hasn't already started running, it
	// won't run.
	TimeoutSeconds *int32

	noSmithyDocumentSerde
}

type SendCommandOutput struct {

	// The request as it was received by Systems Manager. Also provides the command ID
	// which can be used future references to this request.
	Command *types.Command

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationSendCommandMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpSendCommand{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpSendCommand{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "SendCommand"); err != nil {
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
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpSendCommandValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opSendCommand(options.Region), middleware.Before); err != nil {
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
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opSendCommand(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "SendCommand",
	}
}
