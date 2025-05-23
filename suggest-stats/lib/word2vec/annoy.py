import csv
import gzip
import json

import numpy as np

from lib.context import Context
from lib.process.process import query_score
from lib.util.text import clean_phrase

from annoy import AnnoyIndex

dim = 128

stop_words = {
    "и", "в", "во", "не", "что", "он", "на", "я", "с", "со", "как", "а", "то",
    "все", "она", "так", "его", "но", "да", "ты", "к", "у", "же", "вы", "за",
    "бы", "по", "только", "ее", "мне", "было", "вот", "от", "меня", "еще", "нет",
    "о", "из", "ему", "теперь", "когда", "даже", "ну", "вдруг", "ли", "если",
    "уже", "или", "ни", "быть", "был", "него", "до", "вас", "нибудь", "опять",
    "уж", "вам", "ведь", "там", "потом", "себя", "ничего", "ей", "может", "они",
    "тут", "где", "есть", "надо", "ней", "для", "мы", "тебя", "их", "чем", "была",
    "сам", "чтоб", "без", "будто", "чего", "раз", "тоже", "себе", "под", "будет",
    "ж", "тогда", "кто", "этот", "того", "потому", "этого", "какой", "совсем",
    "ним", "здесь", "этом", "один", "почти", "мой", "тем", "чтобы", "нее",
    "были", "куда", "зачем", "всех", "никогда", "можно", "при", "наконец", "два",
    "об", "другой", "хоть", "после", "над", "больше", "тот", "через", "эти", "нас",
    "про", "всего", "них", "какая", "много", "разве", "три", "эту", "моя", "впрочем",
    "хорошо", "свою", "этой", "перед", "иногда", "лучше", "чуть", "том", "нельзя",
    "такой", "им", "более", "всегда", "конечно", "всю", "между", "это", "п", "р",
    "с", "м", "и", "г", "т", "д", "к", "ж", "о", "в", "л", "а", "б", "я", "—", "–",
    "в", "на", "под", "над", "из", "к", "у", "без", "до", "со", "за", "при", "после",
    "между", "ради", "вокруг", "около", "про", "через", "по", "от", "из-за", "из-под",
    "в", "и", "а", "но", "да", "или", "либо", "ни", "не", "же", "что", "чтобы", "как",
    "де", "только", "лишь", "ибо", "так", "тоже", "также", "зато", "притом", "причем",
    "однако", "тем", "не", "все", "всё", "вся", "всю", "всего", "всей", "всем", "всём",
}


def read_tokens_vectors():
    with open('data/vector/tokens_vectors.tsv', 'r') as f:
        reader = csv.reader(f, delimiter='\t')
        tokens_vectors = {}
        for row in reader:
            token = clean_phrase(row[0])
            vector_from_file = row[1][2:len(row[1]) - 2]
            vector = np.array([float(x) for x in vector_from_file.split()])
            tokens_vectors[token] = vector
        return tokens_vectors


def calculate_queries_vectors(ctx):
    tokens_vectors = read_tokens_vectors()

    with gzip.open(ctx.storage.local.queries(), mode='rt', encoding='utf-8') as f:
        queries_json = json.load(f)
        queries_vectors = {}
        idx = 0

        for row in queries_json:
            q = clean_phrase(row['query'])
            searches = int(row['searches'])
            contacts = int(row['contacts'])

            if query_score(searches, contacts) < 500:
                continue

            query_vectors = []
            all_matched = True
            for token in q.split():
                if token in stop_words:
                    continue
                if token not in tokens_vectors:
                    all_matched = False
                    break

                query_vectors.append(tokens_vectors[token])

            if all_matched and len(query_vectors) > 0:
                queries_vectors[q] = {
                    'index': idx,
                    'vector': np.mean(query_vectors, axis=0).tolist(),
                }
                idx += 1

        with open('data/vector/query_vectors.json', 'w') as f:
            json.dump(queries_vectors, f, indent=4, ensure_ascii=False)
        return queries_vectors


def build_annoy_index(queries_vectors):
    annoy_index = AnnoyIndex(dim, 'angular')
    for q, item in queries_vectors.items():
        index = item['index']
        annoy_index.add_item(index, item['vector'])
    annoy_index.build(n_trees=100)
    annoy_index.save('data/vector/annoy_index.ann')


def calculate_vectors_and_build_index(ctx: Context):
    # ctx.logger.info('Calculating queries vectors...')
    # queries_vectors = calculate_queries_vectors(ctx)

    with open('data/vector/query_vectors.json', mode='r') as f:
        queries_vectors = json.load(f)
        ctx.logger.info('Building annoy index...')
        build_annoy_index(queries_vectors)
