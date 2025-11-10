'use client';

import { useState } from 'react';
import { Sparkles, Copy, Check } from 'lucide-react';
import { PitchPoint } from '@/lib/types';

interface LyricsAssistantProps {
  pitches: PitchPoint[];
  musicalKey: string | null;
}

interface LyricsVariation {
  lines: string[];
  syllableCounts: number[];
}

export default function LyricsAssistant({ pitches, musicalKey }: LyricsAssistantProps) {
  const [theme, setTheme] = useState('');
  const [mood, setMood] = useState('neutral');
  const [loading, setLoading] = useState(false);
  const [variations, setVariations] = useState<LyricsVariation[]>([]);
  const [selectedVariation, setSelectedVariation] = useState<number | null>(null);
  const [copied, setCopied] = useState(false);

  const calculateSyllableCount = () => {
    // Estimate syllables based on pitch points
    // Rough estimate: 1 syllable per 0.2-0.3 seconds of audio
    const duration = pitches.length > 0 ? pitches[pitches.length - 1].time - pitches[0].time : 0;
    const estimatedSyllables = Math.round(duration / 0.25);
    return Math.max(estimatedSyllables, 8); // Minimum 8 syllables
  };

  const calculatePhraseLengths = () => {
    // Divide into phrases (roughly 4 phrases)
    const totalSyllables = calculateSyllableCount();
    const phrases = 4;
    const baseLength = Math.floor(totalSyllables / phrases);
    const remainder = totalSyllables % phrases;
    
    const lengths: number[] = [];
    for (let i = 0; i < phrases; i++) {
      lengths.push(baseLength + (i < remainder ? 1 : 0));
    }
    return lengths;
  };

  const generateLyrics = async () => {
    if (!theme.trim()) {
      alert('Please enter a theme');
      return;
    }

    setLoading(true);
    try {
      const syllableCount = calculateSyllableCount();
      const phraseLengths = calculatePhraseLengths();

      const response = await fetch('/api/lyrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          syllableCount,
          phraseLengths,
          musicalKey: musicalKey || 'C Major',
          mood,
          theme: theme.trim(),
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to generate lyrics');
      }

      const data = await response.json();
      setVariations(data.variations || []);
      setSelectedVariation(0);
    } catch (error) {
      console.error('Error generating lyrics:', error);
      alert('Failed to generate lyrics. Please check your API configuration.');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (variation: LyricsVariation) => {
    const text = variation.lines.join('\n');
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-6">
      <h2 className="text-xl font-semibold flex items-center gap-2">
        <Sparkles size={24} className="text-secondary-500" />
        AI Lyric Assist
      </h2>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Theme</label>
          <input
            type="text"
            value={theme}
            onChange={(e) => setTheme(e.target.value)}
            placeholder="e.g., summer, love, adventure"
            className="w-full px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-primary-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">Mood</label>
          <select
            value={mood}
            onChange={(e) => setMood(e.target.value)}
            className="w-full px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-primary-500"
          >
            <option value="happy">Happy</option>
            <option value="sad">Sad</option>
            <option value="energetic">Energetic</option>
            <option value="calm">Calm</option>
            <option value="romantic">Romantic</option>
            <option value="neutral">Neutral</option>
          </select>
        </div>

        <button
          onClick={generateLyrics}
          disabled={loading || !theme.trim()}
          className="w-full px-6 py-3 bg-secondary-500 hover:bg-secondary-600 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {loading ? (
            <>
              <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
              Generating...
            </>
          ) : (
            <>
              <Sparkles size={20} />
              Generate Lyrics
            </>
          )}
        </button>
      </div>

      {variations.length > 0 && (
        <div className="space-y-4">
          <h3 className="font-medium">Generated Variations</h3>
          {variations.map((variation, index) => (
            <div
              key={index}
              className={`p-4 rounded-lg border ${
                selectedVariation === index
                  ? 'bg-primary-500/20 border-primary-500'
                  : 'bg-gray-800 border-gray-700'
              }`}
            >
              <div className="flex items-start justify-between mb-2">
                <span className="text-sm text-gray-400">Variation {index + 1}</span>
                <button
                  onClick={() => copyToClipboard(variation)}
                  className="p-1 hover:bg-gray-700 rounded transition-colors"
                >
                  {copied ? (
                    <Check size={16} className="text-green-400" />
                  ) : (
                    <Copy size={16} className="text-gray-400" />
                  )}
                </button>
              </div>
              <div className="space-y-1">
                {variation.lines.map((line, lineIndex) => (
                  <p key={lineIndex} className="text-white">
                    {line}
                  </p>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

