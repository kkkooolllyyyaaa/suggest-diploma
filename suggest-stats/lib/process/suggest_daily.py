import csv
import gzip
import multiprocessing
import os
from concurrent.futures import ThreadPoolExecutor
from multiprocessing import JoinableQueue

from lib.config import LOCAL_ARTIFACTS_GZ
from lib.context import Context
from lib.process.process import query_score
from lib.util.collections import EffectiveList

# to fix multiprocessing on macOS with arm chip
multiprocessing.set_start_method("fork")


def process_daily_file(ctx: Context, date):
    # 0) check, maybe file was already processed
    queries_processed = ctx.storage.check_file_exists(LOCAL_ARTIFACTS_GZ.queries_daily(date))
    queries_categories_processed = ctx.storage.check_file_exists(LOCAL_ARTIFACTS_GZ.queries_categories_daily(date))
    if queries_processed and queries_categories_processed:
        ctx.logger.info(f'Date {str(date)} was already processed')
        return

    # 1) sort rows by query
    ctx.logger.info(f"sorting queries: {str(date)}")
    sorted_rows = []
    with gzip.open(LOCAL_ARTIFACTS_GZ.queries_daily_raw(date), mode='rt', encoding='utf-8') as f_read:
        reader = csv.reader(f_read, delimiter='\t')
        sorted_rows = sorted(map(
            lambda x: (x[0], int(x[1]), int(x[2]), x[3]),
            reader,
        ), key=lambda row: (row[0], -query_score(row[1], row[2])))

    # 2) match same queries + normalize
    ctx.logger.info(f"matching same queries: {str(date)}")
    grouped_rows = EffectiveList(len(sorted_rows))
    prev_q = None
    total_searches, total_contacts = 0, 0

    for q, searches, contacts, category in sorted_rows:
        if prev_q and prev_q != q:
            normalized_q = ctx.normalizer.strong_normalize(prev_q)
            grouped_rows.append([prev_q, normalized_q, total_searches, total_contacts])
            total_searches = 0
            total_contacts = 0

        total_searches += searches
        total_contacts += contacts
        prev_q = q

    normalized_q = ctx.normalizer.strong_normalize(prev_q)
    grouped_rows.append([prev_q, normalized_q, total_searches, total_contacts])
    grouped_rows = grouped_rows.get()

    # 3) dump queries to file
    ctx.logger.info(f"dumping to file and storage: {str(date)}")
    with gzip.open(LOCAL_ARTIFACTS_GZ.queries_daily(date), mode='wt', encoding='utf-8') as f_write:
        writer = csv.writer(f_write, delimiter='\t')
        for row in grouped_rows:
            writer.writerow(row)

    # 4) process queries categories
    ctx.logger.info(f"processing queries categories: {str(date)}")
    category_rows = EffectiveList(len(sorted_rows) // 10)
    prev_q = ''
    nodes_stats = {}
    for q, searches, contacts, category in sorted_rows:
        agg_node_id = category
        if q != prev_q:
            for node_id in nodes_stats:
                category_rows.append([prev_q, node_id, nodes_stats[node_id][0], nodes_stats[node_id][1]])
            nodes_stats = {agg_node_id: [searches, contacts]}
        else:
            if agg_node_id not in nodes_stats:
                nodes_stats[agg_node_id] = [0, 0]
            nodes_stats[agg_node_id][0] += searches
            nodes_stats[agg_node_id][1] += contacts
        prev_q = q

    for node_id in nodes_stats:
        category_rows.append([prev_q, node_id, nodes_stats[node_id][0], nodes_stats[node_id][1]])

    category_rows = sorted(category_rows.get(), key=lambda row: (row[0], row[1]), reverse=True)

    # 5) dump queries categories to file
    ctx.logger.info(f"dumping query categories: {str(date)}")
    with gzip.open(LOCAL_ARTIFACTS_GZ.queries_categories_daily(date), mode='wt') as f_write:
        writer = csv.writer(f_write, delimiter='\t')
        for row in category_rows:
            writer.writerow(row)


def process_raw_data(ctx: Context, dates):
    processing_queue = JoinableQueue()
    uploading_queue = JoinableQueue()

    # download files
    def download_file_worker(date):
        ctx.storage.download_daily_raw(date)
        processing_queue.put(date)

    storage_workers_cnt = min(ctx.cfg.pipeline.storage_workers, len(dates))
    with ThreadPoolExecutor(max_workers=storage_workers_cnt) as download_executor:
        download_executor.map(download_file_worker, dates)

    # process files
    def process_file_worker():
        new_context = ctx.copy()
        while True:
            date = processing_queue.get()
            if date is None:
                processing_queue.task_done()
                break

            process_daily_file(new_context, date)
            uploading_queue.put(date)
            processing_queue.task_done()

    process_workers = []
    process_workers_cnt = min(ctx.cfg.pipeline.process_workers, len(dates))

    for _ in range(process_workers_cnt):
        p = multiprocessing.Process(target=process_file_worker)
        p.start()
        process_workers.append(p)

    # upload to storage
    def upload_file_worker():
        while True:
            date = uploading_queue.get()
            if date is None:
                uploading_queue.task_done()
                break

            ctx.storage.upload_queries_daily(date)
            ctx.storage.upload_queries_categories_daily(date)
            uploading_queue.task_done()

    upload_workers = []
    for _ in range(storage_workers_cnt):
        p = multiprocessing.Process(target=upload_file_worker)
        p.start()
        upload_workers.append(p)

    # await
    for _ in range(process_workers_cnt):
        processing_queue.put(None)
    processing_queue.join()
    for p in process_workers:
        p.join()
    ctx.logger.info('Process done')

    for _ in range(storage_workers_cnt):
        uploading_queue.put(None)
    uploading_queue.join()
    for p in upload_workers:
        p.join()
    ctx.logger.info('Upload done')

    return


def process(ctx: Context):
    dates_to_process = ctx.storage.get_dates_to_process_raw()
    if len(dates_to_process) > 0:
        ctx.logger.info(f"Going to process {len(dates_to_process)} dates: {[str(x) for x in dates_to_process]}")
        process_raw_data(ctx, dates_to_process)
    else:
        ctx.logger.info("Got 0 days to process, skipping...")


def upload_raw_data_to_storage(ctx: Context):
    existing = set(ctx.storage.s3.list_files("raw"))
    for dir, _, files in os.walk('data/raw'):
        for file in files:
            last2 = file[-2:]
            if last2 != 'gz':
                continue

            local = dir + '/' + file
            remote = 'raw/' + file

            if remote not in existing:
                ctx.logger.debug(f'Загружаем файл {local} на {remote}')
                ctx.storage.s3.upload_file(local, remote)


def upload_processed_data_to_storage(ctx: Context):
    existing = set(ctx.storage.s3.list_files("process"))
    for dir, _, files in os.walk('data/storage/process'):
        for file in files:
            last2 = file[-2:]
            if last2 != 'gz':
                continue

            local = dir + '/' + file
            remote = 'process/' + file

            if remote not in existing:
                ctx.logger.debug(f'Загружаем файл {local} на {remote}')
                ctx.storage.s3.upload_file(local, remote)
