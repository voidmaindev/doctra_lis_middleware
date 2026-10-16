# Tech Stack

- React + TypeScript (frontend), Go (backend), PostgreSQL.
- Go and PostgreSQL are Docker-containerized. React runs locally for fast iteration without rebuilding containers on every change.

# Go Preferences

Always use these libraries — no alternatives unless explicitly told otherwise:

- **ORM:** GORM (`gorm.io/gorm`)
- **HTTP framework:** Fiber v2 (`github.com/gofiber/fiber/v2`)
- **Logging:** slog (stdlib `log/slog`)
- **HTTP client:** req (`github.com/imroc/req/v3`)
- **Authorization:** Casbin v3 (`github.com/casbin/casbin/v3`)
- **CLI:** Cobra (`github.com/spf13/cobra`)
- **Configuration:** Viper (`github.com/spf13/viper`)
- **Migrations:** GORM AutoMigrate
- **Architecture:** Domain-Driven Design (DDD)
- **Project layout:** `/cmd`, `/internal`, `/pkg`

# React Preferences

- **Bundler:** Vite
- **CSS:** Tailwind CSS

# Code Quality

- Always handle errors explicitly — never ignore or silently swallow them.
- No magic strings or numbers — use constants.
- Keep functions under 40 lines. If longer, split.
- Write table-driven tests in Go.
- Name things clearly — no single-letter variables outside loops.
- One function, one job. If it does two things, split it.

# Rules

- Think first, read relevant files before answering. Never speculate about code you haven't opened.
- If a file is referenced, read it before responding. Investigate before answering — no hallucinations.
- Give a high-level explanation of changes at every step.
- Commit and push whenever you see fit — no need to wait for approval.

## CRITICAL: Plan mode file ≠ GitHub Kanban board

- The Claude plan mode file (`.claude/plans/...`) is NOT the same as the GitHub Projects Kanban board.
- Approving the plan mode file does NOT count as approving the project plan.
- Both must exist independently. The plan mode file is temporary; the GitHub Kanban board is permanent and persistent across sessions.

# Kanban Board — MANDATORY EVERY SESSION

**NO CODE until the board is reviewed and up to date. This applies EVERY session, not just kickoff.**

## Session start protocol — do this BEFORE ANY CODE, EVERY TIME:

1. **Read the board state** using `gh project item-list` and `gh issue list` to see all issues and their status.
   - If the GitHub Project board doesn't exist, STOP and run the Kickoff Interview above. No exceptions.
2. **Review** what's done (issues in **Done** column) and what's next (issues in **Todo** column).
3. **Update the board** if needed:
   - Create new issues for tasks discovered since last session and add them to the board.
   - Break down upcoming tasks into sub-issues if they're too big.
   - Close or remove issues that are no longer relevant.
   - If starting a new phase, flesh out its details (field names, validation rules, exact behavior) in the phase issue body before coding it.
4. **Move cards** — shift the next batch of work from **Backlog** to **Todo**.
5. **Tell the user** what's done, what's next, and what you updated on the board.
6. **HARD GATE: Only after steps 1-5 are complete, begin coding.**

## While coding:

- Move issues to **In Progress** when you start working on them.
- Move issues to **Done** and close them immediately after completing them.
- If a task turns out to be bigger than expected, create sub-issues and add them to the board right away.
- If you **or the user** identify new work (new features, bug reports, follow-ups, refactors), create issues on the board BEFORE writing any code for it. See the "New Work Request Protocol" section below.
- Never leave the board stale — it must reflect reality at all times.
- Never commit code without first verifying that the board reflects the current state.

## Sub-agent checkpoint rule:

- **After EVERY sub-agent completes work, you MUST immediately update the Kanban board** BEFORE doing anything else (committing, starting the next phase, launching another sub-agent, etc.).
- Sub-agents cannot update the board. It is YOUR responsibility to move cards and close issues after a sub-agent returns.
- If you batch multiple phases into one sub-agent call, update ALL completed issues when it returns — not later.
- When updating after a sub-agent, don't just close issues — also **create new sub-issues** (up to 3 levels deep) for any work the sub-agent completed that wasn't already broken down, and create sub-issues for the next upcoming tasks.

# New Work Request Protocol — MANDATORY FOR EVERY REQUEST

**This is the rule the kickoff section does NOT cover. Read it carefully.**

