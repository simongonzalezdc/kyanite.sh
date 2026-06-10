# Lyrics Engine

A standalone lyrics writing tool with rhyme engine, syllable counter, AI generation, and song structure templates.

## Features

- **Rhyme Engine**: Find perfect, slant, and multi-syllable rhymes
- **Syllable Counter**: Accurate syllable counting with stress pattern detection
- **AI Generation**: Generate lyrics with rhyme scheme and syllable constraints
- **Song Structure Templates**: Pop, rap, ballad, punk, and singer-songwriter templates
- **Flow Analysis**: Detect stress patterns and rate "singability"
- **Version History**: Track every edit with rollback capability
- **Collaborative Mode**: Simon writes → Liam refines → iterate workflow

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/Pastorsimon1798/lyrics-engine.git
cd lyrics-engine

# Install with uv (recommended)
uv venv
source .venv/bin/activate
uv pip install -e .

# Or install with pip
pip install -e .
```

### Basic Usage

```bash
# Find rhymes for a word
lyrics rhyme "love"

# Count syllables in a line
lyrics syllables "beautiful day"

# Generate lyrics with a theme
lyrics generate "heartbreak" --lines 4 --rhyme ABAB

# Create a new song with a template
lyrics new "My Song" --template pop
```

### Advanced Usage

```bash
# Analyze a song's flow and structure
lyrics analyze song.txt

# Refine lyrics with AI suggestions
lyrics refine song.txt

# Export to different formats
lyrics export song.txt --format pdf
lyrics export song.txt --format json
```

## Project Structure

```
lyrics-engine/
├── lyrics/                    # Main package
│   ├── __init__.py
│   ├── cli.py                # CLI commands (Typer)
│   ├── rhyme_engine.py       # Rhyme finder
│   ├── syllable_counter.py   # Syllable counting
│   ├── flow_analyzer.py      # Flow analysis
│   ├── ai_generator.py       # AI integration
│   ├── structure.py          # Song structure templates
│   ├── storage.py            # SQLite interface
│   └── export.py             # Export formats
├── data/                     # Data files
│   ├── data.sqlite          # Songs, versions, templates
│   └── phonetic_dict.txt    # Rhyme dictionary
├── drafts/                   # Work in progress
├── final/                   # Completed songs
└── exports/                 # Exported files
```

## Development

### Setting up Development Environment

```bash
# Clone the repository
git clone https://github.com/Pastorsimon1798/lyrics-engine.git
cd lyrics-engine

# Create virtual environment
uv venv
source .venv/bin/activate

# Install with dev dependencies
uv pip install -e ".[dev]"

# Run tests
pytest
```

### Running Tests

```bash
# Run all tests
pytest

# Run tests with coverage
pytest --cov=lyrics

# Run specific test module
pytest tests/test_rhyme_engine.py -v
```

### Code Style

This project uses:
- **Black** for code formatting
- **isort** for import sorting
- **flake8** for linting
- **mypy** for type checking

Run code quality checks:
```bash
black lyrics/ tests/
isort lyrics/ tests/
flake8 lyrics/ tests/
mypy lyrics/
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [pronouncing](https://github.com/aparrish/pronouncing) for rhyme detection
- [pyphen](https://github.com/Kozea/Pyphen) for syllable counting
- [spaCy](https://spacy.io/) for NLP capabilities
- [Typer](https://typer.tiangolo.com/) for beautiful CLI interfaces