import gzip
import json
import os

import requests

from lib import context
from lib.process import suggest_daily, suggest_merge
from lib.util.file import compress_to_gz
from lib.word2vec.annoy import calculate_vectors_and_build_index


def sample_morph():
    import csv
    from lib.normalizer.morphology.morpology import Morphology
    queries = []
    morph = Morphology()
    file = 'data/raw/queries-2025-03-02.tsv'
    with open(file, 'r') as f:
        reader = csv.reader(f, delimiter='\t')
        for row in reader:
            queries.append(row[0])
            if len(queries) > 5000:
                break
    with open('sample.txt', 'w') as f:
        writer = csv.writer(f, delimiter='\t')
        for q in queries:
            writer.writerow([q, ' '.join([morph.normalize_word(word) for word in q.split()])])


def filter_raw_data():
    import csv
    from lib.util import nodes
    nodes_by_id = nodes.nodes_map(nodes.read_nodes('data/nodes.json.gz'))
    for dir, _, files in os.walk('data/raw'):
        print(dir, files)
        for file in files:
            path = dir + '/' + file
            rows = []
            with open(path, 'r') as f:
                reader = csv.reader(f, delimiter='\t')
                for row in reader:
                    q, s, c, node = row
                    node_shift = str(int(node) - nodes.ID_SHIFT)
                    if node_shift in nodes_by_id:
                        rows.append([q, s, c, node_shift])

            os.remove(path)
            with open(path, 'w') as f:
                writer = csv.writer(f, delimiter='\t')
                for row in rows:
                    writer.writerow(row)


def create_filtered_nodes():
    from lib.util import nodes
    nodes_to_filter = [
        1047505,  # Услуги
        1047597,  # Животные
        1047789,  # Транспорт
        1054803,  # Недвига
        1057415,  # Готовый бизнес
        1057433,  # Работа
        4291447,  # Путешествия
    ]

    print('creating nodes...')
    nodes.create_nodes(nodes_to_filter)


def test_misspell():
    from lib.normalizer.misspell import misspell
    misspell_corrector = misspell.Misspell()
    test_queries = [
        'тюль без',
        '8110',
        'мерседес 221 двери',
        'новых горизонтов',
        'рефрежииратор',
        'октавия а5 пепельница',
        'противотуманки ф82',
        'алекс фергюсон уроки лидерство',
        'купить macbook pro max сеерый без прашивки безплатно',
        'компьютерный икеа',
        'аренда однжды',
    ]

    for q in test_queries:
        print(q)
        corrected = misspell_corrector.correct(q)
        if q != corrected:
            print(corrected)
        print()

    while True:
        print(misspell_corrector.correct(input()))


def build_annoy_index():
    ctx = context.Context()
    calculate_vectors_and_build_index(ctx)

    filename = 'data/vector/query_vectors.json'
    compress_to_gz(filename)

    ctx.storage.s3.upload_file(filename + '.gz', 'vector/query_vectors.json.gz')
    ctx.storage.s3.upload_file('data/vector/annoy_index.ann', 'vector/query_vectors.json.gz')


def process_suggest():
    ctx = context.Context()
    ctx.logger.info("Starting main suggest process job...")
    suggest_daily.process(ctx)
    suggest_merge.process(ctx)


def getOpenSearchSuggest(query):
    url = 'http://158.160.167.113:9200/search_suggest/_search?pretty=true'
    headers = {'Content-Type': 'application/json'}
    data = {
        "suggest": {
            "search-suggest": {
                "prefix": query,
                "completion": {
                    "field": "suggest",
                    "size": 8
                }
            }
        }
    }
    response = requests.post(url, headers=headers, data=json.dumps(data))
    response_json = response.json()

    suggestions = []
    if 'suggest' in response_json and 'search-suggest' in response_json['suggest']:
        for suggestion in response_json['suggest']['search-suggest']:
            for option in suggestion.get('options', []):
                if 'text' in option:
                    suggestions.append(option['text'])
    return suggestion


