# syntax.sh - Storage Specification

**Version:** 1.0
**Date:** November 2025
**For:** Implementation Reference

---

## Overview

This document specifies the exact file format and storage structure for syntax.sh projects. All implementations must adhere to these specifications to ensure compatibility and data integrity.

---

## Directory Structure

```
~/.local/share/syntax/projects/
└── {project-id}/
    ├── metadata.yaml                 # Project metadata
    ├── characters/
    │   ├── char_abc123def4567890.md
    │   ├── char_xyz789abc1234567.md
    │   └── ...
    ├── locations/
    │   ├── loc_def456ghi7890123.md
    │   ├── loc_pqr987stu6543210.md
    │   └── ...
    ├── scenes/
    │   ├── ch01_sc01.md
    │   ├── ch01_sc02.md
    │   ├── ch02_sc01.md
    │   └── ...
    ├── outline/
    │   └── outline.yaml
    ├── stats/
    │   └── sessions.yaml
    ├── exports/
    │   └── (generated files)
    └── .backups/
        └── (automatic backups)
```

---

## ID Generation

### Character IDs

**Format:** `char_` + 16 hexadecimal characters

**Generation:**
```go
import (
    "crypto/rand"
    "encoding/hex"
)

func GenerateCharacterID() string {
    bytes := make([]byte, 8)  // 8 bytes = 16 hex chars
    rand.Read(bytes)
    return "char_" + hex.EncodeToString(bytes)
}

// Example: char_abc123def4567890
```

**Collision Handling:**
- Check if file exists before creating
- Regenerate if collision detected (astronomically rare)
- Maximum 3 retries before error

### Location IDs

**Format:** `loc_` + 16 hexadecimal characters

**Generation:** Same as character IDs, different prefix

### Scene IDs

**Format:** `ch{chapter}_sc{scene}` (zero-padded to 2 digits)

**Examples:**
- `ch01_sc01.md` (Chapter 1, Scene 1)
- `ch01_sc10.md` (Chapter 1, Scene 10)
- `ch15_sc03.md` (Chapter 15, Scene 3)

**Scene ID Generation:**
```go
func GenerateSceneID(chapter, scene int) string {
    return fmt.Sprintf("ch%02d_sc%02d", chapter, scene)
}
```

---

## Project Metadata

**File:** `metadata.yaml`

**Format:**
```yaml
# Schema version (for future migrations)
schema_version: "1.0"

# Creation info
created_at: "2025-11-10T14:30:00Z"
created_with: "syntax.sh v1.0.0"
last_modified: "2025-11-11T09:15:00Z"

# Project details
id: "proj_abc123def4567890"
title: "The Great Adventure"
author: "Jane Doe"
genre: "fantasy"
status: "draft"  # draft | revising | complete

# Writing goals
target_word_count: 80000
daily_word_goal: 1000

# Statistics (cached for performance)
total_words: 42350
total_scenes: 25
total_characters: 12
total_locations: 8

# Session tracking
current_streak: 7  # days
total_sessions: 42
total_time_seconds: 126000  # 35 hours
```

**Validation Rules:**
- `schema_version`: Required, must be valid semver
- `id`: Required, unique project identifier
- `title`: Required, max 200 characters
- `author`: Optional, max 100 characters
- `genre`: Optional, predefined list or free text
- `status`: Must be one of: draft, revising, complete
- Timestamps: ISO 8601 format (UTC)

---

## Character File Format

**File:** `characters/{character-id}.md`

**Structure:**
```markdown
---
# Required fields
id: char_abc123def4567890
name: "Jane Doe"
created_at: 2025-11-10T14:30:00Z
updated_at: 2025-11-11T09:15:00Z

# Optional fields
aliases: ["JD", "The Detective", "Jay"]
role: "protagonist"  # protagonist | antagonist | supporting | minor
age: 34
occupation: "Homicide Detective"
appearance: "Tall, athletic build, short dark hair, piercing green eyes"
background: "Former military, joined police force after discharge"
arc: "Learning to trust others and work as part of a team"

# Relationships (array of objects)
relationships:
  - character_id: char_xyz789abc1234567
    type: "rival"
    tension: "high"  # low | medium | high
    notes: "Competing for promotion to captain"
  - character_id: char_def456ghi7890123
    type: "mentor"
    tension: "low"
    notes: "Former training officer, now retired"
---

# Jane Doe - Character Biography

Jane grew up in a military family, moving from base to base every few years...

## Personality Traits

- **Determined:** Never gives up on a case
- **Guarded:** Struggles with vulnerability after military trauma
- **Loyal:** Fiercely protective of those she trusts
- **Analytical:** Approaches problems methodically

## Character Development

**Act 1:** Isolated and self-reliant, refuses help from colleagues

**Act 2:** Begins to open up through partnership with Bob Wilson

**Act 3:** Learns to trust her team, becomes a leader

## Key Scenes

- Chapter 1, Scene 3: Introduction during crime scene investigation
- Chapter 5, Scene 2: Confrontation with rival John Smith
- Chapter 12, Scene 7: Breakthrough moment of vulnerability
```

