# Agentic-Init (`ainit`)

> Interactive TUI Harness Engineering Tool for Project Initialization, Architecture Design, Commit/PR Conventions, and Automated Release Pipeline.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

---

## 🚀 Overview

**Agentic-Init (`ainit`)**은 프로젝트의 **Initialization(초기 설계)**부터 **Release(상용 배포)**까지의 전체 개발 생애주기를 관장하는 CLI/TUI 하네스 엔지니어링 도구입니다.

- **바이너리 커맨드 명**: `ainit`

---

## 🛠️ Build & Installation Guide (빌드 및 설치 방법)

### Prerequisites (사전 요구사항)
- **Go**: `1.21` 버전 이상 설치 필요 ([Go 설치 가이드](https://go.dev/doc/install))

### 1. Repository Clone (저장소 복사)
```bash
git clone https://github.com/myeong-han/ainit.git
cd ainit
```

### 2. Dependency Tidy (의존성 동기화)
```bash
go mod tidy
```

### 3. Build Binary (바이너리 빌드)
```bash
# bin/ainit 바이너리 빌드
go build -o bin/ainit ./cmd/ainit
```

### 4. Run TUI (TUI 도구 실행)
```bash
# 빌드된 바이너리 직접 실행
./bin/ainit

# 또는 go run을 통해 소스에서 직접 실행
go run ./cmd/ainit
```

### 5. Install to System Path (시스템 글로벌 설치 - 선택사항)
```bash
# $GOPATH/bin 위치에 ainit 커맨드로 등록
go install ./cmd/ainit

# 이제 어디서든 ainit 실행 가능
ainit
```

### 6. Run Unit Tests (단위 테스트 실행)
```bash
go test -v ./...
```

---

## 📋 System Architecture Diagrams

### 1. Overall 5-Step Pipeline Overview
```mermaid
flowchart LR
    Start([ainit CLI Launch]) --> S0[Step 0: AI Licensing]
    S0 --> S1[Step 1: Architecture Spec]
    S1 --> S2[Step 2: MCP Integrations]
    S2 --> S3[Step 3: Harness Coding]
    S3 --> S4[Step 4: Release Pipeline]
    S4 --> End([Project Released 🎉])
```

### 2. Step 0 & 1: AI Licensing & Architecture Spec Flow
```mermaid
flowchart TD
    subgraph Step0 [Step 0: AI Licensing]
        A[Launch ainit] --> B{Auth Method}
        B -->|OAuth| C[Subscription Device Flow]
        B -->|API Key| D[Direct Provider Key]
        C --> E[LLM Provider & Fallback Setup]
        D --> E
    end

    subgraph Step1 [Step 1: Architecture Spec Generation]
        E --> F[Select Arch Style: MSA / Monolith / EDA]
        F --> G[Generate Mermaid Diagrams]
        G --> G1[Sequence Diagram]
        G --> G2[Traffic & Ingress Flow]
        G --> G3[GitOps Pipeline Flow]
        F --> H[Define API Contract & ERD Spec]
        F --> I[Choose Layout: Monorepo vs Multirepo]
    end
```

### 3. Step 2: MCP & Infrastructure Tooling Topology
```mermaid
flowchart TD
    TUI[ainit TUI Engine] -->|Required| Git[Git Provider: GitHub / Bitbucket]
    
    subgraph Infrastructure [K8s & CI/CD Integrations]
        TUI -->|Health Check| K8s[Kubernetes Cluster]
        TUI -->|Trigger| Jenkins[Jenkins CI]
        TUI -->|Sync| ArgoCD[ArgoCD CD]
        TUI -->|Auth & Push| Registry[Container Registry: Harbor/GHCR/ECR]
    end

    subgraph Operations [Docs & Notifications]
        TUI -->|Update Docs| Notion[Notion API]
        TUI -->|Update Docs| Confluence[Confluence API]
        TUI -->|Alert Webhook| Slack[Slack Channel]
        TUI -->|Alert Webhook| Discord[Discord Channel]
    end
```

### 4. Step 3 & 4: Harness TDD Loop & Release Pipeline
```mermaid
flowchart TD
    subgraph Step3 [Step 3: Harness TDD & Micro-Commit]
        Config[Setup AGENTS.md, Commit/PR Conventions] --> WriteTest[1. Write Failing Test Spec]
        WriteTest --> WriteCode[2. Write Minimal Code]
        WriteCode --> Sandbox{3. Local Sandbox Build & Test}
        Sandbox -->|Fail| WriteCode
        Sandbox -->|Pass| HookCheck{4. Verify Commit Convention}
        HookCheck -->|Pass| Commit[5. Atomic Git Micro-Commit]
    end

    subgraph Step4 [Step 4: Release Pipeline]
        Commit --> SemVer[Auto SemVer Tagging]
        SemVer --> DocsSync[Sync Release Notes to Notion/Confluence]
        DocsSync --> Notify[Send Slack/Discord Release Alert]
    end
```

---

## 📄 Documentation

- [Architecture Design Specification](docs/ARCHITECTURE.md)
- [TUI Questionnaire Form Specification](docs/QUESTIONNAIRE_SPEC.md)

---

> Maintained by **myeong-han**
