# Agentic-Init (`ainit`)

> Interactive TUI Harness Engineering Tool for Project Initialization, Architecture Design, Commit/PR Conventions, and Automated Release Pipeline.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

---

## 📌 Table of Contents (목차)
- [🚀 Overview](#-overview)
- [💬 Slash Commands (슬래시 커맨드)](#-slash-commands-슬래시-커맨드)
- [🛠️ Makefile & Build Guide (메이크파일 빌드 및 설치 방법)](#️-makefile--build-guide-메이크파일-빌드-및-설치-방법)
- [📖 User Manual & Keybindings (사용 설명서 & 조작법)](#-user-manual--keybindings-사용-설명서--조작법)
- [📋 System Architecture Diagrams](#-system-architecture-diagrams)
- [📄 Documentation](#-documentation)
- [📜 License](#-license)

---

## 🚀 Overview

**Agentic-Init (`ainit`)**은 프로젝트의 **Initialization(초기 설계)**부터 **Release(상용 배포)**까지의 전체 개발 생애주기를 관장하는 CLI/TUI 하네스 엔지니어링 도구입니다.

- **바이너리 커맨드 명**: `ainit`

---

## 💬 Slash Commands (슬래시 커맨드)

TUI 대화창/프롬프트 모드에서 `/`를 입력하면 드롭다운 오토컴플릿(Dropdown Autocomplete) 팝업이 노출되며, 다음 슬래시 커맨드를 실행할 수 있습니다:

| Slash Command | Usage / Description | Example |
| :--- | :--- | :--- |
| **`/git-init`** | 원격 레포 존재 시 `git clone` 후 `work-dir` 업데이트, 없으면 신규 레포 초기화 | `/git-init myeong-han/ainit` |
| **`/set-confs`** | TUI 설정을 CLI 인자 형태로 즉시 업데이트 | `/set-confs --provider openai --arch monolith --git github` |
| **`/gen-docs`** | `docs/ARCHITECTURE_SPEC.md` 아키텍처 설계서 및 4종 Mermaid 차트 즉시 생성 | `/gen-docs` |
| **`/gen-codes`** | `AGENTS.md`, `CLAUDE.md`, `.cursorrules` 등 이종 에이전트 룰 파일 생성 | `/gen-codes` |
| **`/help`** | 이용 가능한 슬래시 커맨드 도움말 출력 | `/help` |

---

## 🛠️ Makefile & Build Guide (메이크파일 빌드 및 설치 방법)

```bash
# 바이너리 빌드 (bin/ainit 생성)
make build

# 빌드 및 TUI 실행
make run

# 전체 단위 테스트 실행
make test
```

---

## 📄 Documentation

- [Detailed User Manual (사용자 종합 설명서)](docs/USER_MANUAL.md)
- [Architecture Design Specification](docs/ARCHITECTURE.md)
- [TUI Questionnaire Form Specification](docs/QUESTIONNAIRE_SPEC.md)

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

> Maintained by **myeong-han**
