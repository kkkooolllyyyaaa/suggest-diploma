import csv
import gzip
import os
import random

from lib.util.text import clean_phrase


def gen_from_raw_data(process_row):
    for dir, _, files in os.walk('../../../data/storage/raw'):
        for file in files:
            last2 = file[-2:]
            if last2 != 'gz':
                continue
            with gzip.open(dir + '/' + file, 'rt') as f:
                print('Processing file', file)
                reader = csv.reader(f, delimiter='\t')
                for row in reader:
                    process_row(row)


def gen_test_data():
    strings = []

    def process_row(row):
        searches = int(row[1])
        if random.random() * searches >= 9.999 and random.random() >= 0.99:
            q = clean_phrase(row[0])
            strings.append(q)
            for _ in range(searches):
                pass

    gen_from_raw_data(process_row)

    text = '; '.join(strings)
    with open("../../../data/misspell/data/test.txt", "w") as f:
        f.write(text)


def gen_train_data():
    strings = []

    def process_row(row):
        searches = int(row[1])
        q = clean_phrase(row[0])
        for _ in range(searches):
            strings.append(q)

    gen_from_raw_data(process_row)

    text = '; '.join(strings)
    with open("../../../data/misspell/data/trainv4.txt", "w") as f:
        f.write(text)


if __name__ == '__main__':
    # gen_train_data()
    gen_test_data()
