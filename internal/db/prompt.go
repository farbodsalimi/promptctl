package db

import "fmt"

type Prompt struct {
	ID            int
	Name          string
	Created       string
	LatestVersion int
	VaultName     string
}

func GetPrompts(vaultName string) ([]Prompt, error) {
	query := `
		SELECT p.id, p.name, p.created_at, MAX(pv.version) as latest_version, v.name as vault_name
		FROM prompts p
		JOIN vaults v ON p.vault_id = v.id
		LEFT JOIN prompt_versions pv ON p.id = pv.prompt_id
		WHERE v.name = ?
		GROUP BY p.id, p.name, p.created_at, v.name
		ORDER BY p.created_at DESC
	`
	rows, err := DB.Query(query, vaultName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []Prompt
	for rows.Next() {
		var prompt Prompt
		err := rows.Scan(
			&prompt.ID,
			&prompt.Name,
			&prompt.Created,
			&prompt.LatestVersion,
			&prompt.VaultName,
		)
		if err != nil {
			continue
		}
		prompts = append(prompts, prompt)
	}
	return prompts, nil
}

func GetPromptContent(promptID int) (string, error) {
	query := `
		SELECT content FROM prompt_versions pv
		JOIN prompts p ON pv.prompt_id = p.id
		WHERE p.id = ?
		ORDER BY pv.version DESC
		LIMIT 1
	`
	var content string
	err := DB.QueryRow(query, promptID).Scan(&content)
	return content, err
}

func CreatePrompt(vaultID int, name, content string) error {
	var promptID int
	err := DB.QueryRow("INSERT INTO prompts (vault_id, name) VALUES (?, ?) RETURNING id",
		vaultID, name).Scan(&promptID)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"INSERT INTO prompt_versions (prompt_id, version, content) VALUES (?, 1, ?)",
		promptID,
		content,
	)
	return err
}

func GetPromptVersionContent(promptID int) (int, string, error) {
	query := `
		SELECT pv.id, pv.content FROM prompt_versions pv
		JOIN prompts p ON pv.prompt_id = p.id
		WHERE p.id = ?
		ORDER BY pv.version DESC
		LIMIT 1
	`
	var promptVersionID int
	var content string
	err := DB.QueryRow(query, promptID).Scan(&promptVersionID, &content)
	return promptVersionID, content, err
}

func CreateRun(promptVersionID int, provider, params, response string) error {
	_, err := DB.Exec(
		"INSERT INTO runs (prompt_version_id, provider, params, response) VALUES (?, ?, ?, ?)",
		promptVersionID, provider, params, response)
	return err
}

func GetPromptHistory(promptID int) ([]string, error) {
	query := `
		SELECT pv.version, pv.created_at
		FROM prompt_versions pv
		WHERE pv.prompt_id = ?
		ORDER BY pv.version DESC
	`
	rows, err := DB.Query(query, promptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []string
	for rows.Next() {
		var version int
		var created string
		err := rows.Scan(&version, &created)
		if err != nil {
			continue
		}
		history = append(history, fmt.Sprintf("v%d - %s", version, created))
	}
	return history, nil
}

func GetPromptByName(vaultName, promptName string) (*Prompt, error) {
	query := `
		SELECT p.id, p.name, p.created_at, MAX(pv.version) as latest_version, v.name as vault_name
		FROM prompts p
		JOIN vaults v ON p.vault_id = v.id
		LEFT JOIN prompt_versions pv ON p.id = pv.prompt_id
		WHERE v.name = ? AND p.name = ?
		GROUP BY p.id, p.name, p.created_at, v.name
	`
	var prompt Prompt
	err := DB.QueryRow(query, vaultName, promptName).Scan(
		&prompt.ID,
		&prompt.Name,
		&prompt.Created,
		&prompt.LatestVersion,
		&prompt.VaultName,
	)
	if err != nil {
		return nil, err
	}
	return &prompt, nil
}

func UpdatePrompt(vaultName, promptName, content string) error {
	// Get the prompt
	prompt, err := GetPromptByName(vaultName, promptName)
	if err != nil {
		return err
	}

	// Get the next version number
	nextVersion := prompt.LatestVersion + 1

	// Insert new version
	_, err = DB.Exec(
		"INSERT INTO prompt_versions (prompt_id, version, content) VALUES (?, ?, ?)",
		prompt.ID,
		nextVersion,
		content,
	)
	return err
}

func GetPromptVersionContentByVersion(vaultName, promptName string, version int) (string, error) {
	query := `
		SELECT pv.content
		FROM prompt_versions pv
		JOIN prompts p ON pv.prompt_id = p.id
		JOIN vaults v ON p.vault_id = v.id
		WHERE v.name = ? AND p.name = ? AND pv.version = ?
	`
	var content string
	err := DB.QueryRow(query, vaultName, promptName, version).Scan(&content)
	return content, err
}

type Run struct {
	ID         int
	PromptName string
	VaultName  string
	Provider   string
	Params     string
	Response   string
	Created    string
}

func GetRuns(vaultName, promptName string) ([]Run, error) {
	query := `
		SELECT r.id, p.name as prompt_name, v.name as vault_name, 
		       r.provider, r.params, r.response, r.created_at
		FROM runs r
		JOIN prompt_versions pv ON r.prompt_version_id = pv.id
		JOIN prompts p ON pv.prompt_id = p.id
		JOIN vaults v ON p.vault_id = v.id
	`
	var args []any

	conditions := []string{}
	if vaultName != "" {
		conditions = append(conditions, "v.name = ?")
		args = append(args, vaultName)
	}
	if promptName != "" {
		conditions = append(conditions, "p.name = ?")
		args = append(args, promptName)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY r.created_at DESC LIMIT 20"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []Run
	for rows.Next() {
		var run Run
		err := rows.Scan(&run.ID, &run.PromptName, &run.VaultName,
			&run.Provider, &run.Params, &run.Response, &run.Created)
		if err != nil {
			continue
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func GetRunByID(id int) (*Run, error) {
	query := `
		SELECT r.id, p.name as prompt_name, v.name as vault_name, 
		       r.provider, r.params, r.response, r.created_at
		FROM runs r
		JOIN prompt_versions pv ON r.prompt_version_id = pv.id
		JOIN prompts p ON pv.prompt_id = p.id
		JOIN vaults v ON p.vault_id = v.id
		WHERE r.id = ?
	`
	var run Run
	err := DB.QueryRow(query, id).Scan(&run.ID, &run.PromptName, &run.VaultName,
		&run.Provider, &run.Params, &run.Response, &run.Created)
	if err != nil {
		return nil, err
	}
	return &run, nil
}
