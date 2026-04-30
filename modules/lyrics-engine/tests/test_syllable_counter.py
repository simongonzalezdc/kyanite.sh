"""
Tests for syllable_counter module.
"""

import pytest
from lyrics.syllable_counter import SyllableCounter


@pytest.fixture
def counter():
    """Fixture for SyllableCounter instance."""
    return SyllableCounter()


def test_count_syllables_single_word(counter):
    """Test syllable counting for single word."""
    assert counter.count_syllables("hello") == 2
    assert counter.count_syllables("beautiful") == 3
    assert counter.count_syllables("fire") >= 1  # Edge case: 1 or 2


def test_count_syllables_phrase(counter):
    """Test syllable counting for phrase."""
    count = counter.count_syllables("beautiful day")
    assert count >= 4  # beautiful (3) + day (1)


def test_get_syllable_pattern(counter):
    """Test syllable pattern extraction."""
    pattern = counter.get_syllable_pattern("hello world")
    assert len(pattern) == 2
    assert pattern[0] == 2  # hello
    assert pattern[1] == 1  # world


def test_check_constraints_pass(counter):
    """Test constraint checking (passing)."""
    meets, actual, msg = counter.check_constraints("hello world", target_syllables=3, tolerance=1)
    assert meets is True


def test_check_constraints_fail(counter):
    """Test constraint checking (failing)."""
    meets, actual, msg = counter.check_constraints("hello", target_syllables=10, tolerance=1)
    assert meets is False


# TODO: Add more tests in Phase 1
