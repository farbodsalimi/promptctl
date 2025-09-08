package templates

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

func RenderTemplate(content string, vars map[string]any) (string, error) {
	tmpl, err := template.New("prompt").Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParseVars(varsStr string) (map[string]any, error) {
	vars := make(map[string]any)

	if varsStr == "" {
		return vars, nil
	}

	// Try to parse as JSON first
	if strings.HasPrefix(varsStr, "{") {
		err := json.Unmarshal([]byte(varsStr), &vars)
		return vars, err
	}

	// Parse as comma-separated key=value pairs
	pairs := strings.Split(varsStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			vars[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	return vars, nil
}
