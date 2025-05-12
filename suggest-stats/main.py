import os

from lib import context
from lib.process import suggest_daily, suggest_merge


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


def process_suggest():
    ctx = context.Context()
    suggest_daily.process(ctx)
    suggest_merge.process(ctx)


if __name__ == '__main__':
    process_suggest()
    # sample_morph()
    # test_misspell()
