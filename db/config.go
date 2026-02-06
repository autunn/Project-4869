package db

func SaveConfig(cfg SystemConfig) {
	var current SystemConfig
	DB.First(&current, 1)
	if current.ID == 0 {
		cfg.ID = 1
		DB.Create(&cfg)
	} else {
		DB.Model(&current).Updates(cfg)
	}
}

func GetConfig() SystemConfig {
	var cfg SystemConfig
	DB.First(&cfg, 1)
	return cfg
}