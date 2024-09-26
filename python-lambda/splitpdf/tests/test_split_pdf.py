import io
import unittest

import boto3
from moto import mock_aws
from pypdf import PdfWriter

from splitpdf.split import split_pdf


class TestSplitPDF(unittest.TestCase):
    def test_split_pdf(self):
        mock = mock_aws()
        mock.start()

        s3 = boto3.client("s3", region_name="us-east-1")
        s3.create_bucket(Bucket="input-bucket")
        s3.create_bucket(Bucket="output-bucket")

        pdf_stream = io.BytesIO()
        writer = PdfWriter()
        writer.add_blank_page(width=72, height=72)
        writer.add_blank_page(width=72, height=72)
        writer.add_blank_page(width=72, height=72)
        writer.write(pdf_stream)
        pdf_stream.seek(0)

        s3.upload_fileobj(pdf_stream, "input-bucket", "input.pdf")

        split_pdf(s3, "input-bucket", "input.pdf", "output-bucket", "output-directory")

        uploaded_objects = s3.list_objects_v2(Bucket="output-bucket")["Contents"]
        self.assertEqual(len(uploaded_objects), 3)
        self.assertEqual(uploaded_objects[0]["Key"], "output-directory/00001.pdf")
        self.assertEqual(uploaded_objects[1]["Key"], "output-directory/00002.pdf")
        self.assertEqual(uploaded_objects[2]["Key"], "output-directory/00003.pdf")
        mock.stop()


if __name__ == "__main__":
    unittest.main()
