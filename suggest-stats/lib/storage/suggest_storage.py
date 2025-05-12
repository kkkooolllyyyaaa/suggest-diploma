import datetime
import logging
import os
import re

from lib.config import StorageConfig
from lib.storage.s3_storage import S3

DATES_TO_MERGE = 30


class SuggestStorage:
    def __init__(self, cfg: StorageConfig, logger: logging.Logger):
        self.s3 = S3(cfg.s3_credentials, logger)
        self.remote = cfg.remote_artifacts
        self.local = cfg.local_artifacts
        self.logger = logger

    def get_dates_to_process_raw(self):
        dates_candidates = self.get_dates_from_storage(self.remote.raw_data_folder())
        processed_dated = self.get_dates_from_storage(self.remote.processed_data_folder())

        dates_to_process = []
        for date in dates_candidates:
            if date not in processed_dated:
                dates_to_process.append(date)
        return dates_to_process

    def get_dates_to_merge(self):
        processed_dates = self.get_dates_from_storage(self.remote.processed_data_folder())
        processed_dates = processed_dates[:min(len(processed_dates), DATES_TO_MERGE)]
        self.logger.info(f'Total {len(processed_dates)} days to process merge: {[str(x) for x in processed_dates]}')
        return processed_dates

    def get_dates_to_download_for_merge(self, dates):
        dates_to_download = []
        for date in dates:
            is_queries_processed = self.check_file_exists(self.local.queries_daily(date))
            is_categories_processed = self.check_file_exists(self.local.queries_categories_daily(date))
            if is_queries_processed and is_categories_processed:
                continue
            dates_to_download.append(date)
        return dates_to_download

    def download_daily_raw(self, date):
        local = self.local.queries_daily_raw(date)
        if self.check_file_exists(local):
            return
        self.logger.info(f'Loading raw data for {str(date)}')
        self.s3.download_file(local, self.remote.queries_daily_raw(date))

    def download_queries_daily(self, date):
        local = self.local.queries_daily(date)
        if self.check_file_exists(local):
            return
        self.s3.download_file(local, self.remote.queries_daily(date))

    def download_queries_categories_daily(self, date):
        local = self.local.queries_categories_daily(date)
        if self.check_file_exists(local):
            return
        self.s3.download_file(local, self.remote.queries_categories_daily(date))

    def upload_queries_daily(self, date):
        self.s3.upload_file(
            self.local.queries_daily(date),
            self.remote.queries_daily(date)
        )

    def upload_queries(self):
        self.s3.upload_file(self.local.queries(), self.remote.queries())

    def upload_queries_categories(self):
        self.s3.upload_file(self.local.queries_categories(), self.remote.queries_categories())

    def upload_queries_categories_daily(self, date):
        self.s3.upload_file(
            self.local.queries_categories_daily(date),
            self.remote.queries_categories_daily(date)
        )

    def get_dates_from_storage(self, folder):
        files = self.s3.list_files(folder)
        dates_str = map(lambda x: re.search("([0-9]{4}\-[0-9]{2}\-[0-9]{2})", x).group(1), files)
        dates_set = set(map(lambda x: datetime.datetime.strptime(x, '%Y-%m-%d').date(), dates_str))
        return sorted(list(dates_set), reverse=True)

    def check_file_exists(self, filename):
        exists = os.path.exists(filename)
        if exists:
            self.logger.debug(f'File {filename} already exists')
        return exists
