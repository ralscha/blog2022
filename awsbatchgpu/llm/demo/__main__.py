import base64
import json
import os
from concurrent.futures import ThreadPoolExecutor, as_completed

import boto3
from botocore.exceptions import ClientError
from mypy_boto3_s3.client import S3Client
from vllm import LLM
from vllm.entrypoints.chat_utils import ChatCompletionMessageParam
from vllm.sampling_params import SamplingParams


def download_large_file(
    s3_client: S3Client,
    bucket_name: str,
    object_key: str,
    file_path: str,
    part_size: int = 4 * 1024 * 1024 * 1024,
    max_workers: int = 5,
):
    try:
        response = s3_client.head_object(Bucket=bucket_name, Key=object_key)
        file_size = response["ContentLength"]

        with open(file_path, "wb") as f:
            f.seek(file_size - 1)
            f.write(b"\0")

        parts = [
            (i, min(i + part_size - 1, file_size - 1))
            for i in range(0, file_size, part_size)
        ]

        def download_part(start_byte: int, end_byte: int):
            response = s3_client.get_object(
                Bucket=bucket_name,
                Key=object_key,
                Range=f"bytes={start_byte}-{end_byte}",
            )
            with open(file_path, "r+b") as f:
                f.seek(start_byte)
                f.write(response["Body"].read())

        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = [
                executor.submit(download_part, start, end) for start, end in parts
            ]
            for future in as_completed(futures):
                future.result()
    except ClientError as e:
        print(f"Error downloading file {object_key}: {e}")
        raise


def setu_pmodel(s3_client: S3Client, bucket_name: str, prefix_key: str) -> LLM:
    model_path = "/tmp/pixtral"
    os.makedirs(model_path, exist_ok=True)

    #download_large_file(
    #    s3_client,
    #    bucket_name,
    #    prefix_key + "consolidated.safetensors",
    #    model_path + "/consolidated.safetensors",
    #)
    s3_client.download_file(
        bucket_name, prefix_key + "consolidated.safetensors", model_path + "/consolidated.safetensors"
    )
    s3_client.download_file(
        bucket_name, prefix_key + "params.json", model_path + "/params.json"
    )
    s3_client.download_file(
        bucket_name, prefix_key + "tekken.json", model_path + "/tekken.json"
    )

    return LLM(
        model=model_path,
        tokenizer_mode="mistral",
        tensor_parallel_size=4,
        max_model_len=101792,
    )


def process_image(
    llm: LLM, s3_client: S3Client, bucket_name: str, prompt: str, image: dict[str, str]
) -> dict[str, str]:
    try:
        response = s3_client.get_object(Bucket=bucket_name, Key=image["s3Key"])
        image_data = response["Body"].read()
        image_base64 = base64.b64encode(image_data).decode("utf-8")

        messages: list[ChatCompletionMessageParam] = [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": prompt},
                    {"type": "image_url", "image_url": {"url": image_base64}},
                ],
            },
        ]

        sampling_params = SamplingParams(max_tokens=8192)
        outputs = llm.chat(messages, sampling_params=sampling_params)
        completion_response = outputs[0].outputs[0]
        response_text = completion_response.text

        return {"id": image["id"], "response": response_text}
    except Exception as e:
        print(f"Error processing image {image['id']}: {e}")
        return {"id": image["id"], "response": f"Error: {str(e)}"}


def main():
    s3_client: S3Client = boto3.client("s3")
    bucket_name = os.environ["WORK_BUCKET"]
    input_key = "input.json"
    output_key = "output.json"
    model_prefix_key = "pixtral/"

    try:
        response = s3_client.get_object(Bucket=bucket_name, Key=input_key)
        input_data = json.loads(response["Body"].read().decode("utf-8"))

        llm = setup_model(s3_client, bucket_name, model_prefix_key)

        results = []
        for image in input_data["images"]:
            result = process_image(
                llm, s3_client, bucket_name, input_data["prompt"], image
            )
            results.append(result)
            print(f"Processed image {image['id']}")

        output_data = json.dumps(results, indent=2)
        s3_client.put_object(Bucket=bucket_name, Key=output_key, Body=output_data)

    except Exception as e:
        print(f"An error occurred: {e}")


if __name__ == "__main__":
    main()
