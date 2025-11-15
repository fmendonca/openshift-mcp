# openshfit-mcp

MCP server em Go para interagir com todas as APIs do OpenShift (incluindo CRDs e Operators) e também com recursos do KubeVirt.

## Build

go mod tidy
go build -o openshfit-mcp ./cmd/openshift-mcp


## Uso (fora do cluster)

export KUBECONFIG=~/.kube/config
./openshfit-mcp



## Uso (dentro do OpenShift)

- Crie um ServiceAccount com permissões (RBAC) para os recursos desejados.
- Monte o ServiceAccount no Pod.
- Defina `IN_CLUSTER=true` no deployment.

O servidor fala MCP via stdio, podendo ser integrado a clientes MCP (Claude, etc.) ou exposto via HTTP com `mcp-go/server.ServeHTTP` se preferir.



## Uso (dentro do OpenShift)

- Crie um ServiceAccount com permissões (RBAC) para os recursos desejados.
- Monte o ServiceAccount no Pod.
- Defina `IN_CLUSTER=true` no deployment.

O servidor fala MCP via stdio, podendo ser integrado a clientes MCP (Claude, etc.) ou exposto via HTTP com `mcp-go/server.ServeHTTP` se preferir.

