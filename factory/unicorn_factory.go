package factory

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"unicorn/model"
	"unicorn/utils"
)

type UnicornFactory interface {
	CreateUnicorn() model.Unicorn
}
type RandomUnicornProducer struct {
	names        []string
	adjectives   []string
	capabilities []string
	counter      int
}

func (s *RandomUnicornProducer) CreateUnicorn() model.Unicorn {
	s.counter++
	return model.Unicorn{
		Name:         fmt.Sprintf("%d - %s - %s", s.counter, s.getName(), s.getAdjective()),
		Capabilities: s.getNUniqueCapability(3),
	}
}

func loadData() UnicornFactory {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to get home dir", err)
	}
	dir = filepath.Join(dir, "GolandProjects/unicorn-main")
	names := utils.LoadContentFromFile(filepath.Join(dir, "res/petnames.txt"))
	adj := utils.LoadContentFromFile(filepath.Join(dir, "res/adj.txt"))
	capabilities := utils.LoadContentFromFile(filepath.Join(dir, "res/capabilities.txt"))
	return &RandomUnicornProducer{
		names:        names,
		adjectives:   adj,
		capabilities: capabilities,
	}
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

func NewRandomUnicornProducer() UnicornFactory {
	return loadData()
}
