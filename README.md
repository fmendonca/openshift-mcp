# OpenShift/Kubernetes MCP Server

Model Context Protocol (MCP) server for interacting with OpenShift and Kubernetes clusters.

## Features

- **Pods**: List, get, logs, exec, delete
- **Deployments**: List, get, scale, restart
- **Services**: List, get details
- **Routes**: List and inspect OpenShift routes
- **ImageStreams**: Manage OpenShift image streams
- **Projects**: List and manage OpenShift projects
- **Nodes**: Cluster node information
- **ConfigMaps**: Configuration management
- **Secrets**: Secret management (masked)
- **PVCs**: Persistent volume claims
- **Ingress**: Ingress resources
- **RBAC**: Roles and bindings

## Installation



## Example to Run
```bash
MCP_TRANSPORT=stdio ./openshift-mcp
```
```bash
MCP_TRANSPORT=http MCP_HTTP_ADDR=":8080" ./openshift-mcp
```

## Variable to run

 - **MCP_TRANSPORT=stdio|http**

### Option stdio run local with Agent IA
### Option http run on cluster and receive instruction by api
