# Collection of qbee-cli

## qbee-devices-ansible.py

Utility to generate ansible inventory files based on qbee-cli listing

### Usage:

```
qbee-cli devices list --json --limit <n> | python3 qbee-devices-ansible.py > my_ansible_inventory.yml
```

### Inventory by groups (hiearchical)

```
qbee-cli devices list --json --limit <n> | python3 qbee-devices-ansible.py --by-groups > my_ansible_inventory.yml
```

### Inventory by tags (non-hierarchical)
```
qbee-cli devices list --json --limit <n> | python3 qbee-devices-ansible.py --by-tags > my_ansible_inventory.yml
```