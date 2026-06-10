"""
Rhyme Engine - Find perfect, slant, and multi-syllable rhymes.
"""

import pronouncing
from typing import List, Dict, Optional, Set
import re


class RhymeEngine:
    """Engine for finding rhymes of different types."""
    
    def __init__(self) -> None:
        """Initialize the rhyme engine."""
        # Cache for performance
        self._cache: Dict[str, Dict[str, List[str]]] = {}
    
    def find_rhymes(
        self, 
        word: str, 
        max_results: int = 10, 
        include_slant: bool = True
    ) -> Dict[str, List[str]]:
        """
        Find rhymes for a word.
        
        Args:
            word: The word to find rhymes for
            max_results: Maximum number of results per category
            include_slant: Whether to include slant rhymes
            
        Returns:
            Dictionary with keys: 'perfect', 'slant', 'multi_syllable'
        """
        # Normalize word
        word = word.lower().strip()
        
        # Check cache
        cache_key = f"{word}_{max_results}_{include_slant}"
        if cache_key in self._cache:
            return self._cache[cache_key]
        
        # Get pronunciations
        pronunciations = pronouncing.phones_for_word(word)
        if not pronunciations:
            return {"perfect": [], "slant": [], "multi_syllable": []}
        
        # Use first pronunciation for rhyme finding
        main_pronunciation = pronunciations[0]
        
        # Find perfect rhymes
        perfect_rhymes = pronouncing.rhymes(word)
        
        # Filter to remove the word itself and limit results
        perfect_rhymes = [r for r in perfect_rhymes if r.lower() != word]
        perfect_rhymes = perfect_rhymes[:max_results]
        
        # Find slant rhymes if requested
        slant_rhymes: List[str] = []
        if include_slant:
            slant_rhymes = self._find_slant_rhymes(word, max_results)
        
        # Find multi-syllable rhymes (words with 2+ syllables that rhyme)
        multi_syllable_rhymes = self._find_multi_syllable_rhymes(word, max_results)
        
        result = {
            "perfect": perfect_rhymes,
            "slant": slant_rhymes,
            "multi_syllable": multi_syllable_rhymes,
        }
        
        # Cache result
        self._cache[cache_key] = result
        return result
    
    def _find_slant_rhymes(self, word: str, max_results: int) -> List[str]:
        """
        Find slant (imperfect) rhymes.
        
        Slant rhymes match on the final vowel sound but not the entire ending.
        """
        # Get pronunciations for the word
        pronunciations = pronouncing.phones_for_word(word)
        if not pronunciations:
            return []
        
        # Extract the vowel sounds from the pronunciation
        main_pronunciation = pronunciations[0]
        vowel_sounds = self._extract_vowel_sounds(main_pronunciation)
        
        if not vowel_sounds:
            return []
        
        # Find words with similar vowel sounds in their endings
        slant_matches: List[str] = []
        all_words = pronouncing.search(".*")  # Get all words in dictionary
        
        for other_word in all_words[:5000]:  # Limit search for performance
            if other_word.lower() == word.lower():
                continue
            
            other_pronunciations = pronouncing.phones_for_word(other_word)
            if not other_pronunciations:
                continue
            
            other_pronunciation = other_pronunciations[0]
            other_vowel_sounds = self._extract_vowel_sounds(other_pronunciation)
            
            # Check if the vowel sounds match (for slant rhyme)
            if other_vowel_sounds and vowel_sounds[-1] == other_vowel_sounds[-1]:
                slant_matches.append(other_word)
                if len(slant_matches) >= max_results:
                    break
        
        return slant_matches
    
    def _find_multi_syllable_rhymes(self, word: str, max_results: int) -> List[str]:
        """
        Find multi-syllable rhymes (words with 2+ syllables that rhyme).
        
        Multi-syllable rhymes are perfect rhymes where both words have
        2 or more syllables.
        """
        # Get syllable count for the word
        syllable_count = pronouncing.syllable_count(word)
        if syllable_count < 2:
            return []
        
        # Find perfect rhymes
        perfect_rhymes = pronouncing.rhymes(word)
        
        # Filter for multi-syllable rhymes
        multi_syllable_rhymes = []
        for rhyme_word in perfect_rhymes:
            if rhyme_word.lower() == word.lower():
                continue
            
            rhyme_syllables = pronouncing.syllable_count(rhyme_word)
            if rhyme_syllables >= 2:
                multi_syllable_rhymes.append(rhyme_word)
            
            if len(multi_syllable_rhymes) >= max_results:
                break
        
        return multi_syllable_rhymes
    
    def _extract_vowel_sounds(self, pronunciation: str) -> List[str]:
        """
        Extract vowel sounds from a pronunciation string.
        
        Pronouncing library uses ARPAbet notation where vowel sounds
        are represented by codes like 'AA', 'AE', 'AH', etc.
        """
        # ARPAbet vowel codes (simplified list)
        vowel_codes = {'AA', 'AE', 'AH', 'AO', 'AW', 'AY', 'EH', 'ER', 'EY', 
                      'IH', 'IY', 'OW', 'OY', 'UH', 'UW'}
        
        sounds = pronunciation.split()
        vowel_sounds = [sound for sound in sounds 
                       if any(sound.startswith(vowel) for vowel in vowel_codes)]
        
        return vowel_sounds
    
    def get_rhyme_score(self, word1: str, word2: str) -> float:
        """
        Calculate a rhyme score between two words (0.0 to 1.0).
        
        1.0 = perfect rhyme
        0.5-0.9 = slant rhyme
        0.0-0.5 = no rhyme
        """
        if word1.lower() == word2.lower():
            return 0.0  # Same word, not a rhyme
        
        # Check for perfect rhyme
        perfect_rhymes = pronouncing.rhymes(word1)
        if word2.lower() in [r.lower() for r in perfect_rhymes]:
            return 1.0
        
        # Check for slant rhyme
        pronunciations1 = pronouncing.phones_for_word(word1)
        pronunciations2 = pronouncing.phones_for_word(word2)
        
        if not pronunciations1 or not pronunciations2:
            return 0.0
        
        pron1 = pronunciations1[0]
        pron2 = pronunciations2[0]
        
        vowel_sounds1 = self._extract_vowel_sounds(pron1)
        vowel_sounds2 = self._extract_vowel_sounds(pron2)
        
        if vowel_sounds1 and vowel_sounds2 and vowel_sounds1[-1] == vowel_sounds2[-1]:
            return 0.7  # Slant rhyme score
        
        return 0.0
    
    def suggest_rhyme(self, word: str, context: Optional[str] = None) -> List[str]:
        """
        Suggest rhymes with optional context awareness.
        
        Args:
            word: The word to find rhymes for
            context: Optional context to influence suggestions
            
        Returns:
            List of suggested rhymes
        """
        rhymes = self.find_rhymes(word, max_results=20, include_slant=True)
        
        # Combine all rhymes
        all_rhymes = []
        all_rhymes.extend(rhymes["perfect"])
        all_rhymes.extend(rhymes["slant"])
        all_rhymes.extend(rhymes["multi_syllable"])
        
        # Remove duplicates
        all_rhymes = list(dict.fromkeys(all_rhymes))
        
        # Simple context filtering (will be enhanced in later phases)
        if context:
            context_words = set(context.lower().split())
            # Prefer rhymes that appear in context
            scored_rhymes = []
            for rhyme in all_rhymes:
                score = 1.0 if rhyme.lower() in context_words else 0.5
                scored_rhymes.append((score, rhyme))
            
            scored_rhymes.sort(key=lambda x: x[0], reverse=True)
            all_rhymes = [r for _, r in scored_rhymes]
        
        return all_rhymes[:10]  # Return top 10