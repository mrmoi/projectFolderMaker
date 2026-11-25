---
name: moises-pull-request-description
description: Write a pull request description for this repo. Use when opening or updating a PR, or when asked to draft, write, or fix a PR description.
---

# Pull request description

Fill in the repo's PR template if one exists (`.github/pull_request_template.md`). Otherwise use the layout below.

## Layout

```
## Summary
One or two sentences: what changed and why.

## Changes
- Bullet per notable change, in plain language.

## Testing
- How it was verified (commands run, sample inputs, results).
- Anything not tested, stated plainly.

## Notes
- Behavior changes, follow-ups, or risks reviewers should know about. Omit if none.
```

## Rules

- Describe the diff (`git diff main...HEAD`), not the conversation.
- Lead with the reason for the change, not a file list.
- Be factual: do not claim tests passed unless they were run.
- Keep it short; no filler or marketing language.
- Do not include secrets, internal hostnames, or personal details.
- Reference issues and PRs with full links or `owner/repo#123`, never a bare `#123`.
