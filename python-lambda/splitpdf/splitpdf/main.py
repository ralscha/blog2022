import os

import boto3
from aws_lambda_powertools.utilities.data_classes import S3Event, event_source
from aws_lambda_powertools.utilities.typing import LambdaContext
from mypy_boto3_s3 import S3Client
from mypy_boto3_ssm.client import SSMClient
from slack_sdk import WebClient

from splitpdf.split import split_pdf

s3_client: S3Client = boto3.client("s3")
output_bucket = os.getenv("OUTPUT_BUCKET", "output-bucket")

slack_api_token_ssm_name = os.getenv("SLACK_API_TOKEN_SSM_NAME", "slack-api-token")
ssm_client: SSMClient = boto3.client("ssm")
slack_token = ssm_client.get_parameter(
    Name=slack_api_token_ssm_name, WithDecryption=True
)["Parameter"]["Value"]

slack_channel = os.getenv("SLACK_CHANNEL", "general")
slack_client = WebClient(token=slack_token)


@event_source(data_class=S3Event)
def lambda_handler(event: S3Event, context: LambdaContext) -> None:
    for record in event.records:
        input_bucket = record.s3.bucket.name
        input_key = record.s3.get_object.key
        output_key_prefix = input_key.replace(".pdf", "")

        split_pdf(s3_client, input_bucket, input_key, output_bucket, output_key_prefix)

        slack_client.chat_postMessage(
            channel=slack_channel,
            text=f"Processed file {input_key} from bucket {input_bucket}",
        )
