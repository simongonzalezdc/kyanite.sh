# Contributing to Lyrics Engine

Thank you for considering contributing to Lyrics Engine!

## Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Pastorsimon1798/lyrics-engine.git
   cd lyrics-engine
   ```

2. **Create a virtual environment:**
   ```bash
   uv venv
   source .venv/bin/activate
   ```

3. **Install dependencies:**
   ```bash
   uv pip install -e ".[dev]"
   ```

4. **Download spaCy model (if using NLP features):**
   ```bash
   python -m spacy download en_core_web_sm
   ```

## Code Style

This project uses:
- **Black** for code formatting (line length: 88)
- **isort** for import sorting
- **flake8** for linting
- **mypy** for type checking

Run all checks before committing:
```bash
black lyrics/ tests/
isort lyrics/ tests/
flake8 lyrics/ tests/
mypy lyrics/
```

## Testing

Run tests with:
```bash
pytest
pytest --cov=lyrics  # With coverage
```

## Pull Request Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Add tests for new features
   - Update documentation
   - Follow code style guidelines

3. **Run quality checks:**
   ```bash
   black lyrics/ tests/
   isort lyrics/ tests/
   flake8 lyrics/ tests/
   pytest
   ```

4. **Commit with clear messages:**
   ```bash
   git commit -m "feat: add rhyme scoring algorithm"
   ```

5. **Push and create PR:**
   ```bash
   git push origin feature/your-feature-name
   ```

## Commit Message Format

Use conventional commits:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test additions/changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

## Questions?

Open an issue or reach out to the maintainers.

Thank you for contributing! 🎵
