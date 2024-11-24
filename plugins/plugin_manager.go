package plugins

import "fmt"

type PluginManager struct {
	handlers map[string]func([]string) error
}

func NewPluginManager() *PluginManager {
	return &PluginManager{handlers: map[string]func([]string) error{}}
}

func (pm *PluginManager) RegisterKeyword(keyword string, handler func([]string) error) {
	pm.handlers[keyword] = handler
}

func (pm *PluginManager) ExecuteKeyword(keyword string, args []string) error {
	handler, exists := pm.handlers[keyword]
	if !exists {
		return fmt.Errorf("unknown keyword: %s", keyword)
	}
	return handler(args)
}