package db

type Vault struct {
	ID      int
	Name    string
	Created string
}

func GetVaults() ([]Vault, error) {
	rows, err := DB.Query("SELECT id, name, created_at FROM vaults ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vaults []Vault
	for rows.Next() {
		var vault Vault
		err := rows.Scan(&vault.ID, &vault.Name, &vault.Created)
		if err != nil {
			continue
		}
		vaults = append(vaults, vault)
	}
	return vaults, nil
}

func CreateVault(name string) error {
	_, err := DB.Exec("INSERT INTO vaults (name) VALUES (?)", name)
	return err
}

func DeleteVault(name string) (int64, error) {
	result, err := DB.Exec("DELETE FROM vaults WHERE name = ?", name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetVaultByName(name string) (*Vault, error) {
	var vault Vault
	err := DB.QueryRow("SELECT id, name, created_at FROM vaults WHERE name = ?", name).
		Scan(&vault.ID, &vault.Name, &vault.Created)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}
