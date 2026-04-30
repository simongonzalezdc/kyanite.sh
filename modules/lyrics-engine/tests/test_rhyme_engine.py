"""
Tests for rhyme_engine module.
"""

import pytest
from lyrics.rhyme_engine import RhymeEngine


@pytest.fixture
def rhyme_engine():
    """Fixture for RhymeEngine instance."""
    return RhymeEngine()


def test_find_rhymes_basic(rhyme_engine):
    """Test basic rhyme finding."""
    rhymes = rhyme_engine.find_rhymes("love", max_results=5)
    
    assert "perfect" in rhymes
    assert "slant" in rhymes
    assert "multi_syllable" in rhymes
    assert len(rhymes["perfect"]) > 0


def test_find_rhymes_no_results(rhyme_engine):
    """Test rhyme finding with non-existent word."""
    rhymes = rhyme_engine.find_rhymes("xyzabc", max_results=5)
    
    assert rhymes["perfect"] == []
    assert rhymes["slant"] == []
    assert rhymes["multi_syllable"] == []


def test_rhyme_score_perfect(rhyme_engine):
    """Test rhyme scoring for perfect rhyme."""
    score = rhyme_engine.get_rhyme_score("love", "dove")
    assert score == 1.0


def test_rhyme_score_same_word(rhyme_engine):
    """Test rhyme scoring for same word."""
    score = rhyme_engine.get_rhyme_score("love", "love")
    assert score == 0.0


# TODO: Add more tests in Phase 1
