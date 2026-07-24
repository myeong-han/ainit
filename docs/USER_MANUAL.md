# Agentic-Init (`ainit`) Detailed User Manual

## 📖 Introduction
**Agentic-Init (`ainit`)**은 TUI 인터페이스를 통해 프로젝트 초기 아키텍처 설계부터 TDD 기반 소스코드 작성, 커밋/PR 컨벤션 검증, 그리고 상용 배포 파이프라인까지의 전체 과정을 대화형으로 안내하는 하네스 도구입니다.

---

## ⌨️ Global Keybindings Reference

- `Tab` / `Right Arrow`: Advance to the next step.
- `Shift+Tab` / `Left Arrow`: Move to the previous step.
- `Up` / `Down` (`k` / `j`): Navigate form items within the current step.
- `Enter` / `Space`: Toggle options or trigger health checks.
- `q` / `Ctrl+C`: Quit the application gracefully.

---

## ⚙️ Step-by-Step Configuration Guide

### Step 0: AI Licensing & Model Management
- **Subscription (OAuth)**: Authenticate using OpenCode compatible OAuth device authorization flow.
- **Direct API Keys**: Input your Anthropic, OpenAI, or Google Gemini API key.
- **Model Selection**: Choose your primary model for code generation and fallback model for low-latency tasks.

### Step 1: Architecture Spec & Mermaid Diagrams
- **Project Name**: Set the subdomain-safe project identifier.
- **Architecture Style**: Select MSA (Microservices Architecture), Modular Monolith, or Event-Driven Architecture.
- **Diagram Generation**: Toggle automatic generation for Sequence Diagrams, Traffic Flow, GitOps Pipelines, and ERDs.

### Step 2: MCP Tooling & Integrations
- **Git Provider**: Configure GitHub or Bitbucket organization and repository parameters.
- **Kubernetes**: Test connection to local `kubectl` contexts or remote cluster configurations.
- **CI/CD**: Connect Jenkins CI pipelines and ArgoCD GitOps applications.

### Step 3: Git Initialization, Commit/PR Conventions & TDD
- **Commit Convention**: Choose between Conventional Commits, Gitmoji, Issue Prefix (`PROJ-123`), or custom regex patterns.
- **PR Template**: Generate standard GitHub Pull Request templates with test & security checklists.
- **TDD Mode**: Enforce unit test-first development loop and local sandbox build verification before git commits.

### Step 4: Release Pipeline & Deployment
- **SemVer Tagging**: Automatic Major.Minor.Patch tag calculation.
- **Release Notes Sync**: Automatically publish change logs to Notion and Confluence.
- **Deployment Alerts**: Send real-time notification webhooks to Slack and Discord.
