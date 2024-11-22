# Collection of qbee-cli tools

## qbee-devices-ansible.py

Utility to generate ansible inventory files based on qbee-cli listing. The ansible inventory file will contain
`ProxyCommand` for `ssh` that will use the `qbee-cli` portforwarding functionality to connect to the hosts. 
Group and tag names will be converted to ansible compatible group names. You will need to have a valid `qbee-cli`
session setup through `qbee-cli login` in order for this work. Also, you will need to have a valid ssh public key
distributed to the target devices (which can be done by qbee-agent SSH key distribution). Devices will be identified
with their internal qbee id.

### Inventory by groups (hiearchical)
```
qbee-cli devices list --json --limit <n> | python3 qbee-devices-ansible.py --by-groups > ansible_inventory.yml
```

### Inventory by tags (non-hierarchical)
```
qbee-cli devices list --json --limit <n> | python3 qbee-devices-ansible.py --by-tags > ansible_inventory.yml
```

### Example output:
```
example_tag_1:
  hosts:
    636f6e04f254592cc99a7324f2ba07f2abe3755d5f5c4fd941379c6e5094afd5:
      ansible_ssh_common_args: -o StrictHostKeyChecking=no -o ProxyCommand="qbee-cli
        connect -d 636f6e04f254592cc99a7324f2ba07f2abe3755d5f5c4fd941379c6e5094afd5
        -t stdio:localhost:22"
example_tag_2:
  hosts:
    d9e6887a16f5c2a1f6e8e3a0970368d9de34f5e8685b19b261dbc438d0649dec:
      ansible_ssh_common_args: -o StrictHostKeyChecking=no -o ProxyCommand="qbee-cli
        connect -d d9e6887a16f5c2a1f6e8e3a0970368d9de34f5e8685b19b261dbc438d0649dec
        -t stdio:localhost:22"
```

### Example: Running a test playbook on RaspberryPI

``` ansible_playbook.yml
---
- name: 
  gather_facts: false
  hosts: example_tag_1
  tasks:
    - name: get motd
      ansible.builtin.command: cat /etc/motd
      register: mymotd
    
    - name: debug motd
      debug:
        msg: "This is my motd: {{ mymotd }}"

```

```
ansible-playbook -i ansible_inventory.yml --user pi ansible_playbook.yml
```