**Parsing:**
```go
import "github.com/adrg/frontmatter"

type Character struct {
    ID            string        `yaml:"id"`
    Name          string        `yaml:"name"`
    Aliases       []string      `yaml:"aliases,omitempty"`
    Role          string        `yaml:"role,omitempty"`
    Age           int           `yaml:"age,omitempty"`
    Occupation    string        `yaml:"occupation,omitempty"`
    Appearance    string        `yaml:"appearance,omitempty"`
    Background    string        `yaml:"background,omitempty"`
    Arc           string        `yaml:"arc,omitempty"`
    Relationships []Relationship `yaml:"relationships,omitempty"`
    CreatedAt     time.Time     `yaml:"created_at"`
    UpdatedAt     time.Time     `yaml:"updated_at"`
    Bio           string        // Markdown content after frontmatter
}

type Relationship struct {
    CharacterID string `yaml:"character_id"`
    Type        string `yaml:"type"`
    Tension     string `yaml:"tension"`
    Notes       string `yaml:"notes,omitempty"`
}

func LoadCharacter(path string) (*Character, error) {
    var char Character
    file, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    rest, err := frontmatter.Parse(bytes.NewReader(file), &char)
    if err != nil {
        return nil, err
    }

    char.Bio = string(rest)
    return &char, nil
}
```

---

## Location File Format

**File:** `locations/{location-id}.md`

**Structure:**
```markdown
---
id: loc_def456ghi7890123
name: "The Rusty Anchor Tavern"
type: "tavern"  # city | tavern | castle | forest | planet | etc
created_at: 2025-11-10T15:00:00Z
updated_at: 2025-11-11T10:00:00Z

# Optional fields
region: "Old Port District"
climate: "Temperate coastal"
population: 50  # Average occupancy
significance: "Major meeting point for investigators"

# Connected locations
connections:
  - location_id: loc_pqr987stu6543210
    connection_type: "road"
    distance: "2 miles"
  - location_id: loc_xyz123abc4567890
    connection_type: "alley"
    distance: "500 feet"
---

# The Rusty Anchor Tavern

A weather-beaten establishment near the harbor, frequented by sailors, dock workers, and those seeking information in the seedier parts of town.

## Physical Description

- Two-story wooden building with faded paint
- Creaky floor boards and dim lighting
- Bar along the east wall
- 10-12 tables scattered throughout
- Back room for private conversations

## Atmosphere

Smoky, loud, and rough around the edges. The kind of place where questions aren't asked and secrets are currency.

## Important Events

- Chapter 2, Scene 1: Jane meets her informant
- Chapter 8, Scene 4: Confrontation with smugglers
- Chapter 14, Scene 2: Final piece of evidence discovered

## Notable NPCs

- Barkeep: Old Tom (gruff but knows everything)
- Regular: Sarah the Sailor (potential ally)
```

---

## Scene File Format

**File:** `scenes/ch{XX}_sc{YY}.md`

**Structure:**
```markdown
---
id: ch01_sc01
chapter: 1
scene_number: 1
name: "The Crime Scene"
created_at: 2025-11-10T16:00:00Z
updated_at: 2025-11-11T11:00:00Z

# Scene metadata
pov_character: char_abc123def4567890  # Jane Doe
location: loc_def456ghi7890123        # The Rusty Anchor
time_of_day: "morning"                # morning | afternoon | evening | night
weather: "overcast"

# Characters present in scene
characters:
  - char_abc123def4567890  # Jane Doe
  - char_xyz789abc1234567  # John Smith

# Plot elements
plot_points:
  - "Introduction of protagonist"
  - "Discovery of first clue"
  - "Establishment of mystery"

# Status
status: "done"  # draft | revising | done
word_count: 1847

# Notes for author
notes: "Need to revise the dialogue in the middle section"
---

# Chapter 1, Scene 1: The Crime Scene

The rain had stopped an hour ago, but the cobblestones still glistened...

[Scene content continues...]
```

**Word Count Calculation:**
```go
func CalculateWordCount(text string) int {
    // Remove markdown syntax
    cleaned := stripMarkdown(text)

    // Split by whitespace
    words := strings.Fields(cleaned)

    return len(words)
}

// Update scene metadata on save
func (s *Scene) UpdateContent(content string) {
    s.Content = content
    s.WordCount = CalculateWordCount(content)
    s.UpdatedAt = time.Now().UTC()
}
```

