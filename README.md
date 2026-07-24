# Agentic-Init (`ainit`)

> Interactive TUI Harness Engineering Tool for Project Initialization, Architecture Design, Commit/PR Conventions, and Automated Release Pipeline.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## 🚀 Overview

**Agentic-Init (`ainit`)**은 프로젝트의 **Initialization(초기 설계)**부터 **Release(상용 배포)**까지의 전체 개발 생애주기를 관장하는 CLI/TUI 하네스 엔지니어링 도구입니다.

- **바이너리 커맨드 명**: `ainit`

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

## 📦 Installation & Usage

```bash
# Clone the repository
git clone https://github.com/myeong-han/ainit.git
cd ainit

# Run the TUI Harness Engine
ainit
```

---

## 📄 Documentation

- [Architecture Design Specification](docs/ARCHITECTURE.md)
- [TUI Questionnaire Form Specification](docs/QUESTIONNAIRE_SPEC.md)

---

> Maintained by **myeong-han**
