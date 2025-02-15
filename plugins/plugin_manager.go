package plugins

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type MethodInfo struct {
	pluginName string
	funcName   string
}

type Plugin interface {
	Start(ctx context.Context) error
	Stop() error
	ListMethods() map[string]string
}

type PluginManager struct {
	lock         sync.RWMutex
	plugins      map[string]Plugin
	methodToInfo map[string]MethodInfo
}

var (
	instance *PluginManager
	once     sync.Once
)

func GetInstance() *PluginManager {
	once.Do(func() {
		instance = &PluginManager{
			plugins:      make(map[string]Plugin),
			methodToInfo: make(map[string]MethodInfo),
		}
	})
	return instance
}

func (pm *PluginManager) Register(ctx context.Context, name string, plugin Plugin) error {
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %q already registered", name)
	}
	if err := plugin.Start(ctx); err != nil {
		return fmt.Errorf("failed to init plugin %q: %w", name, err)
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()

	pm.plugins[name] = plugin
	for methodName, methodFunc := range plugin.ListMethods() {
		pm.methodToInfo[methodName] = MethodInfo{
			pluginName: name,
			funcName:   methodFunc,
		}
	}
	return nil
}

func (pm *PluginManager) Unregister(name string) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %q not found", name)
	}
	if err := plugin.Stop(); err != nil {
		return fmt.Errorf("failed to stop plugin %q: %w", name, err)
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()
	delete(pm.plugins, name)
	return nil
}

func (pm *PluginManager) UnregisterAll() error {
	var err error
	for name := range pm.plugins {
		pluginErr := pm.Unregister(name)
		if pluginErr != nil {
			err = pluginErr
		}
	}
	return err
}

// GetPluginInfo returns (pluginName, functionName)
func (pm *PluginManager) GetPluginInfo(methodName string) (string, string, bool) {
	if info, ok := pm.methodToInfo[methodName]; ok {
		return info.pluginName, info.funcName, true
	}
	return "", "", false
}

func (pm *PluginManager) Execute(pluginName string, methodName string, args map[string]interface{}) (interface{}, error) {
	pm.lock.RLock()
	plugin, exists := pm.plugins[pluginName]
	pm.lock.RUnlock()
	if !exists {
		return nil, fmt.Errorf("plugin %q not found", pluginName)
	}
	// Use reflection to get the plugin's method
	pluginValue := reflect.ValueOf(plugin)
	method := pluginValue.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %q not found in plugin %q", methodName, pluginName)
	}

	// Validate and call the method
	// inputArgs := make([]reflect.Value, len(args))
	// for i, arg := range args {
	//     inputArgs[i] = reflect.ValueOf(arg)
	// }
	inputArgs := []reflect.Value{reflect.ValueOf(args)}

	// Call the method and handle results
	results := method.Call(inputArgs)
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) == 1 {
		return results[0].Interface(), nil
	}
	if len(results) == 2 && results[1].IsNil() {
		return results[0].Interface(), nil
	} else if len(results) == 2 && !results[1].IsNil() {
		return results[0].Interface(), results[1].Interface().(error)
	}

	// Return multiple results as a slice
	output := make([]interface{}, len(results))
	for i, result := range results {
		output[i] = result.Interface()
	}
	return output, nil
}
