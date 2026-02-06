package db

func SaveConfig(cfg SystemConfig) {
	var exist SystemConfig
	DB.First(&exist, 1)
	if exist.ID == 0 {
		cfg.ID = 1
		DB.Create(&cfg)
	} else {
		DB.Model(&exist).Updates(cfg)
	}
}

func GetConfig() SystemConfig {
	var cfg SystemConfig
	DB.First(&cfg, 1)
	return cfg
}