import csv
import gzip
import json

from gensim.models import Word2Vec

from lib.context import Context
from lib.util.text import clean_phrase


def prepare_dataset(ctx):
    result_dataset = []
    unique_tokens = set()

    filename = ctx.storage.local.queries()
    with gzip.open('../../' + filename, mode='rt', encoding='utf-8') as f:
        queries_json = json.load(f)
        for row in queries_json:
            q = clean_phrase(row['query'])
            right_q = row['right_query']
            searches = int(row['searches'])
            contacts = int(row['contacts'])

            tokens = q.split()
            if len(tokens) != 1:
                for token in tokens:
                    unique_tokens.add(token)
                for _ in range(searches):
                    result_dataset.append(tokens)

            if right_q == q:
                continue

            tokens = right_q.split()
            if len(tokens) != 1:
                for token in tokens:
                    unique_tokens.add(token)
                for _ in range(searches + contacts):
                    result_dataset.append(tokens)
    return result_dataset, unique_tokens


dim = 128


def train(ctx, model_name):
    ctx.logger.info("Preparing dataset for word2vec")
    dataset, unique_tokens = prepare_dataset(ctx)
    ctx.logger.info(f"Got {len(dataset)} sentences with {len(unique_tokens)} unique tokens")
    model = Word2Vec(
        sentences=dataset,
        vector_size=dim,
        window=7,
        min_count=20,
        workers=8,
        sg=1,
        epochs=8,
    )
    ctx.logger.info("Saving model")
    model.save(model_name)
    return unique_tokens


def train_word2vec(ctx: Context, model_name='word2vec_large.model.v2'):
    ctx.logger.info('Start training Word2vec')
    unique_tokens = train(ctx, model_name)
    ctx.logger.info("Loading model")
    model = Word2Vec.load(model_name)

    with open('../../data/vector/tokens_vectors.tsv', mode='w') as f:
        writer = csv.writer(f, delimiter='\t')
        for token in unique_tokens:
            try:
                vector = model.wv.get_vector(token)
                writer.writerow([token, vector])
            except KeyError:
                pass
