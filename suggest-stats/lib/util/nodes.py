import gzip
import json

ID_SHIFT = 1047503


def read_nodes(path):
    with gzip.open(path, 'rt', encoding='utf-8') as f:
        res = json.load(f)
        return res


def nodes_map(nodes):
    nodes_map = {}
    for node in nodes:
        nodes_map[node['id']] = node
    return nodes_map


def filter_nodes(nodes, filter_nodes):
    filter_nodes_by_id = set()
    for filter_node in filter_nodes:
        filter_nodes_by_id.add(filter_node)

    nodes_by_id = nodes_map(nodes)

    def is_to_filter(node):
        if node in filter_nodes_by_id:
            return True

        parent = nodes_by_id[node]['parentId']
        while parent is not None:
            if parent in filter_nodes_by_id:
                return True
            parent = nodes_by_id[parent]['parentId']
        return False

    new_nodes = []
    for node in nodes:
        if is_to_filter(node['id']):
            continue
        new_nodes.append(node)
    return new_nodes


def filter_fields(nodes):
    new_nodes = []
    for node in nodes:
        new_node = {'id': node['id'], 'parentId': node['parentId'], 'title': node['description']}
        if 'payload' in node:
            if 'title' in node['payload']:
                new_node['title'] = node['payload']['title']

            if 'suggest' in node['payload'] and 'server_icon' in node['payload']['suggest']:
                new_node['server_icon'] = node['payload']['suggest']['server_icon']

        new_nodes.append(new_node)
    return new_nodes


def stringify_ids(nodes):
    for node in nodes:
        node['id'] = str(node['id'])
        node['parentId'] = str(node['parentId'])
    return nodes


def shift_ids(nodes):
    for node in nodes:
        node['id'] = node['id'] - ID_SHIFT
        if 'parentId' in node and node['parentId'] is not None:
            node['parentId'] = node['parentId'] - ID_SHIFT
    return nodes


def create_nodes(nodes_to_filter):
    nodes = stringify_ids(shift_ids(filter_fields(filter_nodes(read_nodes('data/nodes.json.gz'), nodes_to_filter))))

    with open('data/new_nodes.json', 'w') as f:
        json.dump(nodes, f, indent=4, ensure_ascii=False)
    return nodes
