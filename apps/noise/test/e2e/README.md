# End-to-End Testing for noise.sh

This directory contains comprehensive end-to-end tests for the noise.sh application, covering all major components and their interactions.

## Test Structure

### Comprehensive Integration Tests (`comprehensive_integration_test.go`)

These tests validate the complete workflows of the application:

1. **TestCompleteSongCreationWorkflow** - Tests the entire song creation process from creation to export
2. **TestAIIntegrationWorkflow** - Validates AI integration throughout the application
3. **TestCollaborationWorkflow** - Tests real-time collaboration features
4. **TestThemeSystemIntegration** - Validates theme switching and UI consistency
5. **TestErrorHandlingAndRecovery** - Tests comprehensive error handling and recovery mechanisms
6. **TestDataPersistenceAndBackup** - Tests data persistence and backup/recovery functionality
7. **TestPerformanceUnderLoad** - Tests system performance under stress conditions

### Load Testing Framework (`load_testing_framework.go` and `load_testing_framework_test.go`)

A specialized framework for load testing the application with configurable parameters:

- **Concurrent user simulation** with ramp-up time
- **Configurable operations per user**
- **Performance metrics collection**
- **Error tracking and categorization**
- **AI and collaboration integration testing**

## Running Tests

### Run All E2E Tests
```bash
cd test/e2e
go test -v
```

### Run Specific Test
```bash
cd test/e2e
go test -v -run TestCompleteSongCreationWorkflow
```

### Run Load Testing Framework
```bash
cd test/e2e
go test -v -run TestLoadTestingFramework
```

## Test Environment

All tests use isolated temporary environments with:
- Temporary databases for each test run
- Mock AI providers for consistent testing
- Isolated file systems
- Clean setup and teardown

## Performance Targets

The load testing framework uses these default targets (configurable):
- Success rate: 95%
- Maximum operation time: 100ms
- Concurrent users: 5
- Operations per user: 10

## Test Coverage

The e2e tests cover integration between:

- **Database Layer**: Song storage, versioning, and retrieval
- **AI Integration**: Quick idea generation, brainstorming, and content analysis
- **Collaboration**: Session management, participant updates, and conflict resolution
- **UI Components**: Theme switching, dashboard rendering, and editor integration
- **Error Handling**: Error categorization, recovery mechanisms, and graceful degradation
- **Performance**: Concurrent operations, database performance, and memory usage

## Troubleshooting

### Database Lock Errors
If you see "database is locked" errors during load testing, reduce the number of concurrent users or increase the ramp-up time.

### AI Integration Failures
AI integration tests use mock providers, so failures are expected in test environments. The tests are designed to continue with fallback behavior.

### Theme Test Failures
Theme tests validate theme switching behavior. If they fail, check the theme registry and manager implementation.

## Adding New Tests

When adding new e2e tests:

1. Use the `NewE2ETestSetup` function to create a test environment
2. Ensure proper cleanup in a `defer` statement
3. Follow the existing pattern for test structure and assertions
4. Add appropriate logging for debugging
5. Consider adding load testing scenarios if applicable

## Continuous Integration

These tests are designed to run in CI/CD environments:
- They use minimal resources
- They have reasonable time limits
- They provide clear output for debugging
- They handle expected failures gracefully