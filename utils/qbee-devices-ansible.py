#!/usr/bin/env python

import argparse
import json
import sys
import yaml
import re

PROXY_COMMAND = (
    '-o StrictHostKeyChecking=no '
    '-o ProxyCommand="qbee-cli connect -d {} -t stdio:localhost:22"'
)


# create ansible inventory from qbee.io devices by tags (non-hiararchical)
def qbee_devices_ansible_by_tags():
    json_output = json.loads(sys.stdin.read())
    by_tags = {}
    for device in json_output['items']:
        tags = device['tags']
        if tags:
            for tag in tags:
                _add_device_to_group(device, tag, by_tags)
    yaml.dump(by_tags, sys.stdout)


# create ansible inventory from qbee.io devices by groups (hiararchical)
def qbee_devices_ansible_by_groups():

    json_output = json.loads(sys.stdin.read())
    by_groups = {}

    for device in json_output['items']:
        if 'ancestors_titles' in device:
            groups = device['ancestors_titles']
        else:
            groups = device['ancestors']

        if len(groups) == 1:
            _add_device_to_group(device, groups[0], by_groups)
            continue

        group_names = groups[:-1]
        dataref = by_groups
        for group in group_names:
            group = re.sub(r'[ -]+', '_', group)
            # Last group
            if group == group_names[-1]:
                _add_device_to_group(device, group, dataref)
                continue

            # Create group if not exists
            if group not in dataref:
                dataref[group] = {}
                dataref[group]['children'] = {}

            dataref = dataref[group]['children']

    yaml.dump(by_groups, sys.stdout)


# add a device to a group
def _add_device_to_group(device, label, dataref):
    label = re.sub(r'[ -]+', '_', label)
    if label not in dataref:
        dataref[label] = {}

    if 'hosts' not in dataref[label]:
        dataref[label]['hosts'] = {}

    node_id = device['node_id']
    dataref[label]['hosts'][node_id] = {
        'ansible_ssh_common_args': PROXY_COMMAND.format(node_id)
    }


if __name__ == "__main__":

    parser = argparse.ArgumentParser(
                    prog='qbee-devices-ansible',
                    description='Convert qbee.io devices to ansible inventory')

    parser.add_argument(
        '--by-groups',
        action='store_true',
        help='Create ansible inventory by groups'
    )
    parser.add_argument(
        '--by-tags',
        action='store_true',
        help='Create ansible inventory by tags'
    )
    args = parser.parse_args()

    if args.by_groups:
        sys.exit(qbee_devices_ansible_by_groups())

    if args.by_tags:
        sys.exit(qbee_devices_ansible_by_tags())

    parser.print_help()
    sys.exit(1)
