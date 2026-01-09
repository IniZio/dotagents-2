# Vendatta

Vendatta eliminates the "it works on my machine" problem by providing isolated, reproducible development environments that work seamlessly with Coding Agents e.g. Cursor, OpenCode, Claude, etc.

## Key Features

- **Single Binary**: Zero-setup installation with no host dependencies
- **Branch Isolation**: Git worktrees provide unique filesystems for every branch
- **AI Agent Integration**: Automatic configuration for Cursor, OpenCode, Claude, and more via Model Context Protocol (MCP)
- **Service Discovery**: Automatic port mapping and environment variables for multi-service apps
- **Docker-in-Docker**: Run docker-compose projects inside isolated environments

## Quick Start

### Try It Now

Get started in 2 simple steps:

```bash
# 1. Install Vendatta
curl -fsSL https://raw.githubusercontent.com/oursky/vendatta/main/install.sh | bash

# Add ~/.local/bin to your PATH if not already:
# export PATH="$HOME/.local/bin:$PATH"

# 2. Initialize in your project
oursky init

# 3. Start an isolated development session
oursky dev my-feature
```

That's it! Vendatta creates an isolated environment for your `my-feature` branch with automatic AI agent configuration.

#### Alternative: Build from Source

If you prefer to build from source:

```bash
# Requires Go 1.24+
go build -o oursky cmd/oursky/main.go
```

#### Updates

To update to the latest version:

```bash
oursky update
```

### Understanding What Happened

- **Step 1**: Built a single Go binary that manages everything
- **Step 2**: Created a `.oursky/` directory with basic configuration templates
- **Step 3**: Generated a Git worktree at `.oursky/worktrees/my-feature/` and started any configured services

Your AI agents (Cursor, OpenCode, etc.) are now automatically configured to work with this isolated environment.

### Configure for Your Project

Vendatta works with your existing development setup. Edit `.oursky/config.yaml` to define your services:

```yaml
# Example: Full-stack web app
services:
  db:
    command: "docker-compose up -d postgres"
  api:
    command: "cd server && npm run dev"
    depends_on: ["db"]
  web:
    command: "cd client && npm run dev"
    depends_on: ["api"]

# Enable AI agents
agents:
  - name: "cursor"
    enabled: true
  - name: "opencode"
    enabled: true
```

Run `./oursky dev my-feature` again to apply your configuration.

## When to Use Vendatta

### Branch-Based Development
Perfect for teams working on multiple features simultaneously:
- Each branch gets its own isolated filesystem
- No more "git stash" or conflicting dependencies
- Parallel development without environment pollution

### Complex Microservices
When your local setup involves multiple services:
- Databases, APIs, frontend apps
- Docker-compose projects run inside containers
- Automatic service discovery and port mapping

### AI-Assisted Development
Enhance your AI coding experience:
- Secure tool execution for agents
- Project-specific rules and capabilities
- Standardized skills across different AI tools

### Team Standardization
Ensure consistent development environments:
- Version-controlled configurations
- Shared templates for coding standards
- Easy onboarding for new team members

## AI Agent Integration

Vendatta automatically configures your favorite AI coding assistants to work securely with isolated environments.

### Supported Agents

| Agent | Description | Integration |
|-------|-------------|-------------|
| **Cursor** | VS Code with AI | `.cursor/mcp.json` |
| **OpenCode** | Standalone AI assistant | `opencode.json` + `.opencode/` |
| **Claude Desktop** | Anthropic's desktop app | `claude_desktop_config.json` |
| **Claude Code** | CLI tool | `claude_code_config.json` |

### How It Works

1. **Enable agents** in `.oursky/config.yaml`
2. **Start development**: `./oursky dev branch-name`
3. **Open your worktree** in the AI agent of choice
4. **Agents connect automatically** via MCP with full environment access

### Shared Capabilities

All enabled agents get access to:
- **Skills**: Web search, file operations, data analysis
- **Commands**: Spec-to-code workflow commands (design, plan, implement, test, ci)
- **Rules**: Code quality standards, collaboration guidelines

### Template Management

Oursky includes a powerful template management system for easy collaboration:

```bash
# Pull shared templates from your organization
oursky templates pull https://github.com/your-org/dotagents.git

# Edit templates locally
oursky templates checkout dotagents
oursky templates edit  # Shows guidance

# Test changes
oursky templates apply

# Share improvements
oursky templates push dotagents
```

### Agent Config Sharing

Share standardized AI configurations across your team:
- **Remote repositories**: Use GitHub repos for team-wide template sharing
- **Chezmoi-like workflow**: Pull, edit, test, and push template changes
- **Version-controlled setups**: Store agent configurations in git alongside your code
- **Team consistency**: Ensure all developers use the same AI skills, commands, and rules
- **Easy collaboration**: New team members get identical AI assistance setups
- **Cross-agent compatibility**: Configurations work seamlessly across Cursor, OpenCode, Claude, and other agents

