{
  "mcpServers": {
    "oursky": {
      "command": "mcp-remote",
      "args": ["{{.Host}}:{{.Port}}"],
      "env": {
        "AUTH_TOKEN": "{{.AuthToken}}"
      }
    }
  }
}