# Tools do OpenShift/Kubernetes MCP Server

Este MCP server expõe uma série de tools organizados por recurso do cluster.

## Pods

- `list_pods`
  - Lista pods em um namespace ou em todos os namespaces.
  - Parâmetros:
    - `namespace` (string, opcional)
    - `labelSelector` (string, opcional)

- `get_pod`
  - Detalhes de um pod específico.
  - Parâmetros:
    - `name` (string)
    - `namespace` (string)

- `get_pod_logs`
  - Obtém logs de um container em um pod.
  - Parâmetros:
    - `name` (string)
    - `namespace` (string)
    - `container` (string, opcional)
    - `tailLines` (int, padrão 100)
    - `previous` (bool, padrão false)

- `delete_pod`
  - Remove um pod específico.
  - Parâmetros:
    - `name` (string)
    - `namespace` (string)

- `exec_pod`
  - Executa um comando dentro de um container do pod.
  - Parâmetros:
    - `name` (string)
    - `namespace` (string)
    - `container` (string, opcional)
    - `command` ([]string)

## Deployments

- `list_deployments`
- `get_deployment`
- `scale_deployment`
- `restart_deployment`

## Services

- `list_services`
- `get_service`

## Routes (OpenShift)

- `list_routes`
- `get_route`

## ImageStreams (OpenShift)

- `list_imagestreams`
- `get_imagestream`

## Projects / Namespaces (OpenShift)

- `list_projects`
- `get_project`

## Nodes

- `list_nodes`
- `get_node`

## ConfigMaps

- `list_configmaps`
- `get_configmap`

## Secrets

- `list_secrets`
- `get_secret` (dados mascarados nas saídas)

## PVCs

- `list_pvcs`
- `get_pvc`

## Ingress

- `list_ingresses`
- `get_ingress`

## RBAC

- `list_roles`
- `list_rolebindings`
- `list_clusterroles`
- `list_clusterrolebindings`
