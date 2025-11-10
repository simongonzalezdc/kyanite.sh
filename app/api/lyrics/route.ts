import { NextRequest, NextResponse } from 'next/server';
import OpenAI from 'openai';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { syllableCount, phraseLengths, musicalKey, mood, theme } = body;

    // Get API key and base URL from environment
    const apiKey = process.env.OPENAI_API_KEY;
    const baseURL = process.env.OPENAI_API_BASE_URL || 'https://api.openai.com/v1';
    const model = process.env.OPENAI_MODEL || 'gpt-4o-mini';

    if (!apiKey) {
      return NextResponse.json(
        { error: 'API key not configured' },
        { status: 500 }
      );
    }

    // Initialize OpenAI client with custom base URL for OpenAI-compatible APIs
    const openai = new OpenAI({
      apiKey: apiKey,
      baseURL: baseURL,
    });

    // Create prompt for lyrics generation
    const prompt = `Generate song lyrics that match the following requirements:
- Total syllables: ${syllableCount}
- Phrase lengths (syllables per line): ${phraseLengths.join(', ')}
- Musical key: ${musicalKey}
- Mood: ${mood || 'neutral'}
- Theme: ${theme || 'general'}

Generate 3 variations of lyrics. Each variation should:
1. Match the syllable count exactly
2. Fit the musical rhythm
3. Match the mood and theme
4. Be creative and engaging

Format the response as JSON with this structure:
{
  "variations": [
    {
      "lines": ["line 1", "line 2", "line 3", "line 4"],
      "syllableCounts": [8, 8, 8, 8]
    }
  ]
}`;

    const completion = await openai.chat.completions.create({
      model: model,
      messages: [
        {
          role: 'system',
          content: 'You are a professional songwriter. Generate lyrics that match the exact syllable counts and musical requirements provided.'
        },
        {
          role: 'user',
          content: prompt
        }
      ],
      temperature: 0.8,
      max_tokens: 1000,
    });

    const content = completion.choices[0]?.message?.content || '';
    
    // Try to parse JSON from response
    let lyricsData;
    try {
      // Extract JSON from markdown code blocks if present
      const jsonMatch = content.match(/```json\s*([\s\S]*?)\s*```/) || content.match(/```\s*([\s\S]*?)\s*```/);
      const jsonString = jsonMatch ? jsonMatch[1] : content;
      lyricsData = JSON.parse(jsonString);
    } catch (parseError) {
      // If parsing fails, create a simple structure from the text
      const lines = content.split('\n').filter(line => line.trim().length > 0).slice(0, 4);
      lyricsData = {
        variations: [{
          lines: lines,
          syllableCounts: phraseLengths
        }]
      };
    }

    return NextResponse.json(lyricsData);
  } catch (error: any) {
    console.error('Lyrics generation error:', error);
    return NextResponse.json(
      { error: error.message || 'Failed to generate lyrics' },
      { status: 500 }
    );
  }
}

