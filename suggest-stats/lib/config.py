import os


class PipelineConfig:
    def __init__(self):
        self.storage_workers = 1
        self.process_workers = 10


class S3Credentials:
    def __init__(self):
        self.endpoint = os.getenv('S3_ENDPOINT', 'localhost:9000')
        self.access_key = os.getenv('S3_ACCESS_KEY', 'minioadmin')
        self.secret_key = os.getenv('S3_SECRET_KEY', 'minioadmin')
        self.bucket = 'suggest'
        self.secure = False


class StorageConfig:
    def __init__(self):
        self.s3_credentials = S3Credentials()
        self.remote_artifacts = REMOTE_ARTIFACTS
        self.local_artifacts = LOCAL_ARTIFACTS_GZ


class Artifacts:
    def __init__(self, root, suffix):
        self.root = root
        self.suffix = suffix

    def queries(self):
        return self.__build_artifact_path('artifact/queries', 'json' + self.suffix)

    def queries_glued(self):
        return self.__build_artifact_path('artifact/queries', 'tsv.glued' + self.suffix)

    def queries_daily(self, date):
        return self.__build_path('process', 'queries', date, 'tsv' + self.suffix)

    def queries_daily_raw(self, date):
        return self.__build_path('raw', 'queries', date, 'tsv' + self.suffix)

    def queries_categories(self):
        return self.__build_artifact_path('artifact/queries_categories', 'json' + self.suffix)

    def queries_categories_propagated(self):
        return self.__build_artifact_path('artifact/queries_categories_propagated', 'json' + self.suffix)

    def queries_categories_glued(self):
        return self.__build_artifact_path('artifact/queries_categories', 'tsv.glued' + self.suffix)

    def queries_categories_daily(self, date):
        return self.__build_path('process', 'queries_categories', date, 'tsv' + self.suffix)

    def raw_data_folder(self):
        return os.path.join(self.root, 'raw')

    def processed_data_folder(self):
        return os.path.join(self.root, 'process')

    def __build_artifact_path(self, path, extension):
        return self.__build_path(path, None, None, extension)

    def __build_path(self, path, filename, date, extension):
        result_path = os.path.join(self.root, path)
        if filename is not None:
            result_path = os.path.join(result_path, filename)
        if date is not None:
            result_path += '-' + str(date)
        if extension is not None:
            result_path += '.' + extension
        return result_path


remote_root = ''
REMOTE_ARTIFACTS = Artifacts(remote_root, suffix='.gz')

local_root = 'data/storage'
LOCAL_ARTIFACTS_GZ = Artifacts(local_root, suffix='.gz')
LOCAL_ARTIFACTS = Artifacts(local_root, suffix='')


class Config:
    def __init__(self):

        self.storage = StorageConfig()
        self.pipeline = PipelineConfig()
        self.nodes_path = 'data/nodes.json.gz'
