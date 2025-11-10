import { NextRequest, NextResponse } from 'next/server';
import OpenAI from 'openai';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { type, currentLine, selectedText, context, musicalKey } = body;

    // Get API key and base URL from environment
    const apiKey = process.env.OPENAI_API_KEY;
    const baseURL = process.env.OPENAI_API_BASE_URL || 'https://api.openai.com/v1';
    const model = process.env.OPENAI_MODEL || 'gpt-4o-mini';

    if (!apiKey) {
      return NextResponse.json(
        { 
          error: 'API key not configured',
          suggestions: []
        },
        { status: 400 } // Return 400 Bad Request for missing configuration
      );
    }

    // Initialize OpenAI client with custom base URL for OpenAI-compatible APIs
    const openai = new OpenAI({
      apiKey: apiKey,
      baseURL: baseURL,
    });

    let prompt = '';
    let systemMessage = '';

    if (type === 'completion') {
      // Suggest line completion
      systemMessage = 'You are a helpful songwriting assistant. Suggest natural, creative completions for song lyrics that match the style and context.';
      prompt = `The user is writing song lyrics. They've started a line and need help completing it.

Current line: "${currentLine}"

Context (previous lines):
${context ? context.split('\n').slice(-3).join('\n') : '(No previous context)'}

Musical key: ${musicalKey || 'C Major'}

Suggest 3-5 natural completions for this line. Each suggestion should:
1. Complete the thought naturally
2. Match the rhythm and flow of the existing line
3. Fit the context and style of the song
4. Be creative and engaging

Return ONLY a JSON array of completion strings, like:
["completion 1", "completion 2", "completion 3"]`;

    } else if (type === 'improvement') {
      // Suggest improvements
      systemMessage = 'You are a helpful songwriting assistant. Suggest improvements to make lyrics more impactful, poetic, or better fitting the musical context.';
      prompt = `The user has selected a line or phrase from their lyrics and wants suggestions to improve it.

Selected text: "${selectedText}"

Context (full lyrics):
${context || '(No context)'}

Musical key: ${musicalKey || 'C Major'}

Suggest 3-5 improved versions of this line/phrase. Each suggestion should:
1. Maintain the original meaning and intent
2. Be more poetic, impactful, or memorable
3. Better fit the musical rhythm and flow
4. Match the style of the rest of the lyrics

Return ONLY a JSON array of improvement strings, like:
["improved version 1", "improved version 2", "improved version 3"]`;

    } else if (type === 'alternative') {
      // Suggest alternatives
      systemMessage = 'You are a helpful songwriting assistant. Suggest alternative phrasings that maintain the meaning but offer different word choices, rhythms, or styles.';
      prompt = `The user has selected a line or phrase from their lyrics and wants alternative phrasings.

Selected text: "${selectedText}"

Context (full lyrics):
${context || '(No context)'}

Musical key: ${musicalKey || 'C Major'}

Suggest 3-5 alternative phrasings that:
1. Maintain similar meaning and intent
2. Offer different word choices or rhythms
3. Provide variety in style or tone
4. Still fit the musical context

Return ONLY a JSON array of alternative strings, like:
["alternative 1", "alternative 2", "alternative 3"]`;

    } else {
      return NextResponse.json(
        { error: 'Invalid suggestion type' },
        { status: 400 }
      );
    }

    const completion = await openai.chat.completions.create({
      model: model,
      messages: [
        {
          role: 'system',
          content: systemMessage
        },
        {
          role: 'user',
          content: prompt
        }
      ],
      temperature: 0.8,
      max_tokens: 500,
    });

    const content = completion.choices[0]?.message?.content || '';
    
    // Try to parse JSON from response
    let suggestions: string[] = [];
    try {
      // Extract JSON from markdown code blocks if present
      const jsonMatch = content.match(/```json\s*([\s\S]*?)\s*```/) || content.match(/```\s*([\s\S]*?)\s*```/);
      const jsonString = jsonMatch ? jsonMatch[1] : content.trim();
      
      // Try to parse as array
      const parsed = JSON.parse(jsonString);
      if (Array.isArray(parsed)) {
        suggestions = parsed;
      } else if (parsed.suggestions && Array.isArray(parsed.suggestions)) {
        suggestions = parsed.suggestions;
      } else {
        // Fallback: split by lines
        suggestions = content.split('\n')
          .filter(line => line.trim().length > 0)
          .map(line => line.replace(/^[-•*]\s*/, '').replace(/^["']|["']$/g, '').trim())
          .filter(line => line.length > 0)
          .slice(0, 5);
      }
    } catch (parseError) {
      // If parsing fails, try to extract suggestions from text
      const lines = content.split('\n')
        .filter(line => line.trim().length > 0)
        .map(line => line.replace(/^[-•*]\s*/, '').replace(/^["']|["']$/g, '').trim())
        .filter(line => line.length > 0 && !line.match(/^(suggestion|option|alternative)/i))
        .slice(0, 5);
      
      suggestions = lines.length > 0 ? lines : [content.trim()];
    }

    // Ensure we have at least some suggestions
    if (suggestions.length === 0) {
      suggestions = ['No suggestions available'];
    }

    return NextResponse.json({
      suggestions: suggestions.map(text => ({
        text: text.trim(),
        type: type
      }))
    });
  } catch (error: any) {
    console.error('Lyrics suggestion error:', error);
    return NextResponse.json(
      { error: error.message || 'Failed to get suggestions' },
      { status: 500 }
    );
  }
}
