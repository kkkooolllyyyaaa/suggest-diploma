import re

CORRECT_LETTERS = re.compile(r'(?i)[^a-zа-яё0-9 ,.\-]+', re.U)
NON_FRACTIONAL_DOTS = re.compile(r'(?<!\d)[.,](?!\d)')
FRACTIONAL_COMMAS = re.compile(r'(?<=\d)[,](?=\d)')
WHITESPACE_NORMALIZE = re.compile(r'[\s\-]+', re.U)
REPLACE_YO = re.compile(r'ё', re.U)


def clean_phrase(phrase: str) -> str:
    tmp_phrase = phrase.lower()
    tmp_phrase = CORRECT_LETTERS.sub(' ', tmp_phrase)
    tmp_phrase = NON_FRACTIONAL_DOTS.sub(' ', tmp_phrase)
    tmp_phrase = FRACTIONAL_COMMAS.sub('.', tmp_phrase)
    tmp_phrase = WHITESPACE_NORMALIZE.sub(' ', tmp_phrase)
    tmp_phrase = tmp_phrase.strip()
    tmp_phrase = REPLACE_YO.sub('е', tmp_phrase)
    return tmp_phrase
