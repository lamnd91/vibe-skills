package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuongtl1992/vibe-skills/internal/registry"
)

const TargetDir = ".claude/skills"

// SkillProvider defines the interface for skill sources
type SkillProvider interface {
	Find(name string) (*registry.Skill, error)
	List() ([]registry.Skill, error)
	ListByStack(stack string) ([]registry.Skill, error)
	GetContent(skill *registry.Skill) ([]byte, error)
	GetFiles(skill *registry.Skill) (map[string][]byte, error)
}

type Installer struct {
	provider SkillProvider
	baseDir  string
}

func New(provider SkillProvider, baseDir string) *Installer {
	return &Installer{
		provider: provider,
		baseDir:  baseDir,
	}
}

func (i *Installer) Install(skillName string) error {
	skill, err := i.provider.Find(skillName)
	if err != nil {
		return fmt.Errorf("skill not found: %s", skillName)
	}

	// Determine if multi-file skill
	isMultiFile := len(skill.Files) > 0

	if isMultiFile {
		return i.installMultiFile(skill)
	}
	return i.installSingleFile(skill)
}

func (i *Installer) installSingleFile(skill *registry.Skill) error {
	content, err := i.provider.GetContent(skill)
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

func (i *Installer) installMultiFile(skill *registry.Skill) error {
	files, err := i.provider.GetFiles(skill)
	if err != nil {
		return fmt.Errorf("failed to fetch skill files: %w", err)
	}

	skillDir := filepath.Join(i.baseDir, TargetDir, skill.Name)

	for relPath, content := range files {
		fullPath := filepath.Join(skillDir, relPath)

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", relPath, err)
		}

		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", relPath, err)
		}
	}

	return nil
}

func (i *Installer) InstallMultiple(skillNames []string) (installed []string, errors []error) {
	for _, name := range skillNames {
		if err := i.Install(name); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		} else {
			installed = append(installed, name)
		}
	}
	return
}

func (i *Installer) InstallStack(stack string) (installed []string, errors []error) {
	skills, err := i.provider.ListByStack(stack)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to list stack %s: %w", stack, err))
		return
	}
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
	skills, err := i.provider.List()
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to list skills: %w", err))
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

func (i *Installer) Remove(skillName string) error {
	// Check for directory first (multi-file skill)
	dirPath := filepath.Join(i.baseDir, TargetDir, skillName)
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		return os.RemoveAll(dirPath)
	}

	// Fallback to single file
	filePath := filepath.Join(i.baseDir, TargetDir, skillName+".md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("skill not installed: %s", skillName)
	}

	if err := os.Remove(filePath); err != nil {
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
		if entry.IsDir() {
			// Multi-file skill: check for SKILL.md inside
			skillMd := filepath.Join(targetDir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillMd); err == nil {
				installed = append(installed, entry.Name())
			}
		} else if filepath.Ext(entry.Name()) == ".md" {
			// Single-file skill
			name := entry.Name()[:len(entry.Name())-3] // Remove .md extension
			installed = append(installed, name)
		}
	}

	return installed, nil
}

func (i *Installer) IsInstalled(skillName string) bool {
	// Check directory (multi-file skill)
	dirPath := filepath.Join(i.baseDir, TargetDir, skillName)
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		skillMd := filepath.Join(dirPath, "SKILL.md")
		_, err := os.Stat(skillMd)
		return err == nil
	}

	// Check single file
	filePath := filepath.Join(i.baseDir, TargetDir, skillName+".md")
	_, err := os.Stat(filePath)
	return err == nil
}
