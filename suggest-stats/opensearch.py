import csv
import gzip
import json
import os
import re

from opensearchpy import OpenSearch, helpers

CORRECT_LETTERS = re.compile(r'(?i)[^a-zа-яё0-9 ,.\-]+', re.U)
NON_FRACTIONAL_DOTS = re.compile(r'(?<!\d)[.,](?!\d)')
FRACTIONAL_COMMAS = re.compile(r'(?<=\d)[,](?=\d)')
WHITESPACE_NORMALIZE = re.compile(r'[\s\-]+', re.U)
REPLACE_YO = re.compile(r'ё', re.U)

opensearch_queries_len = 18_387_936
suggest_queries_len = 5_227_801


def clean_phrase(phrase: str) -> str:
    tmp_phrase = phrase.lower()
    tmp_phrase = CORRECT_LETTERS.sub(' ', tmp_phrase)
    tmp_phrase = NON_FRACTIONAL_DOTS.sub(' ', tmp_phrase)
    tmp_phrase = FRACTIONAL_COMMAS.sub('.', tmp_phrase)
    tmp_phrase = WHITESPACE_NORMALIZE.sub(' ', tmp_phrase)
    tmp_phrase = tmp_phrase.strip()
    tmp_phrase = REPLACE_YO.sub('е', tmp_phrase)
    return tmp_phrase


INDEX_NAME = "search_suggest"
BATCH_SIZE = 5000

opensearch = OpenSearch(
    hosts=[{"host": "158.160.167.113", "port": 9200}],
    http_compress=True
)


def generate_suggestions(file_path):
    with gzip.open(file_path, mode='rt', encoding="utf-8") as f:
        print('Загружаем json...')
        index_json = json.load(f)
        print('Обработка данных...')
        for element in index_json:
            q = element['query']
            right_q = element['right_query']
            weight = element['searches'] + 15 * element['contacts']

            doc = {
                "suggest": {
                    "input": [q],
                    "weight": weight
                },
                "output": [right_q],
            }

            yield {
                "_index": INDEX_NAME,
                "_source": doc
            }


def index_opensearch(index_file_path):
    suggestions_to_upload = []
    count = 0

    for action in generate_suggestions(index_file_path):
        suggestions_to_upload.append(action)
        count += 1

        if len(suggestions_to_upload) >= BATCH_SIZE:
            helpers.bulk(opensearch, suggestions_to_upload)
            print(f"Загружено {round(count/opensearch_queries_len * 100, 2)}% записей...")
            suggestions_to_upload = []

    if suggestions_to_upload:
        helpers.bulk(opensearch, suggestions_to_upload)
        print(f"Загружено {round(count/opensearch_queries_len * 100, 2)}% записей...")


def prepare_data_to_index():
    grouped_queries = {}
    cnt = 0
    for dir, _, files in os.walk('data'):
        for file in files:
            if file[-2:] != 'gz':
                continue
            cnt += 1
            path = dir + '/' + file
            print(f'processing file {path} №{cnt}')
            with gzip.open(path, mode='rt', encoding="utf-8") as f:
                reader = csv.reader(f, delimiter='\t')
                for row in reader:
                    if len(row) < 2:
                        continue

                    query, searches, contacts, _ = row
                    query = query.strip()
                    searches = int(searches)
                    contacts = int(contacts)

                    if query not in grouped_queries:
                        grouped_queries[query] = {
                            'query': query,
                            'right_query': clean_phrase(query),
                            'searches': searches,
                            'contacts': contacts
                        }
                    grouped_queries[query]['searches'] += searches
                    grouped_queries[query]['contacts'] += searches

    result_to_dump = []
    for q, info in grouped_queries.items():
        result_to_dump.append(info)

    print('dumping res')
    with gzip.open('opensearch_query_info.json.gz', mode='wt', encoding='utf-8') as f:
        json.dump(result_to_dump, f, ensure_ascii=False)


if __name__ == "__main__":
    prepare_data_to_index()
    index_opensearch('opensearch_query_info.json.gz')