Use `oursky templates pull` to get started with shared configurations.

## Configuration Reference

### Project Structure
```
.oursky/
├── config.yaml          # Main project configuration
├── templates/           # Working directory for editing templates
│   ├── skills/          # Reusable AI skills
│   ├── commands/        # Development commands
│   └── rules/           # Coding guidelines
├── remotes/             # Cloned template repositories
│   └── dotagents/       # Remote template repo (e.g., team configs)
├── agents/              # Agent-specific overrides
└── worktrees/           # Auto-generated environments
```

### Main Configuration

The `.oursky/config.yaml` file defines your development environment:

```yaml
# Project settings
name: "my-web-app"

# Services to run
services:
  db:
    command: "docker-compose up -d postgres"
    healthcheck:
      url: "http://localhost:5432/health"
  api:
    command: "cd server && npm run dev"
    depends_on: ["db"]
  web:
    command: "cd client && npm run dev"
    depends_on: ["api"]

# AI agents to configure
agents:
  - name: "cursor"
    enabled: true
  - name: "opencode"
    enabled: true

# MCP server settings
mcp:
  enabled: true
  port: 3001
```

### Customizing Templates

#### Adding AI Skills
Create `.oursky/templates/skills/my-skill.yaml`:
```yaml
name: "my-custom-skill"
description: "Does something useful"
parameters:
  type: object
  properties:
    input: { type: "string" }
execute:
  command: "node"
  args: ["scripts/my-skill.js"]
```

#### Defining Commands
Create `.oursky/templates/commands/my-command.yaml`:
```yaml
name: "deploy"
description: "Deploy to staging"
steps:
  - name: "Build"
    command: "npm run build"
  - name: "Deploy"
    command: "kubectl apply -f k8s/"
```

#### Setting Coding Rules
Create `.oursky/templates/rules/team-standards.md`:
```markdown
---
title: "Team Standards"
applies_to: ["**/*.ts", "**/*.js"]
---

# Code Quality Standards
- Use TypeScript for new code
- Functions should be < 30 lines
- Always add return types
```

### Environment Variables

Use variables for dynamic configuration:

```yaml
# In config.yaml
mcp:
  port: "{{.Env.MCP_PORT}}"
```

```bash
export MCP_PORT=3001
./oursky dev my-branch
```

### Service Discovery & Port Access

Vendatta automatically discovers running services and provides environment variables for easy access:

**Available in worktrees**: When you run `./oursky dev branch-name`, your worktree environment gets these variables:

- `OURSKY_SERVICE_DB_URL` - Database connection URL
- `OURSKY_SERVICE_API_URL` - API service URL
- `OURSKY_SERVICE_WEB_URL` - Web frontend URL
- And more for each service you define

**Example usage in your code**:

```javascript
// In your frontend config
const apiUrl = process.env.OURSKY_SERVICE_API_URL || 'http://localhost:3001';

// In your API config
const dbUrl = process.env.OURSKY_SERVICE_DB_URL;
```

**Check available services**:

```bash
# In your worktree directory
env | grep OURSKY_SERVICE
```

This eliminates manual port management and ensures your services can communicate seamlessly across the isolated environment.

## Example: Full-Stack Development

1. **Set up your project**:
   ```bash
   ./oursky init
   ```

2. **Configure templates** (optional, for shared team configurations):
   ```bash
   oursky templates pull https://github.com/your-org/dotagents.git
   oursky templates checkout dotagents
   oursky templates apply
   ```

4. **Configure services** (edit `.oursky/config.yaml`):
   ```yaml
   services:
     db:
       command: "docker-compose up -d postgres"
     api:
       command: "cd server && npm run dev"
       depends_on: ["db"]
     web:
       command: "cd client && npm run dev"
       depends_on: ["api"]

   agents:
     - name: "cursor"
       enabled: true
   ```

5. **Start development**:
   ```bash
   ./oursky dev new-feature
   ```

6. **Code with AI assistance**:
   - Open `.oursky/worktrees/new-feature/` in Cursor
   - AI agent connects automatically with full environment access

## Customizing Templates

To customize AI agent capabilities:

```bash
# Edit templates locally
oursky templates edit

# Add custom skills in .oursky/templates/skills/
# Add custom commands in .oursky/templates/commands/
# Add coding rules in .oursky/templates/rules/

# Test changes
oursky templates apply

# Share with team
oursky templates push dotagents
```

For team-wide templates, create a GitHub repository and use `oursky templates pull` to share configurations.

---
*Powered by OhMyOpenCode.*