Whenever the user asks for new work — at ANY point in the project, not just kickoff, not just session start — you MUST update the Kanban board BEFORE writing any code. This protocol applies EVERY time the user requests new work. Do NOT assume the session-start protocol covers mid-session requests. Do NOT assume the kickoff board makes this protocol optional.

## Trigger — BROAD on purpose:

This protocol fires on **any** user request for new work. That includes:

- Adding features
- Changing or extending existing features
- Bug fixes
- Refactors
- Tweaks and small improvements
- Follow-ups from previous work

Trigger phrases include (non-exhaustive): "add", "implement", "build", "create", "change", "fix", "refactor", "extend", "now do", "let's do", "can you", "make it", "update", "improve".

**If in doubt whether a request counts as new work — it does. Run the protocol.**

## HARD GATE — NO CODE until issues exist on the board:

**NO CODE may be written until issues for the new work exist on the Kanban board.**

Unlike kickoff, you do NOT need to wait for user approval of the hierarchy. The flow is:

1. Create a top-level issue for the work. Use a label to mark it (`feature`, `bug`, `refactor`, etc.).
2. Break it into sub-issues:
   - **Large features:** 3 levels of depth, matching the kickoff hierarchy rules in the "Issue structure — ALWAYS use sub-issues" section below.
   - **Bug fixes and small tweaks:** 1–2 levels is acceptable.
3. Add the issues to the board. Put the issue you're starting on in **In Progress**; put queued ones in **Todo**.
4. Tell the user in ONE line what you added to the board (e.g., "Added issue #42 'Dark mode toggle' with 3 sub-issues to the board, starting #43 now.").
5. Proceed to code.

If you catch yourself opening an editor or writing code before the issues exist, STOP and create the issues first.

## Issue structure — ALWAYS use sub-issues:

**Every task MUST have sub-issues. No flat boards. Aim for 3 levels of depth.**

- When creating the initial board, break every task into sub-issues immediately.
- When starting a task, break its sub-issues into sub-sub-issues if they're non-trivial.
- When finishing a task or sub-issue, close it AND create/update sub-issues for the next one.

Example structure:

```
Phase 1: Project Setup (issue #1, label: phase)
├── Initialize Go backend (issue #2, sub-issue of #1)
│   ├── Create go.mod with module name (issue #5, sub-issue of #2)
│   ├── Set up /cmd, /internal, /pkg folder structure (issue #6, sub-issue of #2)
│   └── Create main.go entry point (issue #7, sub-issue of #2)
├── Docker Compose for Postgres + Go (issue #3, sub-issue of #1)
│   ├── Postgres container (issue #8, sub-issue of #3)
│   │   ├── Configure volume for data persistence (issue #11, sub-issue of #8)
│   │   └── Set environment variables (issue #12, sub-issue of #8)
│   └── Go container with hot reload (issue #9, sub-issue of #3)
│       ├── Write Dockerfile with multi-stage build (issue #13, sub-issue of #9)
│       ├── Add Air config for hot reload (issue #14, sub-issue of #9)
│       └── Wire up to Docker Compose network (issue #15, sub-issue of #9)
└── Database migrations (issue #4, sub-issue of #1)
    ├── Create initial migration for users table (issue #10, sub-issue of #4)
    └── Add GORM AutoMigrate setup (issue #16, sub-issue of #4)
```

# Documentation — MANDATORY

1. **Check if `docs/` folder exists** in the project root.
   - If it exists: read relevant docs before making changes.
   - If it does NOT exist: create `docs/` and populate it alongside your first feature work.

2. **Required docs:**
   - `docs/architecture.md` — how the app works inside and out. Update this whenever you add or change a system component.
   - Additional docs per feature or domain as needed (e.g., `docs/auth.md`, `docs/api.md`).

3. **Keep docs in sync** — when you add or change a feature, update the relevant doc in the same session. Don't defer it.

# Testing

- Write unit tests for all Go business logic (table-driven).
- Write integration tests for API endpoints using Go's httptest.
- Frontend: use Vitest + React Testing Library for component tests.
- After medium+ features, test in browser using Playwright — logins and main user paths.
- Tests live next to the code they test, not in a separate folder.

# Docker

- Always rebuild containers with `--no-cache` when changing anything container-based.

# Tools

- For any frontend work, always use the `frontend-design` skill.
- Always use Context7 MCP whenever needed.
- Always use Playwright MCP for browser testing.
