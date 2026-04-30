"""
CLI interface for Lyrics Engine using Typer.
"""

import typer
from typing import Optional, List

from lyrics.rhyme_engine import RhymeEngine
from lyrics.syllable_counter import SyllableCounter

app = typer.Typer(help="Lyrics Engine - Standalone lyrics writing tool")
rhyme_engine = RhymeEngine()
syllable_counter = SyllableCounter()


@app.command()
def rhyme(
    word: str = typer.Argument(..., help="Word to find rhymes for"),
    max_results: int = typer.Option(10, "--max", "-m", help="Maximum number of results"),
    slant: bool = typer.Option(True, "--no-slant", help="Include slant rhymes"),
) -> None:
    """
    Find rhymes for a word.
    
    Example:
        lyrics rhyme "love" --max 5
    """
    rhymes = rhyme_engine.find_rhymes(word, max_results=max_results, include_slant=slant)
    
    if rhymes["perfect"]:
        typer.echo("\nPerfect rhymes:")
        for r in rhymes["perfect"]:
            typer.echo(f"  {r}")
    
    if slant and rhymes["slant"]:
        typer.echo("\nSlant rhymes:")
        for r in rhymes["slant"]:
            typer.echo(f"  {r}")
    
    if rhymes["multi_syllable"]:
        typer.echo("\nMulti-syllable rhymes:")
        for r in rhymes["multi_syllable"]:
            typer.echo(f"  {r}")
    
    if not any(rhymes.values()):
        typer.echo(f"No rhymes found for '{word}'")


@app.command()
def syllables(
    text: str = typer.Argument(..., help="Text to count syllables in"),
) -> None:
    """
    Count syllables in text.
    
    Example:
        lyrics syllables "beautiful day"
    """
    count = syllable_counter.count_syllables(text)
    words = text.split()
    
    if len(words) > 1:
        word_counts = []
        for word in words:
            word_count = syllable_counter.count_syllables(word)
            word_counts.append(f"{word}: {word_count}")
        
        typer.echo(f"Total syllables: {count}")
        typer.echo(f"Per word: {', '.join(word_counts)}")
    else:
        typer.echo(f"Syllables in '{text}': {count}")


@app.command()
def generate(
    theme: str = typer.Argument(..., help="Theme for lyric generation"),
    lines: int = typer.Option(4, "--lines", "-l", help="Number of lines to generate"),
    rhyme_scheme: str = typer.Option("ABAB", "--rhyme", "-r", help="Rhyme scheme"),
    syllables_per_line: Optional[int] = typer.Option(
        None, "--syllables", "-s", help="Syllables per line (optional)"
    ),
) -> None:
    """
    Generate lyrics with AI.
    
    Example:
        lyrics generate "heartbreak" --lines 4 --rhyme ABAB --syllables 8
    """
    # This is a skeleton - will be implemented in Phase 1
    typer.echo(f"Generating {lines} lines about '{theme}'...")
    typer.echo(f"Rhyme scheme: {rhyme_scheme}")
    if syllables_per_line:
        typer.echo(f"Syllables per line: {syllables_per_line}")
    
    typer.echo("\n[AI generation will be implemented in Phase 1]")
    typer.echo("See ai_generator.py for implementation.")


@app.command()
def new(
    title: str = typer.Argument(..., help="Title of the new song"),
    template: str = typer.Option(
        "pop", "--template", "-t", help="Template type (pop, rap, ballad, punk, singer-songwriter)"
    ),
) -> None:
    """
    Create a new song with a template.
    
    Example:
        lyrics new "My Song" --template pop
    """
    # This is a skeleton - will be implemented in Phase 2
    typer.echo(f"Creating new song: '{title}'")
    typer.echo(f"Template: {template}")
    typer.echo("\n[Song structure templates will be implemented in Phase 2]")
    typer.echo("See structure.py for implementation.")


@app.command()
def analyze(
    file: str = typer.Argument(..., help="File to analyze"),
) -> None:
    """
    Analyze lyrics for flow, rhyme, and structure.
    
    Example:
        lyrics analyze song.txt
    """
    # This is a skeleton - will be implemented in Phase 3
    typer.echo(f"Analyzing: {file}")
    typer.echo("\n[Flow analysis will be implemented in Phase 3]")
    typer.echo("See flow_analyzer.py for implementation.")


@app.command()
def refine(
    file: str = typer.Argument(..., help="File to refine"),
    output: Optional[str] = typer.Option(None, "--output", "-o", help="Output file"),
) -> None:
    """
    Refine lyrics with AI suggestions.
    
    Example:
        lyrics refine song.txt --output refined.txt
    """
    # This is a skeleton - will be implemented in Phase 3
    typer.echo(f"Refining: {file}")
    if output:
        typer.echo(f"Output: {output}")
    
    typer.echo("\n[Refinement mode will be implemented in Phase 3]")
    typer.echo("See flow_analyzer.py and ai_generator.py for implementation.")


@app.command()
def export(
    file: str = typer.Argument(..., help="File to export"),
    format: str = typer.Option("txt", "--format", "-f", help="Export format (txt, json, pdf)"),
) -> None:
    """
    Export lyrics to different formats.
    
    Example:
        lyrics export song.txt --format pdf
    """
    # This is a skeleton - will be implemented in Phase 4
    typer.echo(f"Exporting {file} to {format.upper()} format")
    typer.echo("\n[Export functionality will be implemented in Phase 4]")
    typer.echo("See export.py for implementation.")


@app.command()
def version() -> None:
    """Show the version of Lyrics Engine."""
    from lyrics import __version__
    typer.echo(f"Lyrics Engine v{__version__}")


if __name__ == "__main__":
    app()