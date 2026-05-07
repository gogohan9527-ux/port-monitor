# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## 5. Session Handoff and Recovery

**Assume another agent may have worked before you. Resume from recorded state.**

At the start of a new session or takeover:
- Read `.claude/handoff.md` first. If it is missing, read `.claude/progress.log`.
- Check `MEMORY.md` when present for recurring project issues.
- Run `git status` and the project's quick test command when available.
- Compare the local state with the handoff notes before editing.
- Continue from `To-Do` or `Blocked` items. Do not rewrite modules marked `Completed`.

During multi-step or complex work:
- Use subagents when the task can be split into independent frontend, backend, QA, or research work.
- Keep changes idempotent. If prior work is partial, reconcile it before adding more.
- Do not guess missing critical context. Ask instead of inventing placeholders.
- Keep UI changes flat and white-background unless the project already defines a different visual system.
- Use consistent status tags when helpful: `[TODO-AGENT-A]`, `[BLOCKED]`, `[VERIFIED]`.

Before ending, compacting, or handing off:
- Update `.claude/handoff.md` with completed work, current state, blockers, changed files, verification run, and exact next steps.
- Keep `CLAUDE.md` focused on agent behavior. Put formatting and lint rules in Prettier, ESLint, hooks, or project config.
- If context is getting large, summarize progress into the handoff file and start a fresh session.
- If a resumed agent goes off track, prefer rewind plus clearer instructions over repeated correction in chat.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
