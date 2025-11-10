'use client';

import { useState, useRef } from 'react';
import { Sparkles, Wand2, Lightbulb, Copy, Check, ArrowRight } from 'lucide-react';
import { PitchPoint } from '@/lib/types';

interface LyricsAssistantProps {
  pitches: PitchPoint[];
  musicalKey: string | null;
}

interface Suggestion {
  text: string;
  type: 'completion' | 'improvement' | 'alternative';
}

export default function LyricsAssistant({ pitches, musicalKey }: LyricsAssistantProps) {
  const [lyrics, setLyrics] = useState('');
  const [selectedLine, setSelectedLine] = useState<number | null>(null);
  const [selectedText, setSelectedText] = useState('');
  const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
  const [loading, setLoading] = useState(false);
  const [suggestionType, setSuggestionType] = useState<'completion' | 'improvement' | 'alternative'>('completion');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const getCurrentLine = (): string => {
    const lines = lyrics.split('\n');
    if (selectedLine !== null && selectedLine < lines.length) {
      return lines[selectedLine];
    }
    // Get the last line
    return lines[lines.length - 1] || '';
  };

  const getCursorPosition = (): { line: number; column: number } => {
    if (!textareaRef.current) return { line: 0, column: 0 };
    
    const textarea = textareaRef.current;
    const text = textarea.value;
    const cursorPos = textarea.selectionStart;
    
    const textBeforeCursor = text.substring(0, cursorPos);
    const lines = textBeforeCursor.split('\n');
    
    return {
      line: lines.length - 1,
      column: lines[lines.length - 1].length
    };
  };

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setLyrics(e.target.value);
    setSuggestions([]);
  };

  const handleTextSelection = () => {
    if (!textareaRef.current) return;
    
    const textarea = textareaRef.current;
    const selected = textarea.value.substring(textarea.selectionStart, textarea.selectionEnd);
    
    if (selected.trim()) {
      setSelectedText(selected);
      const cursorPos = getCursorPosition();
      setSelectedLine(cursorPos.line);
    } else {
      setSelectedText('');
      setSelectedLine(null);
    }
  };

  const suggestCompletion = async () => {
    const currentLine = getCurrentLine();
    if (!currentLine.trim()) {
      alert('Please start typing a line first');
      return;
    }

    setLoading(true);
    setSuggestionType('completion');
    
    try {
      const response = await fetch('/api/lyrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          type: 'completion',
          currentLine: currentLine.trim(),
          context: lyrics,
          musicalKey: musicalKey || 'C Major',
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to get suggestions');
      }

      const data = await response.json();
      setSuggestions(data.suggestions || []);
    } catch (error) {
      console.error('Error getting suggestions:', error);
      alert('Failed to get suggestions. Please check your API configuration.');
    } finally {
      setLoading(false);
    }
  };

  const suggestImprovement = async () => {
    if (!selectedText.trim()) {
      alert('Please select a line or phrase to improve');
      return;
    }

    setLoading(true);
    setSuggestionType('improvement');
    
    try {
      const response = await fetch('/api/lyrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          type: 'improvement',
          selectedText: selectedText.trim(),
          context: lyrics,
          musicalKey: musicalKey || 'C Major',
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to get suggestions');
      }

      const data = await response.json();
      setSuggestions(data.suggestions || []);
    } catch (error) {
      console.error('Error getting suggestions:', error);
      alert('Failed to get suggestions. Please check your API configuration.');
    } finally {
      setLoading(false);
    }
  };

  const suggestAlternative = async () => {
    if (!selectedText.trim()) {
      alert('Please select a line or phrase to get alternatives');
      return;
    }

    setLoading(true);
    setSuggestionType('alternative');
    
    try {
      const response = await fetch('/api/lyrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          type: 'alternative',
          selectedText: selectedText.trim(),
          context: lyrics,
          musicalKey: musicalKey || 'C Major',
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to get suggestions');
      }

      const data = await response.json();
      setSuggestions(data.suggestions || []);
    } catch (error) {
      console.error('Error getting suggestions:', error);
      alert('Failed to get suggestions. Please check your API configuration.');
    } finally {
      setLoading(false);
    }
  };

  const applySuggestion = (suggestion: string) => {
    if (!textareaRef.current) return;

    const textarea = textareaRef.current;
    const text = textarea.value;
    const cursorPos = textarea.selectionStart;
    const selectionEnd = textarea.selectionEnd;

    if (suggestionType === 'completion') {
      // Insert completion at cursor
      const beforeCursor = text.substring(0, cursorPos);
      const afterCursor = text.substring(cursorPos);
      const newText = beforeCursor + suggestion + afterCursor;
      setLyrics(newText);
      
      // Set cursor position after inserted text
      setTimeout(() => {
        if (textareaRef.current) {
          const newPos = cursorPos + suggestion.length;
          textareaRef.current.setSelectionRange(newPos, newPos);
          textareaRef.current.focus();
        }
      }, 0);
    } else {
      // Replace selected text
      const beforeSelection = text.substring(0, textarea.selectionStart);
      const afterSelection = text.substring(selectionEnd);
      const newText = beforeSelection + suggestion + afterSelection;
      setLyrics(newText);
      
      // Set cursor position after replaced text
      setTimeout(() => {
        if (textareaRef.current) {
          const newPos = beforeSelection.length + suggestion.length;
          textareaRef.current.setSelectionRange(newPos, newPos);
          textareaRef.current.focus();
        }
      }, 0);
    }

    setSuggestions([]);
  };

  const copyLyrics = () => {
    navigator.clipboard.writeText(lyrics);
    alert('Lyrics copied to clipboard!');
  };

  return (
    <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold flex items-center gap-2">
          <Sparkles size={24} className="text-secondary-500" />
          AI Lyric Assistant
        </h2>
        {lyrics && (
          <button
            onClick={copyLyrics}
            className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors text-sm"
          >
            <Copy size={16} />
            Copy
          </button>
        )}
      </div>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Write Your Lyrics</label>
          <textarea
            ref={textareaRef}
            value={lyrics}
            onChange={handleTextChange}
            onSelect={handleTextSelection}
            placeholder="Start writing your lyrics here...&#10;&#10;The AI will help you complete lines, improve phrases, and suggest alternatives."
            className="w-full h-64 px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-primary-500 resize-none font-mono text-sm"
          />
          <p className="text-xs text-gray-400 mt-2">
            {lyrics.split('\n').length} lines • {lyrics.length} characters
          </p>
        </div>

        <div className="flex flex-wrap gap-3">
          <button
            onClick={suggestCompletion}
            disabled={loading || !lyrics.trim()}
            className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
          >
            {loading && suggestionType === 'completion' ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                Suggesting...
              </>
            ) : (
              <>
                <Wand2 size={16} />
                Complete Line
              </>
            )}
          </button>

          <button
            onClick={suggestImprovement}
            disabled={loading || !selectedText.trim()}
            className="flex items-center gap-2 px-4 py-2 bg-secondary-500 hover:bg-secondary-600 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
          >
            {loading && suggestionType === 'improvement' ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                Improving...
              </>
            ) : (
              <>
                <Lightbulb size={16} />
                Improve Selected
              </>
            )}
          </button>

          <button
            onClick={suggestAlternative}
            disabled={loading || !selectedText.trim()}
            className="flex items-center gap-2 px-4 py-2 bg-purple-500 hover:bg-purple-600 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
          >
            {loading && suggestionType === 'alternative' ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                Finding...
              </>
            ) : (
              <>
                <Sparkles size={16} />
                Alternatives
              </>
            )}
          </button>
        </div>

        {selectedText && (
          <div className="p-3 bg-gray-800 rounded-lg border border-gray-700">
            <p className="text-xs text-gray-400 mb-1">Selected text:</p>
            <p className="text-sm text-white font-mono">"{selectedText}"</p>
          </div>
        )}

        {suggestions.length > 0 && (
          <div className="space-y-3">
            <h3 className="font-medium text-sm">
              {suggestionType === 'completion' && 'Completion Suggestions'}
              {suggestionType === 'improvement' && 'Improvement Suggestions'}
              {suggestionType === 'alternative' && 'Alternative Suggestions'}
            </h3>
            <div className="space-y-2">
              {suggestions.map((suggestion, index) => (
                <button
                  key={index}
                  onClick={() => applySuggestion(suggestion.text)}
                  className="w-full p-4 bg-gray-800 border border-gray-700 rounded-lg hover:border-primary-500 hover:bg-gray-750 transition-colors text-left group"
                >
                  <div className="flex items-start justify-between gap-3">
                    <p className="text-white font-mono text-sm flex-1">{suggestion.text}</p>
                    <ArrowRight size={16} className="text-gray-400 group-hover:text-primary-500 transition-colors flex-shrink-0" />
                  </div>
                </button>
              ))}
            </div>
            <button
              onClick={() => setSuggestions([])}
              className="text-sm text-gray-400 hover:text-gray-300 transition-colors"
            >
              Dismiss suggestions
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
