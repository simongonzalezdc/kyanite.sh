# Knowledge Base Stub Implementation Guide

## Overview

The Knowledge Base Stub is a lightweight implementation of the Knowledge Base interface that provides basic songwriting knowledge and graceful degradation when a full RAG (Retrieval-Augmented Generation) system is not available. This implementation is part of Enhancement #6.5 Week 2 tasks.

## Features

### Core Functionality
- **Basic Songwriting Knowledge**: Pre-populated with essential songwriting information including rhyme schemes, chord progressions, lyrical techniques, song structure, and inspiration themes
- **Graceful Degradation**: Continues to function when the full knowledge base is unavailable
- **Search Capabilities**: Simple keyword matching across cards with category and tag filtering
- **Enhancement Provider**: Provides lyric and pattern suggestions based on knowledge base content

### Knowledge Categories
1. **Rhyme Schemes**: AABB, ABAB, and other common patterns
2. **Chord Progressions**: 50s progression, four-chord progression, etc.
3. **Lyrical Techniques**: Imagery, metaphor, show-don't-tell
4. **Song Structure**: Verse-chorus, bridge, etc.
5. **Inspiration**: Theme-based creative prompts

## Architecture

### Interface Definition

The Knowledge Base is defined by the `KnowledgeBase` interface in `internal/app/knowledge/interface.go`:

```go
type KnowledgeBase interface {
    Search(ctx context.Context, query string, options SearchOptions) (*SearchResult, error)
    AddCard(ctx context.Context, card Card) error
    GetCard(ctx context.Context, id string) (*Card, error)
    UpdateCard(ctx context.Context, card Card) error
    DeleteCard(ctx context.Context, id string) error
    GetCategories(ctx context.Context) ([]string, error)
    GetTags(ctx context.Context) ([]string, error)
    GetStatus(ctx context.Context) *KnowledgeStatus
    IsAvailable(ctx context.Context) bool
    Initialize(ctx context.Context) error
    Close(ctx context.Context) error
}
```

### Enhancement Provider

The `EnhancementProvider` interface adds AI enhancement capabilities:

```go
type EnhancementProvider interface {
    EnhanceLyrics(ctx context.Context, lyrics string, options SearchOptions) (*LyricSuggestion, error)
    EnhancePatterns(ctx context.Context, pattern string, options SearchOptions) (*PatternSuggestion, error)
    GetInspirationCards(ctx context.Context, theme string, options SearchOptions) (*SearchResult, error)
    GetKnowledgeBase() KnowledgeBase
}
```

## Integration with QuickIdeaAgent

The Knowledge Base Stub integrates with the QuickIdeaAgent to enhance AI suggestions:

1. **Context Enhancement**: The agent searches the knowledge base for relevant cards before generating suggestions
2. **Prompt Enrichment**: Knowledge base content is added to the AI prompts for better context
3. **Fallback Enhancement**: When AI is unavailable, knowledge base insights are added to fallback responses
4. **Status Indication**: The agent reports knowledge base availability and status

### Usage Example

```go
// Create agent with knowledge base
agent := ai.NewQuickIdeaAgent()

// Check knowledge base status
status := agent.GetKnowledgeBaseStatus(ctx)
if status.Available {
    fmt.Printf("KB available with %d cards\n", status.CardCount)
}

// Generate enhanced suggestions
req := ai.QuickRequest{
    Mode:    ai.QuickIdeaModeUnstick,
    Context: "I'm walking down the street tonight",
    Options: map[string]string{},
}

resp, err := agent.Generate(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Suggestions:", resp.Suggestions)
```

## UI Integration

### Status Bar Indicators

The status bar shows knowledge base availability with visual indicators:

