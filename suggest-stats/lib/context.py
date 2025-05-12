import logging

from lib import config
from lib.nodes.tree import NodesTree
from lib.normalizer.normalizer import Normalizer
from lib.storage.suggest_storage import SuggestStorage
from lib.util import nodes

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('app.log'),
        logging.StreamHandler(),
    ]
)


class Context:
    def __init__(self, context=None):
        if context:
            self.logger = context.logger
            self.cfg = context.cfg
            self.storage = context.storage
            self.normalizer = Normalizer()
            self.tree = context.tree
            return

        self.logger = logging.getLogger(__name__)
        self.cfg = config.Config()
        self.storage = SuggestStorage(self.cfg.storage, self.logger)
        self.normalizer = Normalizer()
        self.tree = NodesTree(nodes.read_nodes(self.cfg.nodes_path))

    def copy(self):
        return Context(self)
