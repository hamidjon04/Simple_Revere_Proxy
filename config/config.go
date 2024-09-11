package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"proxy/model"
)

func LoadConfig(logger *slog.Logger) (*model.ConfigSettings, error) {
	var settings model.ConfigSettings
	file, err := os.Open("config/config.json")
	if err != nil {
		logger.Error(fmt.Sprintf("Configuratsiya file ochilmadi: %v", err))
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&settings); err != nil {
		logger.Error(fmt.Sprintf("Configuratsiya filedagi ma'lumotlarni o'qib bo'lmadi: %v", err))
		return nil, err
	}
	return &settings, nil
}