- **Green/Bold**: Knowledge base is available (even if it's the stub)
- **Gray/Italic**: Knowledge base is unavailable
- **Status Text**: Shows "KB: Stub (N cards)" for stub implementation

### Example Status Display

```
Untitled | Ln 12, Col 24 | lyrics | 42 words | 280 chars | 12 lines | L W I B | KB: Stub (9 cards) | Saved 14:32:05
```

## Graceful Degradation

The stub implementation provides several layers of graceful degradation:

1. **Always Available**: The stub returns `true` for `IsAvailable()` to ensure the system continues to function
2. **Clear Status**: The status clearly indicates this is a stub implementation
3. **Fallback Content**: When no knowledge matches are found, generic suggestions are provided
4. **Error Handling**: All operations handle errors gracefully without breaking the UI

## Performance Characteristics

- **Search Time**: ~10ms simulated processing time
- **Memory Usage**: Minimal, stores cards in memory map
- **Card Count**: 9 pre-populated cards with essential songwriting knowledge
- **Response Time**: Consistent regardless of query complexity

## Testing

### Unit Tests

The stub implementation includes comprehensive unit tests in `internal/app/knowledge/stub_test.go`:

- Search functionality with various queries and filters
- Card management (add, get, update, delete)
- Category and tag retrieval
- Status reporting
- Enhancement provider functionality

### Integration Tests

Integration tests in `internal/app/ai/knowledge_integration_test.go` verify:

- QuickIdeaAgent integration with knowledge base
- Graceful degradation when knowledge base is unavailable
- Timeout handling
- Content type detection
- Knowledge-based enhancement of suggestions

### Running Tests

```bash
# Run knowledge base tests
go test ./internal/app/knowledge/...

# Run AI integration tests
go test ./internal/app/ai/...

# Run all tests
go test ./...
```

## Future Enhancement

The stub implementation is designed to be easily replaced by a full RAG implementation:

1. **Interface Compliance**: Full implementation of the KnowledgeBase interface
2. **Drop-in Replacement**: Same integration points with QuickIdeaAgent
3. **Status Indication**: Clear differentiation between stub and full implementation
4. **Migration Path**: All existing code will work with the full implementation

### Implementation Steps for Full KB

1. Implement ChromaDB integration for vector storage
2. Add embedding generation for semantic search
3. Implement proper RAG pipeline with LLM
4. Add real-time synchronization
5. Implement caching and performance optimization
6. Update status to reflect real KB availability

## Configuration

### Default Values

```go
// Search options
DefaultLimit:        5
DefaultMinRelevance: 0.6
DefaultUseCache:     true

// Knowledge base status
Version:     "stub-1.0.0"
Available:   false // Indicates this is a stub
CardCount:   9      // Number of pre-populated cards
```

### Customization

The stub can be customized by adding new cards:

```go
kb := knowledge.NewStubKnowledgeBase()

newCard := knowledge.Card{
    ID:       "custom-card",
    Title:    "Custom Knowledge",
    Content:  "Custom content for specific use case",
    Category: "custom",
    Tags:     []string{"custom", "specific"},
    Relevance: 0.8,
}

err := kb.AddCard(ctx, newCard)
```

## Troubleshooting

### Common Issues

1. **Knowledge Base Shows Unavailable**
   - This is expected for the stub implementation
   - Check that the knowledge base is properly initialized
   - Verify the status bar shows "KB: Stub" not "KB: Error"

2. **No Knowledge Enhancement in Suggestions**
   - Ensure the content type is detected correctly
   - Check that relevant cards exist for the query
   - Verify the search options aren't too restrictive

3. **Tests Fail with Import Errors**
   - Ensure all dependencies are properly imported
   - Check that the Go module is up to date
   - Verify the test files are in the correct directories

### Debug Information

Enable debug logging to troubleshoot issues:

```go
// Get detailed status
status := agent.GetKnowledgeBaseStatus(ctx)
fmt.Printf("KB Status: %+v\n", status)

// Check search results
result, err := kb.Search(ctx, query, options)
if err != nil {
    fmt.Printf("Search error: %v\n", err)
} else {
    fmt.Printf("Found %d cards in %v\n", len(result.Cards), result.Duration)
}
```

## Conclusion

The Knowledge Base Stub provides a solid foundation for songwriting assistance while maintaining graceful degradation. It enhances the AI suggestions with domain-specific knowledge and provides clear visual indicators about its availability. The implementation is designed to be easily replaced by a full RAG system when ready, requiring no changes to the integration points.