from lib.normalizer.misspell.misspell import Misspell
from lib.normalizer.morphology.morpology import Morphology
from lib.util.text import clean_phrase


class Normalizer:
    def __init__(self):
        self.morphology = Morphology()
        self.misspell = Misspell()

    def strong_normalize(self, query):
        clean_query = clean_phrase(query)
        misspell_corrected = self.misspell.correct(clean_query)

        tokens = [self.morphology.normalize_word(word) for word in misspell_corrected.split()]
        return ' '.join(sorted(tokens))

    def soft_normalize(self, query):
        return clean_phrase(query)
