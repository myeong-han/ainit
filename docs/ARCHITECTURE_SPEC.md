# ainit-project Architecture Specification

## 1. Overview & Requirements
**AI Provider**: anthropic (subscription)
**Primary Model**: claude-3-5-sonnet-20241022
**Architecture Style**: MSA
**Repository Structure**: monorepo

### Plain Text Architecture Prompt Requirements
> Received architecture prompt: '아하'. Press 'Ctrl+S' to generate docs & code.

---

## 2. Generated Mermaid Diagrams

### 2.1 Microservices & Middleware Overview
```mermaid
flowchart LR
    User([Client]) --> Gateway[API Gateway]
    Gateway --> Auth[Auth Service]
    Gateway --> Core[Core App Domain]
    Core --> DB[(Database)]
    Core --> Kafka[Event Broker]
```

### 2.2 Sequence Diagram
```mermaid
sequenceDiagram
    autonumber
    Client->>Gateway: HTTP Request
    Gateway->>Auth: Validate JWT Token
    Auth-->>Gateway: 200 OK
    Gateway->>Service: Process Business Logic
    Service-->>Client: HTTP 200 Success
```

### 2.3 Ingress Traffic Flow
```mermaid
flowchart TD
    Internet --> Ingress[Nginx Ingress / Emissary]
    Ingress --> ServiceMesh[Istio Service Mesh]
    ServiceMesh --> PodA[App Pod 1]
    ServiceMesh --> PodB[App Pod 2]
```

### 2.4 GitOps & CI/CD Pipeline Flow
```mermaid
flowchart LR
    GitPush([Git Commit Push]) --> Jenkins[Jenkins CI Test]
    Jenkins --> Harbor[Container Registry Push]
    Harbor --> ArgoCD[ArgoCD GitOps Sync]
    ArgoCD --> K8s[Kubernetes Cluster]
```

---

## 3. Harness Code Strategy
- **TDD Mode**: true
- **Commit Convention**: conventional
- **Local Sandbox Check**: true