---

## Outline File Format

**File:** `outline/outline.yaml`

**Structure:**
```yaml
structure: "three-act"  # three-act | hero-journey | custom

acts:
  - number: 1
    name: "Setup"
    goal: "Introduce characters and establish mystery"
    beats:
      - id: "beat_001"
        number: 1
        name: "Opening Image"
        status: "done"  # todo | active | done
        scene_ref: "ch01_sc01"
        notes: "Crime scene discovery"

      - id: "beat_002"
        number: 2
        name: "Theme Stated"
        status: "done"
        scene_ref: "ch01_sc03"
        notes: "Trust issues introduced"

      - id: "beat_003"
        number: 3
        name: "Setup"
        status: "active"
        scene_ref: "ch02_sc01"
        notes: "Establish detective's methods"

  - number: 2
    name: "Confrontation"
    goal: "Complications and rising tension"
    beats:
      - id: "beat_004"
        number: 4
        name: "Catalyst"
        status: "todo"
        scene_ref: ""
        notes: "Second murder changes everything"

      # ... more beats

  - number: 3
    name: "Resolution"
    goal: "Climax and resolution"
    beats:
      - id: "beat_012"
        number: 12
        name: "Finale"
        status: "todo"
        scene_ref: ""
        notes: "Confrontation with killer"
```

---

## Session Statistics

**File:** `stats/sessions.yaml`

**Structure:**
```yaml
# Session history
sessions:
  - session_id: "sess_001"
    start_time: "2025-11-01T09:00:00Z"
    end_time: "2025-11-01T10:30:00Z"
    words_written: 847
    scenes_worked:
      - ch01_sc01
      - ch01_sc02

  - session_id: "sess_002"
    start_time: "2025-11-02T14:00:00Z"
    end_time: "2025-11-02T15:15:00Z"
    words_written: 623
    scenes_worked:
      - ch01_sc02

# Daily statistics
daily_stats:
  "2025-11-01": 847
  "2025-11-02": 623
  "2025-11-03": 1205

# Streaks
current_streak: 7      # consecutive days with writing
longest_streak: 14     # all-time best
last_write_date: "2025-11-11"

# Goals tracking
goals:
  daily_target: 1000
  project_target: 80000
  deadline: "2026-03-01"
```

---

## Backup Strategy

**Directory:** `.backups/`

**Backup Triggers:**
1. Before any data migration
2. Before major version update
3. On user request (`syntax backup`)

**Backup Format:**
```
.backups/
├── backup_20251110_143000.tar.gz
├── backup_20251111_091500.tar.gz
└── backup_migration_v1_to_v2.tar.gz
```

**Backup Contents:**
- All project files (metadata, characters, locations, scenes, outline)
- Exclude: exports/, .backups/ (avoid recursion)

**Retention:**
- Keep last 5 automatic backups
- Keep all migration backups
- User-requested backups never auto-deleted

---

## Data Validation

### On Project Load

```go
func ValidateProject(projectDir string) error {
    // Check metadata exists
    metadataPath := filepath.Join(projectDir, "metadata.yaml")
    if !fileExists(metadataPath) {
        return errors.New("missing metadata.yaml")
    }

    // Parse and validate metadata
    metadata, err := LoadMetadata(metadataPath)
    if err != nil {
        return fmt.Errorf("invalid metadata: %w", err)
    }

    // Validate schema version
    if !isSupportedVersion(metadata.SchemaVersion) {
        return fmt.Errorf("unsupported schema version: %s", metadata.SchemaVersion)
    }

    // Validate required directories exist
    requiredDirs := []string{"characters", "locations", "scenes", "outline"}
    for _, dir := range requiredDirs {
        dirPath := filepath.Join(projectDir, dir)
        if !dirExists(dirPath) {
            // Create missing directory
            os.MkdirAll(dirPath, 0700)
        }
    }

    return nil
}
```

### On File Save

```go
func SaveCharacter(char *Character, projectDir string) error {
    // Validate required fields
    if char.ID == "" {
        return errors.New("character ID required")
    }
    if char.Name == "" {
        return errors.New("character name required")
    }

    // Sanitize user input
    char.Name = SanitizeInput(char.Name, 200)
    char.Occupation = SanitizeInput(char.Occupation, 100)

    // Set timestamps
    if char.CreatedAt.IsZero() {
        char.CreatedAt = time.Now().UTC()
    }
    char.UpdatedAt = time.Now().UTC()

    // Serialize to YAML + Markdown
    data, err := SerializeCharacter(char)
    if err != nil {
        return err
    }

    // Write to file
    path := filepath.Join(projectDir, "characters", char.ID+".md")
    return os.WriteFile(path, data, 0600)
}
```

