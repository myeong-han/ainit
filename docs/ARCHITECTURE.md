# Agentic-Init (`ainit`) Architecture Specification

## Overview
**Agentic-Init (`ainit`)** is an interactive TUI harness engineering tool that manages the entire lifecycle of a software project from architecture design to release.

## System Architecture Diagrams

### 1. Overall Pipeline Flow
```mermaid
flowchart LR
    Start([ainit CLI Launch]) --> S0[Step 0: AI Licensing]
    S0 --> S1[Step 1: Architecture Spec]
    S1 --> S2[Step 2: MCP Integrations]
    S2 --> S3[Step 3: Harness Coding]
    S3 --> S4[Step 4: Release Pipeline]
    S4 --> End([Project Released 🎉])
```

### 2. Step 0 & 1: AI Licensing & Architecture Design Flow
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

### 3. Step 2: MCP Integrations Topology
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

### 4. Step 3 & 4: Harness TDD Sandbox & Release Pipeline
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
