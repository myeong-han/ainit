# Agentic-Init (`ainit`) TUI Questionnaire & Form Specification

## TUI Form & Field Mapping

| Step | Form Item | Component Type | Options / Validation |
| :--- | :--- | :--- | :--- |
| **Step 0** | Licensing Mode | Single Select | Subscription (OAuth) / Direct API Key / Local LLM |
| **Step 0** | Primary Model | Single Select | Claude 3.5 Sonnet / GPT-4o / Gemini 1.5 Pro |
| **Step 1** | Project Name | Input Text | Subdomain-safe string (e.g. `my-awesome-msa`) |
| **Step 1** | Architecture | Single Select | Microservices (MSA) / Modular Monolith / EDA |
| **Step 1** | Repo Structure | Single Select | Monorepo (pnpm/go) / Multirepo |
| **Step 2** | Git Provider | Single Select | GitHub / Bitbucket (Required) |
| **Step 2** | K8s Target | Single Select | Local Context / Remote Kubeconfig / None |
| **Step 2** | CI/CD Options | Multi Select | Jenkins CI / ArgoCD CD / Container Registry |
| **Step 2** | Notification | Multi Select | Slack Webhook / Discord Webhook / Notion |
| **Step 3** | **Commit Convention** | Single Select | Conventional Commits / Gitmoji / Issue Prefix / Custom Regex |
| **Step 3** | **Issue Key Format** | Input Text | Regex pattern (e.g. `PROJ-\d+` for Jira/GitHub Issues) |
| **Step 3** | **PR Template** | Single Select | Standard Feature Checklist / Minimal / Jira Integrated / Custom |
| **Step 3** | **Auto PR Labeling** | Toggle | Auto-assign labels based on commit types (Default: ON) |
| **Step 3** | Code Strategy | Single Select | TDD First / Feature Driven / Spec First |
| **Step 3** | Local Sandbox | Toggle | Run test & build before commit (Default: ON) |
