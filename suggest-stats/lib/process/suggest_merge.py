import csv
import gzip
import json
import multiprocessing
import os
import subprocess
import typing as T

from lib import context
from lib.config import LOCAL_ARTIFACTS_GZ, LOCAL_ARTIFACTS
from lib.process.process import query_score, QueriesCategoriesInfo, QueriesCategoriesEncoder
from lib.util.collections import EffectiveList
from lib.util.file import decompress_from_gz

MAX_QUERY_LEN = 100
MIN_QUERY_SCORE = 12


def unix_sort(params: T.List[str]):
    env = os.environ.copy()
    env['LC_ALL'] = 'C'  # say to `sort`: consider input files as bytes, not as text
    subprocess.call(['sort'] + params, env=env)


def merge_queries(ctx: context.Context, dates_to_process):
    ctx.logger.info('Gluing queries file')
    glued_filename = LOCAL_ARTIFACTS.queries_glued()
    filenames = [LOCAL_ARTIFACTS.queries_daily(date) for date in dates_to_process]
    unix_sort(['-s', '-t', '\t', '-rk', '1,1', '-T', '.', '-o', glued_filename] + filenames)

    ctx.logger.info('Merging glued queries file')
    aggregated_rows = EffectiveList(10_000_000)

    prev_q = ''
    prev_normalized_q = ''
    sum_searches = 0
    sum_contacts = 0
    glued_cnt = 0
    with open(glued_filename, 'r', encoding='utf-8') as f_read:
        reader = csv.reader(f_read, delimiter='\t')
        for row in reader:
            if len(row) != 4:
                continue
            glued_cnt += 1

            q, normalized_q, searches, contacts = row
            if len(q) > MAX_QUERY_LEN:
                continue

            searches = int(searches)
            contacts = int(contacts)

            if q != prev_q:
                if query_score(sum_searches, sum_contacts) > 0:
                    aggregated_rows.append([prev_q, prev_normalized_q, sum_searches, sum_contacts])

                prev_q = q
                prev_normalized_q = normalized_q
                sum_searches = searches
                sum_contacts = contacts
            else:
                sum_searches += searches
                sum_contacts += contacts

        if query_score(sum_searches, sum_contacts) > 0:
            aggregated_rows.append([prev_q, prev_normalized_q, sum_searches, sum_contacts])

    aggregated_rows = aggregated_rows.get()
    ctx.logger.info(f'Got {len(aggregated_rows)} unique queries from {glued_cnt} glued after merging')

    ctx.logger.info(f'Getting most score query for normalization')
    grouped_by_normalized_form = sorted(
        aggregated_rows,
        key=lambda row: (row[1], query_score(row[2], row[3])),
        reverse=True,
    )

    result_to_dump = EffectiveList(len(grouped_by_normalized_form))
    prev_normalized_q = ''
    most_score_q = ''
    max_score = 0
    for q, normalized_q, searches, contacts in grouped_by_normalized_form:
        if prev_normalized_q != normalized_q:
            most_score_q = ctx.normalizer.soft_normalize(q)
            max_score = query_score(searches, contacts)

        if max_score > MIN_QUERY_SCORE:
            result_to_dump.append(
                {
                    'query': q,
                    'right_query': most_score_q,
                    'searches': searches,
                    'contacts': contacts,
                }
            )
        prev_normalized_q = normalized_q

    result_to_dump = result_to_dump.get()
    ctx.logger.info(f'Got {len(result_to_dump)} queries from {len(grouped_by_normalized_form)} grouped')
    ctx.logger.info(f'Dumping queries to file')
    with gzip.open(LOCAL_ARTIFACTS_GZ.queries(), mode='wt', encoding='utf-8') as f_write:
        json.dump(result_to_dump, f_write, ensure_ascii=False)


def merge_queries_categories(ctx: context.Context, dates_to_process):
    ctx.logger.info('Gluing queries categories file')
    filenames = [LOCAL_ARTIFACTS.queries_categories_daily(date) for date in dates_to_process]
    glued_filename = LOCAL_ARTIFACTS.queries_categories_glued()
    unix_sort(['-s', '-t', '\t', '-rk', '1,1', '-rk', '2,2', '-T', '.', '-o', glued_filename] + filenames)

    glued_cnt = 0
    result = QueriesCategoriesInfo(ctx.tree)

    with open(glued_filename, 'r', encoding='utf-8') as f_read:
        reader = csv.reader(f_read, delimiter='\t')
        for row in reader:
            if len(row) != 4:
                continue
            glued_cnt += 1

            q, category, searches, contacts = row
            q = ctx.normalizer.soft_normalize(q)
            searches = int(searches)
            contacts = int(contacts)
            result.add(q, category, searches, contacts)

    ctx.logger.info(f'Got {len(result.queries_categories)} unique queries from {glued_cnt} queries categories rows')

    ctx.logger.info(f'Dumping queries categories to file')
    with gzip.open(LOCAL_ARTIFACTS_GZ.queries_categories(), mode='wt', encoding='utf-8') as f_write:
        json.dump(result, f_write, cls=QueriesCategoriesEncoder, ensure_ascii=False, indent=4)

    ctx.logger.info(f'Propagating stats for queries categories')
    result.propagate_all()
    ctx.logger.info(f'Calculating features for queries categories propagated')
    result.calc_features_all()

    ctx.logger.info(f'Dumping queries categories propagated to file')
    with gzip.open(LOCAL_ARTIFACTS_GZ.queries_categories_propagated(), mode='wt', encoding='utf-8') as f_write:
        json.dump(result, f_write, cls=QueriesCategoriesEncoder, ensure_ascii=False, indent=4)


def process(ctx: context.Context):
    dates_to_process = ctx.storage.get_dates_to_merge()
    dates_to_download = ctx.storage.get_dates_to_download_for_merge(dates_to_process)

    # download files
    if len(dates_to_download) > 0:
        ctx.logger.info(f"Downloading {len(dates_to_download)} files: {[str(x) for x in dates_to_download]}")
        for date in dates_to_download:
            ctx.storage.download_queries_daily(date)
            ctx.storage.download_queries_categories_daily(date)
    else:
        ctx.logger.debug('Nothing to download for merge')

    # uncompress files
    for date in dates_to_process:
        file = LOCAL_ARTIFACTS.queries_daily(date)
        if not ctx.storage.check_file_exists(file):
            ctx.logger.info(f'Decompressing queries file for date {str(date)}')
            decompress_from_gz(LOCAL_ARTIFACTS_GZ.queries_daily(date))

        file = LOCAL_ARTIFACTS.queries_categories_daily(date)
        if not ctx.storage.check_file_exists(file):
            ctx.logger.info(f'Decompressing queries categories file for date {str(date)}')
            decompress_from_gz(LOCAL_ARTIFACTS_GZ.queries_categories_daily(date))

    def queries_worker():
        ctx.logger.info("Merging queries...")
        merge_queries(ctx, dates_to_process)
        ctx.logger.info("Uploading queries...")
        ctx.storage.upload_queries()

    def queries_categories_worker():
        ctx.logger.info("Merging queries categories...")
        merge_queries_categories(ctx, dates_to_process)
        ctx.logger.info("Uploading queries categories...")
        ctx.storage.upload_queries_categories()

    p1 = multiprocessing.Process(target=queries_worker)
    p1.start()

    p2 = multiprocessing.Process(target=queries_categories_worker)
    p2.start()

    ctx.logger.info("Awaiting merge processes")
    p1.join()
    p2.join()
    ctx.logger.info("merge processes done")
