import gzip
import os
import shutil


def delete_file(file):
    os.remove(file)


def rename_file(old, new):
    os.rename(old, new)


def compress_to_gz(file):
    with open(file, 'rb') as f_in:
        with gzip.open(file + '.gz', 'wb') as f_out:
            shutil.copyfileobj(f_in, f_out)


def decompress_from_gz(input_gz):
    assert len(input_gz) > 3 and input_gz[-3:] == '.gz'

    with gzip.open(input_gz, 'rb') as f_in:
        with open(input_gz[:len(input_gz) - 3], 'wb') as f_out:
            shutil.copyfileobj(f_in, f_out)
