package ocservGroup

import (
	"context"
	"errors"
	"fmt"
	"github.com/mmtaee/go-oc-utils/logger"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// OcGroup ocserv group
type OcGroup struct{}

// OcGroupInterface ocserv group methods
// All methods in this interface need reload server config except List and NameList.
// Use from Occtl module Reload method to reload server config in a schedule.
type OcGroupInterface interface {
	List(c context.Context) (*[]OcGroupConfigInfo, error)
	NameList(c context.Context) (*[]string, error)
	UpdateDefault(c context.Context, config *map[string]interface{}) error
	Create(c context.Context, name string, config *map[string]interface{}) error
	Update(c context.Context, name string, config *map[string]interface{}) error
	Delete(c context.Context, name string) error
}

var (
	groupDir     = "/etc/ocserv/groups"              // ocserv group configs directory path
	defaultGroup = "/etc/ocserv/defaults/group.conf" // ocserv defaults group file path
)

// NewOcGroup create new ocserv group obj
func NewOcGroup() *OcGroup {
	return &OcGroup{}
}

// List a list og ocserv group info with config data
func (g *OcGroup) List(c context.Context) (*[]OcGroupConfigInfo, error) {
	var (
		result []OcGroupConfigInfo
		wg     sync.WaitGroup
	)
	err := WithContext(c, func() error {
		err := filepath.Walk(groupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				result = append(result, OcGroupConfigInfo{
					Name: info.Name(),
					Path: path,
				})
			}
			return nil
		})
		if err != nil {
			return err
		}

		for i := range result {
			wg.Add(1)
			go func(data *OcGroupConfigInfo) {
				defer wg.Done()
				config, err := ParseConfFile(data.Path)
				if err != nil {
					fmt.Printf("Error parsing file %s: %v\n", data.Path, err)
					return
				}
				data.Config = config
			}(&result[i])
		}
		wg.Wait()
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return &result, err
}

// NameList a list of ocserv group's names
func (g *OcGroup) NameList(c context.Context) (*[]string, error) {
	var names []string
	err := WithContext(c, func() error {
		err := filepath.Walk(groupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				names = append(names, info.Name())
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return &names, nil
}

// UpdateDefault update default ocserv group configs
func (g *OcGroup) UpdateDefault(c context.Context, config *map[string]interface{}) error {
	return WithContext(c, func() error {
		file, err := os.Open(defaultGroup)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)

		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				logger.Log(logger.ERROR, fmt.Sprintf("failed to close file: %v", closeErr))
			}
		}()
		return Writer(file, config)
	})
}

// Create ocserv group creating with configs
func (g *OcGroup) Create(c context.Context, name string, config *map[string]interface{}) error {
	return WithContext(c, func() error {
		file, err := os.Create(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				logger.Log(logger.ERROR, fmt.Sprintf("failed to close file: %v", closeErr))
			}
		}()
		return Writer(file, config)
	})
}

// Update ocserv group updating with configs
func (g *OcGroup) Update(c context.Context, name string, config *map[string]interface{}) error {
	return WithContext(c, func() error {
		file, err := os.Open(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				logger.Log(logger.ERROR, fmt.Sprintf("failed to close file: %v", closeErr))
			}
		}()
		return Writer(file, config)
	})
}

// Delete ocserv group delete
func (g *OcGroup) Delete(c context.Context, name string) error {
	return WithContext(c, func() error {
		if name == "defaults" {
			return errors.New("default group cannot be deleted")
		}
		err := os.Remove(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("group %s does not exist", name)
			}
			return fmt.Errorf("failed to delete group %s: %w", name, err)
		}
		return nil
	})
}