---

## Data Integrity

### Checksums

**Purpose:** Detect file corruption

**Implementation:**
```go
import "crypto/sha256"

// Calculate checksum for a file
func CalculateChecksum(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }

    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:]), nil
}

// Store checksums in metadata
type Metadata struct {
    // ... other fields
    Checksums map[string]string `yaml:"checksums"`
}

// Validate on load
func ValidateChecksums(projectDir string, metadata *Metadata) error {
    for filePath, expectedChecksum := range metadata.Checksums {
        actualChecksum, err := CalculateChecksum(filePath)
        if err != nil {
            return fmt.Errorf("failed to read %s: %w", filePath, err)
        }

        if actualChecksum != expectedChecksum {
            return fmt.Errorf("corruption detected in %s", filePath)
        }
    }
    return nil
}
```

**Checksum Strategy (v1.0):**
- Not implemented in v1.0 (performance concern)
- Planned for v1.1
- User opt-in via config

---

## Migration Between Versions

### Schema Version v1.0 → v1.1 (Example)

**Change:** Add spell-check fields

**Migration:**
```go
func MigrateV1_0_to_V1_1(projectDir string) error {
    // Load metadata
    metadata, err := LoadMetadata(projectDir)
    if err != nil {
        return err
    }

    // Check current version
    if metadata.SchemaVersion != "1.0" {
        return errors.New("not a v1.0 project")
    }

    // Create backup before migration
    if err := CreateBackup(projectDir, "migration_v1_to_v1.1"); err != nil {
        return err
    }

    // Add new fields with defaults
    metadata.SchemaVersion = "1.1"
    metadata.SpellCheckEnabled = true
    metadata.SpellCheckLanguage = "en_US"

    // Save updated metadata
    if err := SaveMetadata(projectDir, metadata); err != nil {
        return err
    }

    log.Println("Migration complete: v1.0 → v1.1")
    return nil
}
```

---

## Platform-Specific Considerations

### Windows

- File paths use backslashes (handled by `filepath.Join`)
- Line endings: CRLF (Go handles automatically)
- File permissions: ACLs instead of Unix permissions

### macOS

- File paths are case-insensitive (but case-preserving)
- Extended attributes may be present

### Linux

- File paths are case-sensitive
- Watch for NFS/network filesystem issues

### All Platforms

- Always use `filepath.Join()` for path construction
- Always use `os.WriteFile()` with appropriate permissions
- Use `github.com/adrg/xdg` for config directories

---

## File Size Limits

| File Type | Max Size | Reason |
|-----------|----------|--------|
| Scene | 10MB | ~500,000 words (unrealistic for single scene) |
| Character | 1MB | ~50,000 words (bio) |
| Location | 1MB | ~50,000 words (description) |
| Metadata | 100KB | Structured data only |
| Outline | 500KB | Reasonable for large projects |

**Enforcement:**
```go
const MaxSceneSize = 10 * 1024 * 1024  // 10MB

func LoadScene(path string) (*Scene, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if info.Size() > MaxSceneSize {
        return nil, errors.New("scene file too large")
    }

    // Continue loading...
}
```

---

## Atomic Operations

**Problem:** Prevent data loss during save

**Solution:** Write-then-rename pattern

```go
func AtomicWriteFile(path string, data []byte, perm os.FileMode) error {
    // Write to temporary file first
    tmpPath := path + ".tmp"
    if err := os.WriteFile(tmpPath, data, perm); err != nil {
        return err
    }

    // Atomic rename (overwrites existing file)
    if err := os.Rename(tmpPath, path); err != nil {
        os.Remove(tmpPath)  // Clean up temp file
        return err
    }

    return nil
}
```

**Benefits:**
- If write fails, original file unchanged
- If power loss during write, temp file discarded
- Atomic rename ensures no partial writes visible

---

## Summary

This specification defines:
- ✅ Exact directory structure
- ✅ File formats (YAML frontmatter + Markdown)
- ✅ ID generation schemes
- ✅ Data validation rules
- ✅ Backup strategies
- ✅ Migration procedures
- ✅ Platform compatibility

**Implementation Checklist:**
- [ ] Create directory structure on project init
- [ ] Generate IDs using crypto/rand
- [ ] Parse files using frontmatter library
- [ ] Validate required fields on save/load
- [ ] Implement atomic file writes
- [ ] Create backups before migrations
- [ ] Use cross-platform paths (filepath.Join)

---

**Reference Implementation:** See `internal/storage/` in source code.
