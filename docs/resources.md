# Resources do OpenShift/Kubernetes MCP Server

Além dos tools, o servidor expõe resources somente‑leitura via URIs.

## Cluster

- `cluster://info`
  - Retorna um JSON com informações gerais do cluster:
    - Versão do Kubernetes (`gitVersion`, `major`, `minor`)
    - Plataforma, build date, Go version
    - Quantidade de API groups disponíveis

- `cluster://openshift/version`
  - Em clusters OpenShift, retorna:
    - `desiredVersion`
    - `channel`
    - Histórico de updates (`history`)
  - Em clusters puros Kubernetes, indica que OpenShift não está disponível.

- `cluster://apigroups`
  - Lista todos os API groups do cluster, com:
    - `name`
    - `versions`
    - `preferredVersion`

## Namespaces

- `namespaces://all`
  - Lista todos os namespaces com:
    - `name`
    - `phase`
    - `labels`

- `namespaces://detail?name=<NAMESPACE>`
  - Detalhes de um namespace específico:
    - `name`
    - `phase`
    - `labels`
    - `annotations`
  - Se `name` não for informado, retorna um JSON de erro com mensagem e usage.
