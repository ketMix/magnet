package data

import (
	"path"
)

type Config interface {
	LoadFromFile(p string)
}

var (
	PlayerInit    EntityConfig
	CoreConfig    EntityConfig
	TurretConfigs map[string]EntityConfig
	EnemyConfigs  map[string]EntityConfig
)

func LoadConfigurations() error {
	// Load the player configuration
	config, err := NewPlayerConfig()
	PlayerInit = config
	if err != nil {
		return err
	}

	// Load the core configuration
	CoreConfig, err = NewCoreConfig()
	if err != nil {
		return err
	}

	// Traverse the turret config folder and load all turret configurations
	TurretConfigs = make(map[string]EntityConfig)
	turretFiles, err := getPathFiles(path.Join("entities", "turrets"))
	println("Loading turret configs:")
	if err != nil {
		return err
	}
	for _, fileName := range turretFiles {
		println("\t", fileName)
		TurretConfigs[fileName], err = NewTurretConfig(fileName)
		if err != nil {
			return err
		}
	}

	// Traverse the enemy config folder and load all enemy configurations
	EnemyConfigs = make(map[string]EntityConfig)
	enemyFiles, err := getPathFiles(path.Join("entities", "enemies"))
	println("Loading enemy configs:")
	if err != nil {
		return err
	}
	for _, fileName := range enemyFiles {
		println("\t", fileName)
		EnemyConfigs[fileName], err = NewEnemyConfig(fileName)
		if err != nil {
			return err
		}
	}
	return nil
}
