import json
import typing as T
from collections import OrderedDict

from lib.nodes.tree import NodesTree, ROOT_ID

MIN_CAT_SCORE = 4
MIN_QUERY_CATEGORIES_SCORE = 80
ROUND_DIGITS = 3


class QueriesCategoriesInfo:
    def __init__(self, tree: NodesTree):
        self.queries_categories: T.Dict[str, QueryCategories] = {}
        self.tree = tree

    def add(self, q, node_id, searches, contacts):
        if q not in self.queries_categories:
            self.queries_categories[q] = QueryCategories()
        self.queries_categories[q].add(node_id, searches, contacts)

    def propagate_all(self):
        for _, nodes_stats in self.queries_categories.items():
            nodes_stats.propagate_stats_to_parents(self.tree)
            nodes_stats.sort_by_score(self.tree)

    def calc_features_all(self):
        for nodeId, nodes_stats in self.queries_categories.items():
            nodes_stats.calc_features()
            nodes_stats.sort_by_score(self.tree)


class QueryCategories:
    __slots__ = ['nodes', 'total_stats']

    def __init__(self):
        self.nodes = OrderedDict()
        self.total_stats: QueryCategories.NodeFreqs = None

    def add(self, node_id, searches, contacts):
        if node_id not in self.nodes:
            self.nodes[node_id] = QueryCategories.NodeFreqs(0, 0)
        self.nodes[node_id].searches += searches
        self.nodes[node_id].contacts += contacts

    def filter_small_nodes(self):
        new_dict = OrderedDict()
        sum_score = 0
        for node, stats in self.nodes.items():
            score = stats.score()
            if score >= MIN_CAT_SCORE:
                new_dict[node] = stats
                sum_score += score
        self.nodes = new_dict
        return sum_score

    def sort_by_score(self, tree: NodesTree):
        ordered_nodes = sorted(
            [node for node in self.nodes],
            reverse=True, key=lambda x: (self.nodes[x].score(), -tree.depth(x)),
        )
        new_dict = OrderedDict()
        for node in ordered_nodes:
            new_dict[node] = self.nodes[node]
        self.nodes = new_dict

    def propagate_stats_to_parents(self, tree: NodesTree):
        propagated_nodes = OrderedDict()
        total_stats = QueryCategories.NodeFreqs(0, 0)

        for node_id, stats in self.nodes.items():
            total_stats.searches += stats.searches
            total_stats.contacts += stats.contacts

            propagated_nodes[node_id] = stats

            parents = tree.get_parents(node_id)
            for parent in parents:
                if parent not in propagated_nodes:
                    propagated_nodes[parent] = QueryCategories.NodeFreqs(0, 0)

                propagated_nodes[parent].searches += stats.searches
                propagated_nodes[parent].contacts += stats.contacts

        self.total_stats = total_stats
        self.nodes = propagated_nodes

    def calc_features(self):
        total_contacts = self.total_stats.contacts
        total_searches = self.total_stats.searches
        total_score = query_score(total_searches, total_contacts)
        for _, stats in self.nodes.items():
            stats.calc_features(total_contacts, total_searches, total_score)

    class NodeFreqs:
        __slots__ = ['searches', 'contacts', 'features']

        def __init__(self, searches, contacts, features=None):
            self.searches = searches
            self.contacts = contacts
            self.features: QueryCategories.Features = features

        def score(self):
            return query_score(self.searches, self.contacts)

        def calc_features(self, total_contacts, total_searches, total_score):
            score = self.score()
            self.features = QueryCategories.Features(
                total_contacts=total_contacts,
                total_searches=total_searches,
                total_score=total_score,

                node_contacts=self.contacts,
                node_searches=self.searches,
                node_score=score,

                node_contact_rate=self.contacts / total_contacts if total_contacts != 0 else 0.0,
                node_search_rate=self.searches / total_searches if total_searches != 0 else 0.0,
                node_score_rate=score / total_score if total_score != 0 else 0.0,
            )

    class Features:
        __slots__ = [
            'total_contacts', 'node_contacts', 'node_contact_rate',
            'total_searches', 'node_searches', 'node_search_rate',
            'total_score', 'node_score', 'node_score_rate',
        ]

        def __init__(
                self,
                total_contacts, node_contacts, node_contact_rate,
                total_searches, node_searches, node_search_rate,
                total_score, node_score, node_score_rate,
        ):
            self.total_contacts = total_contacts
            self.node_contacts = node_contacts
            self.node_contact_rate = node_contact_rate

            self.total_searches = total_searches
            self.node_searches = node_searches
            self.node_search_rate = node_search_rate

            self.total_score = total_score
            self.node_score = node_score
            self.node_score_rate = node_score_rate


class QueriesCategoriesEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, QueriesCategoriesInfo):
            return {
                q: self.default(node_stats)
                for q, node_stats in obj.queries_categories.items()
            }
        elif isinstance(obj, QueryCategories):
            has_not_features = any(True for _, stats in obj.nodes.items() if stats.features is None)
            if has_not_features:
                return [
                    {
                        'node_id': node_id,
                        'total_searches': stats.searches,
                        'total_contacts': stats.contacts,
                    }
                    for node_id, stats in obj.nodes.items() if node_id != ROOT_ID
                ]

            return [
                {
                    'node_id': node_id,

                    'total_contacts': stats.features.total_contacts,
                    'node_contacts': stats.features.node_contacts,
                    'node_contact_rate': round(stats.features.node_contact_rate, ROUND_DIGITS),

                    'total_searches': stats.features.total_searches,
                    'node_searches': stats.features.node_searches,
                    'node_search_rate': round(stats.features.node_search_rate, ROUND_DIGITS),

                    'total_score': stats.features.total_score,
                    'node_score': stats.features.node_score,
                    'node_score_rate': round(stats.features.node_score_rate, ROUND_DIGITS),
                }
                for node_id, stats in obj.nodes.items() if node_id != ROOT_ID
            ]
        return super().default(obj)


def query_score(searches, contacts):
    return searches + 15 * contacts
