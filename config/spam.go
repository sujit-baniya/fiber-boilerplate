package config

import (
	"github.com/sujit-baniya/shield/spam"
	"path/filepath"
)

type SpamConfig struct {
	*spam.DetectionEngine
	TrainFile string `yaml:"train_file" env-default:"spam.csv"`
}

func (s *SpamConfig) Load(assetPath string) {
	spamDetectEngine := spam.NewSpamDetectionEngine()
	spamDetectEngine.TrainFromFile(filepath.Join(assetPath, s.TrainFile))
	s.DetectionEngine = spamDetectEngine
}
