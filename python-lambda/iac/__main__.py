import json

import pulumi
import pulumi_aws as aws
from pulumi import ResourceOptions
from pulumi_aws import iam, lambda_, s3, ssm

slack_api_token = pulumi.Config().require_secret("slack_api_token")

# Create an S3 bucket for the input and output
input_bucket = s3.Bucket("rasc-input-bucket")
output_bucket = s3.Bucket("rasc-output-bucket")

# Store the Slack token in the AWS Systems Manager Parameter Store
slack_api_token_parameter = ssm.Parameter(
    "slack_api_token", type="SecureString", value=slack_api_token
)

# Create an IAM role for the Lambda function
splitpdf_lambda_role = aws.iam.Role(
    "lambdaRole",
    assume_role_policy="""{
        "Version": "2012-10-17",
        "Statement": [{
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }]
    }""",
)
# Attach policy to the Lambda role
role_policy_attachment = aws.iam.RolePolicyAttachment(
    "lambdaRolePolicyAttachment",
    role=splitpdf_lambda_role.name,
    policy_arn=aws.iam.ManagedPolicy.AWS_LAMBDA_BASIC_EXECUTION_ROLE,
)
# Lambda's policy to access S3 buckets
split_pdf_lambda_policy = iam.RolePolicy(
    "splitPdfLambdaRole",
    role=splitpdf_lambda_role.id,
    policy=pulumi.Output.all(
        input_bucket.arn, output_bucket.arn, slack_api_token_parameter.arn
    ).apply(
        lambda args: json.dumps(
            {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Action": ["s3:GetObject"],
                        "Effect": "Allow",
                        "Resource": f"{args[0]}/*",
                    },
                    {
                        "Action": ["s3:PutObject"],
                        "Effect": "Allow",
                        "Resource": f"{args[1]}/*",
                    },
                    {
                        "Action": ["ssm:GetParameter"],
                        "Effect": "Allow",
                        "Resource": f"{args[2]}",
                    },
                    {
                        "Action": ["kms:Decrypt"],
                        "Effect": "Allow",
                        "Resource": "*",
                    },
                ],
            }
        )
    ),
)

# Create a CloudWatch log group
log_group = aws.cloudwatch.LogGroup(
    "splitpdfLogGroup",
    name="/aws/lambda/splitpdf",
    retention_in_days=30,
)

lambda_zip = pulumi.FileArchive("../splitpdf/dist/lambda.zip")

# Lambda function
lambda_func = lambda_.Function(
    "splitpdf",
    name="splitpdf",
    code=lambda_zip,
    handler="splitpdf.main.lambda_handler",
    role=splitpdf_lambda_role.arn,
    runtime="python3.11",
    architectures=["arm64"],
    environment=lambda_.FunctionEnvironmentArgs(
        variables={
            "SLACK_CHANNEL": "#general",
            "OUTPUT_BUCKET": output_bucket.bucket.apply(lambda bucket: bucket),
            "SLACK_API_TOKEN_SSM_NAME": slack_api_token_parameter.name.apply(
                lambda arn: arn
            ),
        }
    ),
    timeout=60,  # seconds
    memory_size=256,  # MB
    opts=ResourceOptions(depends_on=[log_group]),
)

# Lambda permission for S3 to invoke the function
lambda_permission = lambda_.Permission(
    "splitpdfLambdaPermission",
    action="lambda:InvokeFunction",
    function=lambda_func.name,
    principal="s3.amazonaws.com",
    source_arn=input_bucket.arn,
    statement_id="AllowS3Event",
)

# S3 bucket notification for the input bucket to invoke the Lambda function
bucket_notification = s3.BucketNotification(
    "bucketNotification",
    bucket=input_bucket.id,
    lambda_functions=[
        s3.BucketNotificationLambdaFunctionArgs(
            lambda_function_arn=lambda_func.arn,
            events=["s3:ObjectCreated:*"],
            filter_suffix=".pdf",
        )
    ],
)

# Output the names of the buckets and the lambda function
pulumi.export("input_bucket_name", input_bucket.bucket)
pulumi.export("output_bucket_name", output_bucket.bucket)
pulumi.export("lambda_function_name", lambda_func.name)
