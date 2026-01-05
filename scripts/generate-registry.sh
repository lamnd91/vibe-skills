#!/bin/bash

# Generate registry.json from SKILL.md files
# Usage: ./scripts/generate-registry.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
SKILLS_DIR="$ROOT_DIR/skills"
OUTPUT_FILE="$SKILLS_DIR/registry.json"

echo "Scanning skills in $SKILLS_DIR..."

# Start JSON
echo '{' > "$OUTPUT_FILE"
echo '  "version": "1.0",' >> "$OUTPUT_FILE"
echo '  "skills": [' >> "$OUTPUT_FILE"

first=true
skill_count=0

for skill_file in $(find "$SKILLS_DIR" -name "SKILL.md" -type f | sort); do
  # Extract stack and name from path
  # skills/common/commit-convention/SKILL.md -> stack=common, name=commit-convention
  relative_path="${skill_file#$SKILLS_DIR/}"
  stack=$(echo "$relative_path" | cut -d'/' -f1)
  name=$(echo "$relative_path" | cut -d'/' -f2)
  path="${relative_path%/SKILL.md}/SKILL.md"

  # Check if file has YAML frontmatter (starts with ---)
  if head -1 "$skill_file" | grep -q '^---$'; then
    # Extract frontmatter content between first and second ---
    frontmatter=$(sed -n '2,/^---$/p' "$skill_file" | sed '$d')

    # Extract name from frontmatter (override folder name if present)
    fm_name=$(echo "$frontmatter" | grep '^name:' | sed 's/^name:[[:space:]]*//')
    if [ -n "$fm_name" ]; then
      name="$fm_name"
    fi

    # Extract description from frontmatter
    fm_desc=$(echo "$frontmatter" | grep '^description:' | sed 's/^description:[[:space:]]*//')
    if [ -n "$fm_desc" ]; then
      description="$fm_desc"
    else
      # Fallback: first non-header, non-empty line after frontmatter
      description=$(sed -n '/^---$/,/^---$/d; /^#/d; /^$/d; /^\`\`\`/d; p' "$skill_file" | head -1)
    fi
  else
    # No frontmatter: extract description from first non-empty, non-header line
    description=$(grep -v '^#' "$skill_file" | grep -v '^$' | grep -v '^\`\`\`' | head -1)
  fi

  # Escape quotes and truncate description
  description=$(echo "$description" | sed 's/"/\\"/g' | cut -c1-200)

  # Find additional files in skill directory (excluding SKILL.md and hidden files)
  skill_dir=$(dirname "$skill_file")
  additional_files=""
  while IFS= read -r -d '' file; do
    # Get relative path from skill directory
    rel_path="${file#$skill_dir/}"
    if [ "$rel_path" != "SKILL.md" ]; then
      if [ -z "$additional_files" ]; then
        additional_files="\"$rel_path\""
      else
        additional_files="$additional_files, \"$rel_path\""
      fi
    fi
  done < <(find "$skill_dir" -type f ! -name ".*" ! -name ".DS_Store" -print0 2>/dev/null | sort -z)

  # Build files array
  if [ -n "$additional_files" ]; then
    files_json="[\"SKILL.md\", $additional_files]"
  else
    files_json="[]"
  fi

  if [ "$first" = true ]; then
    first=false
  else
    echo ',' >> "$OUTPUT_FILE"
  fi

  # Write skill entry (without trailing newline for comma handling)
  printf '    {\n' >> "$OUTPUT_FILE"
  printf '      "name": "%s",\n' "$name" >> "$OUTPUT_FILE"
  printf '      "stack": "%s",\n' "$stack" >> "$OUTPUT_FILE"
  printf '      "description": "%s",\n' "$description" >> "$OUTPUT_FILE"
  printf '      "path": "%s",\n' "$path" >> "$OUTPUT_FILE"
  printf '      "files": %s\n' "$files_json" >> "$OUTPUT_FILE"
  printf '    }' >> "$OUTPUT_FILE"

  skill_count=$((skill_count + 1))
  echo "  Found: $stack/$name"
done

# Close JSON
echo '' >> "$OUTPUT_FILE"
echo '  ]' >> "$OUTPUT_FILE"
echo '}' >> "$OUTPUT_FILE"

echo ""
echo "Generated $OUTPUT_FILE with $skill_count skill(s)"
