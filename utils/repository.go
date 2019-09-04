package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Component interface {
	GetComponentId() string
}

type Repository struct {
	Name       string
	components map[string]Component
}

func (r *Repository) Register(com Component) {
	comId := strings.ToLower(com.GetComponentId())
	if _, ok := r.components[comId]; ok {
		log.Fatalf("Component '%s' already registed in %s", comId, r.Name)
		return
	}
	r.components[comId] = com
	log.Printf("registered %s in %s", comId, r.Name)
}

func (r *Repository) Get(insId string) (Component, error) {
	if component, ok := r.components[insId]; !ok {
		return nil, errors.New(fmt.Sprintf("invalid inspectorId %s", insId))
	} else {
		return component, nil
	}
}

func (r *Repository) GetComponentIds() []string {
	var ids []string
	for id := range r.components {
		ids = append(ids, id)
	}
	return ids
}

func CreateRepository(name string) Repository {
	return Repository{
		Name:       name,
		components: make(map[string]Component),
	}
}
