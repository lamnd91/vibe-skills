package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuongtl/vibe-skills/internal/registry"
)

const TargetDir = ".claude/skills"

type Installer struct {
	registry *registry.Registry
	baseDir  string
}

func New(reg *registry.Registry, baseDir string) *Installer {
	return &Installer{
		registry: reg,
		baseDir:  baseDir,
	}
}

func (i *Installer) Install(skillName string) error {
	skill := i.registry.Find(skillName)
	if skill == nil {
		return fmt.Errorf("skill not found: %s", skillName)
	}

	content, err := i.registry.GetContent(skill)
	if err != nil {
		return fmt.Errorf("failed to read skill content: %w", err)
	}

	targetPath := filepath.Join(i.baseDir, TargetDir, skill.Name+".md")

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write skill file
	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	return nil
}

func (i *Installer) InstallMultiple(skillNames []string) (installed []string, errors []error) {
	for _, name := range skillNames {
		if err := i.Install(name); err != nil {
			errors = append(errors, err)
		} else {
			installed = append(installed, name)
		}
	}
	return
}

func (i *Installer) InstallStack(stack string) (installed []string, errors []error) {
	skills := i.registry.ListByStack(stack)
	if len(skills) == 0 {
		errors = append(errors, fmt.Errorf("no skills found in stack: %s", stack))
		return
	}

	for _, skill := range skills {
		if err := i.Install(skill.Name); err != nil {
			errors = append(errors, err)
		} else {
			installed = append(installed, skill.Name)
		}
	}
	return
}

func (i *Installer) InstallAll() (installed []string, errors []error) {
	skills := i.registry.List()
	for _, skill := range skills {
		if err := i.Install(skill.Name); err != nil {
			errors = append(errors, err)
		} else {
			installed = append(installed, skill.Name)
		}
	}
	return
}

func (i *Installer) Remove(skillName string) error {
	targetPath := filepath.Join(i.baseDir, TargetDir, skillName+".md")

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("skill not installed: %s", skillName)
	}

	if err := os.Remove(targetPath); err != nil {
		return fmt.Errorf("failed to remove skill: %w", err)
	}

	return nil
}

func (i *Installer) ListInstalled() ([]string, error) {
	targetDir := filepath.Join(i.baseDir, TargetDir)

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var installed []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			name := entry.Name()[:len(entry.Name())-3] // Remove .md extension
			installed = append(installed, name)
		}
	}

	return installed, nil
}

func (i *Installer) IsInstalled(skillName string) bool {
	targetPath := filepath.Join(i.baseDir, TargetDir, skillName+".md")
	_, err := os.Stat(targetPath)
	return err == nil
}
