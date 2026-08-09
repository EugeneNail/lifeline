---
name: make-commits
description: Plan and create semantic git commits from the current worktree. Use when the user asks to make commits, split changes into commit groups, review unstaged or staged git changes, or order dependent changes before dependent commits.
---

# Make Commits

## Workflow

1. Inspect the current worktree with `git status --short`, `git diff --stat`, and `git diff`.
2. Pull context from recent history with `git log --oneline -n 10`.
3. Group changes by meaning, not by file count.
4. Keep unrelated edits out of the current commit set.
5. Order commits by dependency. Commit prerequisites first, then commit the files that rely on them.
6. Create one commit per semantic group.
7. Use Conventional Commits format: `type(scope): short message`.
8. Write the message in plain conversational English. Keep it direct and specific. Do not optimize for git folklore about imperative mood or 72-character limits.

## Commit Selection Rules

- Prefer the smallest commit that still leaves a coherent change.
- If a decomposition has already been made for the current commit run, treat it as binding: create commits strictly according to that split, without regrouping, merging, or inventing extra commits.
- If one change only makes sense after another, commit the dependency first.
- Do not mix unrelated cleanup with feature work.
- Leave stray or user-owned changes untouched unless the user explicitly asked to include them.

## Output Expectations

- Show the proposed commit split before committing when the split is not obvious.
- Create the commits directly once the split is clear.
- Report the commit hashes and messages after finishing.
