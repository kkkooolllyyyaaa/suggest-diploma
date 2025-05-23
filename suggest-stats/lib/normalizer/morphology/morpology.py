import pymorphy2
from pylru import lrucache


class Morphology:
    def __init__(self):
        self.__morph = pymorphy2.MorphAnalyzer()
        self.__morph_cache = lrucache(30_000)

    def normalize_word(self, word: str) -> str:
        normalized_word = self.__morph_cache.get(word, None)

        if normalized_word:
            return normalized_word

        parsed = self.__morph.parse(word)[0]
        parsed_sing_nomn = parsed.inflect({'sing', 'nomn'})
        parsed_sing_nomn_masc = parsed.inflect({'sing', 'nomn', 'masc'})
        if parsed_sing_nomn_masc:
            normalized_word = parsed_sing_nomn_masc[0]
        elif parsed_sing_nomn:
            normalized_word = parsed_sing_nomn[0]
        else:
            normalized_word = parsed[0]

        self.__morph_cache[word] = normalized_word
        return normalized_word
