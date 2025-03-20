package storage

import (
	"fmt"
	"log"
	"math/rand"
	"unicorn/model"
	"unicorn/utils"
)

type UnicornProducer interface {
	CreateUnicorn() model.Unicorn
}
type RandomUnicornProducer struct {
	names        []string
	adjectives   []string
	capabilities []string
}

func (s *RandomUnicornProducer) CreateUnicorn() model.Unicorn {
	return model.Unicorn{
		Name:         fmt.Sprintf("%s - %s", s.getName(), s.getAdjective()),
		Capabilities: s.getNUniqueCapability(3),
	}
}

func loadData() UnicornProducer {
	names := loadContentFromFile("res/petnames.txt")
	adj := loadContentFromFile("res/adj.txt")
	capabilities := loadContentFromFile("res/capabilities.txt")
	return &RandomUnicornProducer{
		names:        names,
		adjectives:   adj,
		capabilities: capabilities,
	}
}

func loadContentFromFile(filepath string) []string {
	content := utils.GetFileContent(filepath)
	if content == nil {
		log.Fatalf("unable to load content from file: %s", filepath)
	}
	return content
}

func (s *RandomUnicornProducer) getName() string {
	return s.names[rand.Intn(len(s.names)-1)]
}

func (s *RandomUnicornProducer) getAdjective() string {
	return s.adjectives[rand.Intn(len(s.adjectives)-1)]
}

func (s *RandomUnicornProducer) getNUniqueCapability(n int) []string {
	if n >= len(s.capabilities) {
		return []string{}
	}
	rand.Shuffle(len(s.capabilities), func(i, j int) {
		s.capabilities[i], s.capabilities[j] = s.capabilities[j], s.capabilities[i]
	})
	capabilitiesList := make([]string, 3)
	copy(capabilitiesList, s.capabilities[:n])
	return capabilitiesList
}

func NewRandomUnicornProducer() UnicornProducer {
	return loadData()
}
