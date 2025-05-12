ROOT_ID = '1'


class NodesTree:
    def __init__(self, nodes):
        nodes_map = {}
        parents_map = {}
        for node in nodes:
            id = node['id']
            parentId = node['parentId']
            if parentId != 'None':
                parents_map[id] = parentId
            nodes_map[id] = node

        self.nodes_map = nodes_map
        self.parents_map = parents_map

    def exists(self, nodeId):
        return nodeId in self.nodes_map

    def depth(self, nodeId):
        return len(self.get_parents(nodeId))

    def get_parents(self, node_id):
        result = []
        while node_id in self.parents_map:
            parent = self.parents_map[node_id]
            result.append(parent)
            node_id = parent
        return result
