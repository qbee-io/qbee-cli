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
    ansible_inventory = {}
    for device in json_output['items']:
        tags = device['tags']
        if tags:
            for tag in tags:
                tag_label = re.sub(r'[ -]+', '_', tag)
                _add_device_to_group(device, tag_label, ansible_inventory)
    yaml.dump(ansible_inventory, sys.stdout)


# create ansible inventory from qbee.io devices by groups (hiararchical)
def qbee_devices_ansible_by_groups():

    json_output = json.loads(sys.stdin.read())
    ansible_inventory = {}

    for device in json_output['items']:
        if 'ancestors_titles' in device:
            groups = device['ancestors_titles']
        else:
            groups = device['ancestors']

        if len(groups) == 1:
            group_label = re.sub(r'[ -]+', '_', groups[0])
            _add_device_to_group(device, group_label, ansible_inventory)
            continue

        group_names = groups[:-1]
        dataref = ansible_inventory
        for group in group_names:
            group_label = re.sub(r'[ -]+', '_', group)
            # Last group
            if group == group_names[-1]:
                _add_device_to_group(device, group_label, dataref)
                break

            # Create group if not exists
            if group_label not in dataref:
                dataref[group_label] = {}
                dataref[group_label]['children'] = {}

            dataref = dataref[group_label]['children']

    yaml.dump(ansible_inventory, sys.stdout)


# add a device to a group
def _add_device_to_group(device, label, dataref):
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
