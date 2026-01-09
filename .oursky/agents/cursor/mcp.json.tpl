{
  "mcpServers": {
    "oursky": {
      "command": "oursky",
      "args": ["agent"],
      "env": {
        "HOST": "{{.Host}}",
        "PORT": "{{.Port}}",
        "TOKEN": "{{.AuthToken}}"
      }
    }
  }
}