import io

from mypy_boto3_s3 import S3Client
from pypdf import PdfReader, PdfWriter


def split_pdf(
    s3_client: S3Client,
    input_bucket: str,
    input_key: str,
    output_bucket: str,
    output_key_prefix: str,
) -> None:
    pdf_stream = io.BytesIO()
    s3_client.download_fileobj(input_bucket, input_key, pdf_stream)
    pdf_stream.seek(0)

    reader = PdfReader(pdf_stream)

    for page_num in range(reader.get_num_pages()):
        writer = PdfWriter()
        writer.add_page(reader.get_page(page_num))

        page_stream = io.BytesIO()
        writer.write(page_stream)
        page_stream.seek(0)

        output_key = f"{output_key_prefix}/{page_num + 1:05d}.pdf"

        s3_client.upload_fileobj(page_stream, output_bucket, output_key)
