# Integration Testing Guide for noise.sh

This guide provides comprehensive information about the integration testing framework for noise.sh, including best practices, troubleshooting, and extending the test suite.

## Overview

The integration testing framework validates the complete functionality of noise.sh by testing the interactions between all major components:

- **Database Layer**: Song storage, versioning, and retrieval
- **AI Integration**: Quick idea generation, brainstorming, and content analysis
- **Collaboration**: Session management, participant updates, and conflict resolution
- **UI Components**: Theme switching, dashboard rendering, and editor integration
- **Error Handling**: Error categorization, recovery mechanisms, and graceful degradation
- **Performance**: Concurrent operations, database performance, and memory usage

## Test Architecture

### Test Environment Setup

All tests use the `E2ETestSetup` struct which provides:
- Isolated temporary directories for each test run
- In-memory databases with full schema initialization
- Mock AI providers for consistent testing
- Clean setup and teardown mechanisms

```go
setup := NewE2ETestSetup(t)
defer setup.Cleanup()
```

### Component Integration

The test framework ensures realistic integration between components:
- Database transactions are properly committed and rolled back
- AI providers simulate real-world behavior with controlled responses
- Collaboration sessions manage concurrent access correctly
- Theme switches propagate through all UI components
- Error handling recovers from various failure scenarios

## Test Categories

### 1. Workflow Tests (`comprehensive_integration_test.go`)

These tests validate complete user workflows:

#### TestCompleteSongCreationWorkflow
- Creates a new song through the editor service
- Adds multiple versions with different content
- Tests version retrieval and integrity
- Validates export functionality

#### TestAIIntegrationWorkflow
- Tests AI knowledge base availability
- Validates quick idea generation
- Tests editor AI integration (brainstorming, continue mode)
- Validates content type detection

#### TestCollaborationWorkflow
- Creates collaboration sessions
- Adds participants with different roles
- Tests concurrent participant updates
- Validates session lifecycle management

#### TestThemeSystemIntegration
- Tests initial theme loading
- Validates theme switching across all available themes
- Tests theme persistence
- Validates fallback behavior for invalid themes

#### TestErrorHandlingAndRecovery
- Tests basic error handling
- Validates error recovery with retry mechanisms
- Tests error categorization and prioritization
- Validates concurrent error handling

#### TestDataPersistenceAndBackup
- Creates test data with multiple versions
- Tests auto-save functionality
- Validates version retrieval and integrity
- Tests backup creation and recovery

#### TestPerformanceUnderLoad
- Tests concurrent song creation and editing
- Validates performance targets
- Tests database performance under load
- Tests concurrent error handling

### 2. Load Testing Framework (`load_testing_framework.go`)

A specialized framework for performance testing:

#### Features
- Configurable concurrent users with ramp-up time
- Configurable operations per user
- Performance metrics collection
- Error tracking and categorization
- AI and collaboration integration testing

#### Configuration
```go
config := LoadTestConfig{
    ConcurrentUsers:      5,
    OperationsPerUser:    10,
    TestDuration:         30 * time.Second,
    RampUpTime:           5 * time.Second,
    TargetSuccessRate:    0.95,
    MaxOperationTime:     100 * time.Millisecond,
    EnableAIIntegration:  true,
    EnableCollaboration:  true,
}
```

#### Metrics Collected
- Total operations and success rate
- Operation time statistics (min, max, average)
- Operations per second
- AI requests and collaboration sessions
- Database operations
- Error counts by type

## Best Practices

### Writing Tests

1. **Use Descriptive Names**: Test names should clearly describe what is being tested
2. **Follow AAA Pattern**: Arrange, Act, Assert structure for clarity
3. **Test Edge Cases**: Include tests for error conditions and boundary cases
4. **Use Isolated Environments**: Each test should be independent and not rely on others
5. **Clean Up Resources**: Always clean up resources in defer statements

### Test Data Management

1. **Use Unique Identifiers**: Avoid constraint violations by using unique data
2. **Validate Data Integrity**: Check that data is correctly stored and retrieved
3. **Test Data Relationships**: Validate foreign key constraints and relationships
4. **Clean Up Test Data**: Ensure test data doesn't leak between tests

### Performance Testing

1. **Set Realistic Targets**: Configure performance targets based on expected usage
2. **Test Concurrent Operations**: Validate behavior under concurrent load
3. **Monitor Resource Usage**: Track memory and CPU usage during tests
4. **Test Error Scenarios**: Include error conditions in performance tests

## Troubleshooting

### Common Issues

#### Database Lock Errors
```
database is locked (5) (SQLITE_BUSY)
```

**Solution**: Reduce concurrent operations or increase ramp-up time in load tests.

#### Theme Test Failures
```
Expected theme X, got Y
```

**Solution**: Check theme registry and manager implementation. Verify theme names match exactly.

#### AI Integration Failures
```
AI generation failed
```

**Solution**: These are expected in test environments. Tests should continue with fallback behavior.

#### Version Content Mismatch
```
Version X content mismatch: expected Y, got Z
```

**Solution**: Check version ordering. Versions are returned in DESC order by created_at.

### Debugging Tips

1. **Enable Verbose Logging**: Use `-v` flag to see detailed test output
2. **Check Test Logs**: Review log messages for warnings and errors
3. **Isolate Failing Tests**: Run individual tests to identify specific issues
4. **Use Test Breakpoints**: Add temporary logging to debug complex issues

## Extending the Test Suite

### Adding New Workflow Tests

1. Create a new test function following the naming pattern `TestXxxWorkflow`
2. Use `NewE2ETestSetup` to create a test environment
3. Implement the test following the AAA pattern
4. Add appropriate assertions and logging
5. Include cleanup in a defer statement

### Adding Load Test Scenarios

1. Define a new configuration in the test
2. Create a new test function with the `TestLoadTestingXxx` naming pattern
3. Use the `LoadTester` framework with your configuration
4. Validate the results against your targets
5. Add appropriate logging for debugging

### Adding New Component Tests

1. Identify the component interactions to test
2. Create mock implementations if needed
3. Write tests that validate the integration points
4. Include both success and failure scenarios
5. Add performance tests if applicable

## Continuous Integration

### CI Configuration

The tests are designed to run in CI/CD environments:

1. **Minimal Resource Usage**: Tests use minimal memory and CPU
2. **Reasonable Time Limits**: Tests complete within acceptable timeframes
3. **Clear Output**: Tests provide clear output for debugging
4. **Graceful Failure Handling**: Tests handle expected failures gracefully

### CI Best Practices

1. **Run Tests in Parallel**: Configure CI to run tests in parallel when possible
2. **Cache Dependencies**: Cache Go modules to speed up builds
3. **Use Fixed Versions**: Pin dependency versions for consistent builds
4. **Generate Test Reports**: Generate test reports for tracking trends

## Conclusion

The integration testing framework for noise.sh provides comprehensive coverage of all major components and their interactions. By following this guide, you can effectively extend and maintain the test suite to ensure the reliability and performance of the application.

For more information, refer to the test files in this directory and the main application documentation.