"""
Syllable Counter - Count syllables in words and text.
"""

import pyphen
from typing import List, Tuple, Dict
import re


class SyllableCounter:
    """Counter for syllables in English text."""
    
    def __init__(self, lang: str = "en_US") -> None:
        """
        Initialize the syllable counter.
        
        Args:
            lang: Language code for hyphenation (default: en_US)
        """
        self.dic = pyphen.Pyphen(lang=lang)
        
    def count_syllables(self, text: str) -> int:
        """
        Count syllables in text.
        
        Args:
            text: Text to count syllables in
            
        Returns:
            Number of syllables
        """
        # Clean and normalize text
        text = self._clean_text(text)
        
        if not text:
            return 0
        
        # Split into words
        words = re.findall(r"\b\w+\b", text)
        
        # Count syllables for each word
        total = 0
        for word in words:
            total += self._count_word_syllables(word.lower())
        
        return total
    
    def _count_word_syllables(self, word: str) -> int:
        """
        Count syllables in a single word.
        
        Args:
            word: Single word to count syllables in
            
        Returns:
            Number of syllables
        """
        # Use pyphen to hyphenate and count hyphens
        hyphenated = self.dic.inserted(word)
        if hyphenated:
            # Count syllables: number of hyphens + 1
            return hyphenated.count('-') + 1
        
        # Fallback for words pyphen can't handle
        return self._fallback_syllable_count(word)
    
    def _fallback_syllable_count(self, word: str) -> int:
        """
        Fallback syllable counting using simple vowel pattern matching.
        
        Args:
            word: Single word to count syllables in
            
        Returns:
            Estimated number of syllables
        """
        # Remove non-alphabetic characters
        word = re.sub(r'[^a-z]', '', word.lower())
        
        if len(word) <= 3:
            # Short words are usually 1 syllable
            return 1
        
        # Count vowel groups (basic algorithm)
        # This is a simplified version - not perfect but good enough for many cases
        vowels = "aeiouy"
        
        # Convert to list of characters
        chars = list(word)
        
        # Special cases
        if word.endswith(('es', 'ed')):
            # Remove common suffixes that don't add syllables
            word = word[:-2]
        
        # Count vowel groups
        count = 0
        prev_char_vowel = False
        
        for i, char in enumerate(chars):
            is_vowel = char in vowels
            
            # 'y' as vowel (not at start of word)
            if char == 'y' and i > 0:
                is_vowel = True
            
            if is_vowel:
                if not prev_char_vowel:
                    count += 1
                prev_char_vowel = True
            else:
                prev_char_vowel = False
        
        # Ensure at least 1 syllable
        return max(1, count)
    
    def _clean_text(self, text: str) -> str:
        """
        Clean text for syllable counting.
        
        Args:
            text: Raw text input
            
        Returns:
            Cleaned text
        """
        # Remove punctuation (keep apostrophes for contractions)
        text = re.sub(r'[^\w\s\']', ' ', text)
        # Normalize whitespace
        text = re.sub(r'\s+', ' ', text)
        return text.strip()
    
    def get_syllable_pattern(self, text: str) -> List[int]:
        """
        Get syllable pattern for each word in text.
        
        Args:
            text: Text to analyze
            
        Returns:
            List of syllable counts per word
        """
        # Clean and normalize text
        text = self._clean_text(text)
        
        if not text:
            return []
        
        # Split into words
        words = re.findall(r"\b\w+\b", text)
        
        # Get syllable count for each word
        pattern = []
        for word in words:
            pattern.append(self._count_word_syllables(word.lower()))
        
        return pattern
    
    def detect_meter(self, text: str) -> Dict[str, float]:
        """
        Detect poetic meter in text.
        
        Args:
            text: Text to analyze
            
        Returns:
            Dictionary with meter analysis
        """
        # Get syllable pattern
        pattern = self.get_syllable_pattern(text)
        
        if not pattern:
            return {"error": "No text to analyze"}
        
        # Calculate average syllables per line
        # For now, we'll treat the entire text as one "line"
        avg_syllables = sum(pattern) / max(1, len(pattern))
        
        # Simple meter detection based on syllable count patterns
        meter_scores = {
            "iambic": 0.0,  # x / (unstressed-stressed)
            "trochaic": 0.0,  # / x (stressed-unstressed)
            "anapestic": 0.0,  # x x / (unstressed-unstressed-stressed)
            "dactylic": 0.0,  # / x x (stressed-unstressed-unstressed)
        }
        
        # For now, return basic analysis
        # Full meter detection requires stress pattern analysis (Phase 3)
        return {
            "total_syllables": sum(pattern),
            "words": len(pattern),
            "avg_syllables_per_word": avg_syllables,
            "syllable_pattern": pattern,
            "meter_scores": meter_scores,
            "note": "Full meter detection requires stress pattern analysis (Phase 3)"
        }
    
    def check_constraints(
        self, 
        text: str, 
        target_syllables: int,
        tolerance: int = 1
    ) -> Tuple[bool, int, str]:
        """
        Check if text meets syllable count constraints.
        
        Args:
            text: Text to check
            target_syllables: Target number of syllables
            tolerance: Acceptable deviation from target
            
        Returns:
            Tuple of (meets_constraint, actual_count, message)
        """
        actual_count = self.count_syllables(text)
        difference = abs(actual_count - target_syllables)
        
        if difference <= tolerance:
            return True, actual_count, f"✓ Syllables: {actual_count} (target: {target_syllables})"
        else:
            return False, actual_count, f"✗ Syllables: {actual_count} (target: {target_syllables}, difference: {difference})"
    
    def suggest_corrections(
        self, 
        text: str, 
        target_syllables: int
    ) -> List[str]:
        """
        Suggest corrections to meet syllable count target.
        
        Args:
            text: Text to correct
            target_syllables: Target number of syllables
            
        Returns:
            List of suggested corrections
        """
        suggestions = []
        actual_count = self.count_syllables(text)
        difference = target_syllables - actual_count
        
        if difference == 0:
            suggestions.append("Text already meets syllable target.")
            return suggestions
        
        words = re.findall(r"\b\w+\b", text)
        word_syllables = [self._count_word_syllables(w.lower()) for w in words]
        
        if difference > 0:
            # Need to add syllables
            suggestions.append(f"Need to add {difference} syllable(s). Suggestions:")
            suggestions.append("1. Add descriptive words (adjectives, adverbs)")
            suggestions.append("2. Use multi-syllable synonyms")
            suggestions.append("3. Add articles (the, a) or prepositions")
            
            # Find words that could be expanded
            for i, (word, count) in enumerate(zip(words, word_syllables)):
                if count == 1:
                    # Single-syllable words that could be multi-syllable synonyms
                    synonyms = self._get_multi_syllable_synonyms(word)
                    if synonyms:
                        suggestions.append(f"   - Replace '{word}' with: {', '.join(synonyms[:3])}")
        
        else:
            # Need to remove syllables
            suggestions.append(f"Need to remove {abs(difference)} syllable(s). Suggestions:")
            suggestions.append("1. Remove unnecessary words")
            suggestions.append("2. Use contractions (can't, won't, don't)")
            suggestions.append("3. Use single-syllable synonyms")
            
            # Find words that could be shortened
            for i, (word, count) in enumerate(zip(words, word_syllables)):
                if count >= 2:
                    # Multi-syllable words that could be single-syllable synonyms
                    synonyms = self._get_single_syllable_synonyms(word)
                    if synonyms:
                        suggestions.append(f"   - Replace '{word}' with: {', '.join(synonyms[:3])}")
        
        return suggestions
    
    def _get_multi_syllable_synonyms(self, word: str) -> List[str]:
        """Get multi-syllable synonyms for a word (simplified)."""
        # This is a simplified version - would use a thesaurus in production
        synonyms_dict = {
            "big": ["enormous", "gigantic", "colossal", "massive"],
            "small": ["tiny", "minuscule", "miniature", "petite"],
            "fast": ["rapid", "swift", "speedy", "accelerated"],
            "slow": ["gradual", "leisurely", "unhurried", "deliberate"],
            "happy": ["joyful", "delighted", "ecstatic", "elated"],
            "sad": ["melancholy", "despondent", "disheartened", "mournful"],
        }
        
        return synonyms_dict.get(word.lower(), [])
    
    def _get_single_syllable_synonyms(self, word: str) -> List[str]:
        """Get single-syllable synonyms for a word (simplified)."""
        # This is a simplified version - would use a thesaurus in production
        synonyms_dict = {
            "enormous": ["big", "huge", "vast"],
            "gigantic": ["big", "huge"],
            "colossal": ["big", "huge"],
            "massive": ["big", "huge"],
            "tiny": ["small", "little"],
            "miniature": ["small", "little"],
            "rapid": ["fast", "quick"],
            "swift": ["fast", "quick"],
            "gradual": ["slow"],
            "leisurely": ["slow"],
            "joyful": ["glad", "happy"],
            "delighted": ["glad", "happy"],
            "melancholy": ["sad", "blue"],
            "despondent": ["sad", "low"],
        }
        
        return synonyms_dict.get(word.lower(), [])