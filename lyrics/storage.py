"""
Storage - SQLite interface for songs, versions, and templates.

This module will be implemented in Phase 1.
"""

import sqlite3
from typing import Dict, List, Optional
from pathlib import Path


class LyricsStorage:
    """SQLite storage for lyrics and metadata."""
    
    def __init__(self, db_path: str = "data/data.sqlite") -> None:
        """
        Initialize storage.
        
        Args:
            db_path: Path to SQLite database
        """
        self.db_path = Path(db_path)
        self.db_path.parent.mkdir(parents=True, exist_ok=True)
        # TODO: Initialize database schema
    
    def save_song(
        self,
        title: str,
        lyrics: str,
        metadata: Optional[Dict] = None
    ) -> int:
        """
        Save a song to the database.
        
        Args:
            title: Song title
            lyrics: Lyrics text
            metadata: Optional metadata (rhyme scheme, etc.)
            
        Returns:
            Song ID
        """
        # TODO: Implement in Phase 1
        raise NotImplementedError("Storage will be implemented in Phase 1")
    
    def load_song(self, song_id: int) -> Dict:
        """
        Load a song from the database.
        
        Args:
            song_id: Song ID
            
        Returns:
            Dictionary with song data
        """
        # TODO: Implement in Phase 1
        raise NotImplementedError("Storage will be implemented in Phase 1")
    
    def save_version(
        self,
        song_id: int,
        lyrics: str,
        change_description: Optional[str] = None
    ) -> int:
        """
        Save a new version of a song.
        
        Args:
            song_id: Original song ID
            lyrics: Updated lyrics
            change_description: Description of changes
            
        Returns:
            Version ID
        """
        # TODO: Implement in Phase 4
        raise NotImplementedError("Version control will be implemented in Phase 4")
    
    def list_songs(self, limit: int = 10) -> List[Dict]:
        """
        List recent songs.
        
        Args:
            limit: Maximum number of songs to return
            
        Returns:
            List of song dictionaries
        """
        # TODO: Implement in Phase 1
        raise NotImplementedError("Storage will be implemented in Phase 1")
