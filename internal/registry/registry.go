package registry

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed skills
var embeddedSkills embed.FS

type Skill struct {
	Name        string
	Stack       string
	Path        string
	Description string
}

type Registry struct {
	skills []Skill
}

func New() (*Registry, error) {
	r := &Registry{}
	if err := r.load(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) load() error {
	r.skills = []Skill{}

	err := fs.WalkDir(embeddedSkills, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if d.Name() != "SKILL.md" {
			return nil
		}

		// Extract stack and skill name from path
		// skills/common/commit-convention/SKILL.md -> stack=common, name=commit-convention
		parts := strings.Split(path, string(filepath.Separator))
		if len(parts) < 4 {
			return nil
		}

		stack := parts[1]
		name := parts[2]

		// Read description from first line of SKILL.md
		content, err := embeddedSkills.ReadFile(path)
		if err != nil {
			return nil
		}

		description := extractDescription(string(content))

		r.skills = append(r.skills, Skill{
			Name:        name,
			Stack:       stack,
			Path:        path,
			Description: description,
		})

		return nil
	})

	return err
}

func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Return first non-empty, non-header line
		if len(line) > 100 {
			return line[:100] + "..."
		}
		return line
	}
	return ""
}

func (r *Registry) List() []Skill {
	return r.skills
}

func (r *Registry) ListByStack(stack string) []Skill {
	var result []Skill
	for _, s := range r.skills {
		if s.Stack == stack {
			result = append(result, s)
		}
	}
	return result
}

func (r *Registry) GetStacks() []string {
	stackMap := make(map[string]bool)
	for _, s := range r.skills {
		stackMap[s.Stack] = true
	}

	var stacks []string
	for stack := range stackMap {
		stacks = append(stacks, stack)
	}
	return stacks
}

func (r *Registry) Find(name string) *Skill {
	for _, s := range r.skills {
		if s.Name == name {
			return &s
		}
		// Also match full path like "common/commit-convention"
		if s.Stack+"/"+s.Name == name {
			return &s
		}
	}
	return nil
}

func (r *Registry) Search(query string) []Skill {
	query = strings.ToLower(query)
	var result []Skill
	for _, s := range r.skills {
		if strings.Contains(strings.ToLower(s.Name), query) ||
			strings.Contains(strings.ToLower(s.Description), query) ||
			strings.Contains(strings.ToLower(s.Stack), query) {
			result = append(result, s)
		}
	}
	return result
}

func (r *Registry) GetContent(skill *Skill) ([]byte, error) {
	return embeddedSkills.ReadFile(skill.Path)
}
