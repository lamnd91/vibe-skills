---
name: review-code-android
description: Review Android and Kotlin code for correctness, lifecycle safety, threading, memory, and platform best practices. Use when reviewing Android or Kotlin code, pull requests in Android projects, or when the user asks for an Android code review.
---

# Android Code Review

Review Android/Kotlin code by following the checklist below.

## When to Apply

- User asks for an Android or Kotlin code review
- Reviewing pull requests or changes in Android modules
- Reviewing UI, ViewModels, repositories, or platform code

## General Framework

Follow the checklist: apply the sections that match the changes under review (UI, ViewModel, repository, platform, etc.). Use the full checklist by area in [references/checklists.md](references/checklists.md).

## Response Format

For each significant finding use:

1. **Observation**: What you noticed (with file/line or snippet if helpful)
2. **Explanation**: Why it matters on Android
3. **Recommendation**: Concrete fix or pattern (e.g. use `viewModelScope`, move to repository)
4. **Priority**: Critical / Moderate / Minor

Summarize at the end: 2–3 strengths and 2–3 main improvements. Keep tone constructive and educational; tie feedback to Android/Kotlin docs or project rules where relevant.

## Project Conventions

- Prefer Kotlin; follow existing style (indentation, naming) in the repo
- Respect project rules under `.cursor/rules/` (e.g. RecyclerView adapters, event handling, database fields)
- Do not suggest adding libraries without explicit user approval (per general rule)