def getSuggest(query):
    url = 'http://localhost:8080/v1/api/suggest'
    headers = {'Content-Type': 'application/json', 'userId': '1'}
    data = {
        "query": query,
    }
    response = requests.post(url, headers=headers, data=json.dumps(data))
    response_json = response.json()
    return list([x['title'] for x in response_json['items']])


miss_cnt = 0


def test_suggest():
    with open('data/misspell/data/test.txt', mode='r') as f:
        test_queries = set(f.read().split(';'))
        print(len(test_queries))

    full_c = 0
    full_s = 0
    with gzip.open('data/storage/artifact/queries.json.gz', mode='rt') as f:
        queries = json.load(f)
        query_to_json = {}
        for el in queries:
            q = el['query']
            if q in test_queries:
                full_s += el['searches']
                full_c += el['contacts']
            query_to_json[q] = el
    print(full_c, full_s)

    with gzip.open('data/storage/artifact/queries_categories_propagated.json.gz', mode='rt') as f:
        queries_categories_propagated = json.load(f)

    total_tests = 0
    total_searches = 0
    total_contacts = 0
    total_conversion = 0.0
    total_diversity_s = 0.0
    total_diversity_c = 0.0
    total_mis_cnt = 0

    def calc_diversities(suggestions):
        cat_s_stats = {}
        cat_c_stats = {}
        for q in suggestions:
            if q not in queries_categories_propagated:
                continue
            for node_stats in queries_categories_propagated[q]:
                node_id = node_stats['node_id']
                if node_id not in cat_s_stats:
                    cat_s_stats[node_id] = 0.0
                if node_id not in cat_c_stats:
                    cat_c_stats[node_id] = 0.0

                cat_s_stats[node_id] += node_stats['node_searches']
                cat_c_stats[node_id] += node_stats['node_contacts']

        sum_searches = sum(searches for cat, searches in cat_s_stats.items())
        searches_div = 0
        if sum_searches > 0:
            searches_div = 1 - sum([(searches / sum_searches) ** 2 for cat, searches in cat_s_stats.items()])

        contacts_div = 0
        sum_contacts = sum(contacts for cat, contacts in cat_c_stats.items())
        if sum_contacts > 0:
            contacts_div = 1 - sum([(contacts / sum_contacts) ** 2 for cat, contacts in cat_c_stats.items()])
        return searches_div, contacts_div

    def calc_metrics(suggestions):
        miss_cnt = 0
        serp_searches = 0
        serp_contacts = 0
        serp_conversion = 0.0
        for q in suggestions:
            if q not in query_to_json:
                miss_cnt += 1
                continue
            if q not in queries_categories_propagated:
                miss_cnt += 1
                continue
            qe = query_to_json[q]
            serp_searches += qe['searches']
            serp_contacts += qe['contacts']
            serp_conversion += 0.0 if serp_searches <= 0 else serp_contacts / serp_searches

        serp_conversion /= len(suggestions)
        s_div, c_div = calc_diversities(suggestions)
        return serp_searches, serp_contacts, serp_conversion, s_div, c_div, miss_cnt

    print('Start calculating metrics')
    for q in test_queries:
        sgt = getOpenSearchSuggest(q)
        serp_searches, serp_contacts, serp_conversion, serp_diversity_s, serp_diversity_c, miss_cnt = calc_metrics(sgt)
        total_searches += serp_searches
        total_contacts += serp_contacts
        total_conversion += serp_conversion
        total_diversity_s += serp_diversity_s
        total_diversity_c += serp_diversity_c
        total_mis_cnt += miss_cnt

    print(total_searches, total_contacts, total_conversion, total_diversity_s, total_diversity_c, total_mis_cnt)
    print(
        total_searches / len(test_queries),
        total_contacts / len(test_queries),
        total_conversion / len(test_queries),
        total_diversity_s / len(test_queries),
        total_diversity_c / len(test_queries),
        total_mis_cnt / len(test_queries),
    )


if __name__ == '__main__':
    process_suggest()
    # sample_morph()
    # test_misspell()
    # test_suggest()
