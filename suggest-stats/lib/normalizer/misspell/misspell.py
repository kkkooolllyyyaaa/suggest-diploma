import jamspell


class Misspell:
    def __init__(self):
        self.jamspell_corrector = jamspell.TSpellCorrector()
        self.jamspell_corrector.SetPenalty(15.0, 3.75)
        assert self.jamspell_corrector.LoadLangModel('data/misspell/models/model_queriesv4.bin')

    def correct(self, query):
        return self.jamspell_corrector.FixFragment(query)

# Got results:
# python3 evaluate/evaluate.py -a data/misspell/alphabet.txt -jsp data/misspell/models/model_queriesv4.bin -mx 50000 data/misspell/data/test3.txt
# [info] loading models
# [info] loading text
# [info] generating typos
# [info] total words: 30334
# [info]               errRate   fixRate    broken   topNerr   topNfix      time
# [info]   jamspell      2.02%    89.07%     0.38%     1.09%    91.06%   136.24s
# [info]      empty     15.60%     0.00%     0.00%    15.60%     0.00%     0.01s
