# DotAgents

Agent configuration templates for OurSky/Vendatta.

## Structure

```
.oursky/
└── templates/
    ├── commands/    # Agent command specifications
    └── rules/       # Coding standards and guidelines
```

## Commands

- `/plan` - Generate dependency-ordered tasks
- `/implement` - Execute tasks with memory-first approach
- `/test` - Test features with real app
- `/ci` - Run CI checks and fix issues
- `/summarize` - Generate quick context summaries
- `/consolidate` - Update memory documentation

## Rules

- `portal.md` - React portal development guidelines
- `spec-to-code.md` - General spec-to-code workflow guidelines