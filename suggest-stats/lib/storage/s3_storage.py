import logging

from minio import Minio
from minio.error import S3Error

from lib.config import S3Credentials


class S3:
    def __init__(self, s3_credentials: S3Credentials, logger: logging.Logger):
        self.logger = logger

        self.s3 = Minio(
            s3_credentials.endpoint,
            access_key=s3_credentials.access_key,
            secret_key=s3_credentials.secret_key,
            secure=s3_credentials.secure,
        )
        self.bucket = s3_credentials.bucket
        self.ensure_bucket_exists()

    def upload_file(self, local_filepath, remote_filepath):
        try:
            self.s3.fput_object(bucket_name=self.bucket, file_path=local_filepath, object_name=remote_filepath)
            self.logger.info(f"File '{local_filepath}' was uploaded as '{remote_filepath}'")
        except S3Error as e:
            self.logger.error(f'Error uploading file {local_filepath}: {e}')

    def list_files(self, prefix):
        try:
            objects = self.s3.list_objects(self.bucket, prefix=prefix, recursive=True)
            return list(map(lambda x: x.object_name, objects))
        except S3Error as e:
            self.logger.error(f'Error getting files list: {e}')

    def download_file(self, local_filepath, remote_filepath):
        try:
            self.s3.fget_object(self.bucket, remote_filepath, local_filepath)
            self.logger.info(f"File '{remote_filepath}' downloaded as '{local_filepath}'")
        except S3Error as e:
            self.logger.error(f"Error downloading file: {e}")

    def ensure_bucket_exists(self):
        found = self.s3.bucket_exists(self.bucket)
        if not found:
            self.s3.make_bucket(self.bucket)
            self.logger.info(f'{self.bucket} was created')
        else:
            self.logger.debug(f'{self.bucket} already exists')
