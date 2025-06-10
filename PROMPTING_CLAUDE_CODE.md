# Effective Prompting for Claude Code

## The Problem: Literal Implementation

Claude Code tends to be extremely literal when implementing features. If it sees a TODO comment or a placeholder, it will often preserve it rather than implementing the actual functionality. This can lead to incomplete implementations where UI elements exist but don't work.

## Example of What Went Wrong

**Given:** A codebase with UI buttons for edit/delete and corresponding TODO comments
**Result:** Claude implemented the UI but left the functionality as TODO stubs
**Discovery:** Only found when user tried to use the features

## Better Prompting Strategies

### 1. Explicit Completeness Requirements

**Instead of:**
```
"read all the md files and start working"
```

**Use:**
```
"Implement ALL features shown in the UI - if there's a button, it must work. No TODO stubs."
```

### 2. Upfront Feature Checklist

**Prompt:**
```
"First, list all features mentioned in the docs/UI that need implementation, then implement them ALL"
```

This forces Claude to acknowledge incomplete features before starting.

### 3. Test-First Approach

**Prompt:**
```
"For each UI element, write a test that verifies it actually works, not just displays"
```

Would catch non-functional buttons immediately.

### 4. Incremental Verification

**Prompt:**
```
"After implementing each feature, show me a test that proves it works"
```

Rather than implementing everything then discovering gaps.

### 5. Anti-TODO Directive

**Prompt:**
```
"Replace ALL TODO comments with working implementations"
```

Claude tends to preserve TODOs unless explicitly told otherwise.

### 6. Feature Completion Check

**Prompt:**
```
"Before committing, verify every button/command in the UI actually does something"
```

## The Optimal Initial Prompt

```
Read all the md files and implement a FULLY FUNCTIONAL budget tracker where:
- EVERY feature mentioned in the UI actually works
- No TODO stubs or placeholder implementations
- Test each feature after implementing to prove it works
- If you see a button or command, implement its functionality
- Commit every step
```

## Key Insights

1. **Claude is overly literal** - It implements exactly what's asked, nothing more
2. **Claude preserves existing patterns** - If it sees TODOs, it keeps them
3. **Claude doesn't infer intent** - Won't think "edit/delete are core CRUD features that obviously need to work"
4. **Claude needs explicit completeness criteria** - Tell it when something should be "complete"

## How to Think About Claude Code

Prompt Claude like it's a very capable but extremely literal junior developer who will:
- Do exactly what's asked
- Not make assumptions about what "should" work
- Preserve existing code patterns (including TODOs)
- Need explicit instructions for completeness

## Red Flags to Watch For

When reviewing Claude's work, look for:
- TODO comments that weren't replaced
- UI elements without backend implementation
- Features mentioned in docs but not implemented
- Placeholder or example data instead of real functionality

## Verification Prompts

After implementation, use prompts like:
```
"List all UI elements and confirm each one has working functionality"
"Show me that the edit button actually edits data in the database"
"Prove that delete actually removes records"
```

This document represents lessons learned from implementing partial functionality when full functionality was expected.