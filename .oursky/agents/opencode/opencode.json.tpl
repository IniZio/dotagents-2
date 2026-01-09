{
  "name": "oursky",
  "command": "oursky",
  "args": ["agent"],
  "env": {
    "HOST": "{{.Host}}",
    "PORT": "{{.Port}}",
    "TOKEN": "{{.AuthToken}}"
  },
  "rules": "{{.RulesConfig}}",
  "skills": "{{.SkillsConfig}}",
  "commands": "{{.CommandsConfig}}"
}