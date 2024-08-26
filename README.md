# Nagios-ES

## Description

This is a Nagios plugin that checks the status of Elasticsearch cluster.

## Build

```
go build -o check_es .
```

## Usage

```
check_es --es_url=http://<domain>:<port> --check=<check_name> --w=50 --c=90
```

### Options

- `es_url`: Elasticsearch URL (format: http://<domain>:<port>)
- `check`: Check name
  - `health`: Check cluster health
  - `node_count`: Overall nodes count (W/C required)
  - `data_node_count`: Check data nodes count (W/C required)
  - `cpu_usage`: Check node CPU usage (W/C required)
    - If filter will be used, it will check the CPU usage of the specific node
    - If filter is not used, it will check the maximum CPU usage of all nodes
  - `heap_size`: Check node Head usage (W/C required)
    - If filter will be used, it will check the Heap usage of the specific node
    - If filter is not used, it will check the maximum Heap usage of all nodes
  - `disk_usage`: Check node Disk usage (W/C required)
    - If filter will be used, it will check the Disk usage of the specific node
    - If filter is not used, it will check the maximum Disk usage of all nodes
- `w`: Warning threshold
- `c`: Critical threshold
For filtering node specific checks, you can use the following options:
- `node_ip`: Node IP address for node specific checks
- `node_name`: Node name for node specific checks

## Example

```
./check_es --es_url=http://<domain>:<port> --check=health
```

OK Output:
```
OK: Cluster health is green | 'time'=19ms;;;;
```

Warning Output:
```
WARNING: Cluster health is yellow, relocating shards: 0 | 'active_shards'=5;;;; 'relocating_shards'=0;;;; 'time'=24ms;;;; 'unassigned_shards'=3;;;;
```
