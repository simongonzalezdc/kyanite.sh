# Collaboration Feature (Future)

## Overview

The collaboration feature enables real-time multi-user songwriting sessions. This is a **planned future feature** that is currently disabled by default for single-user mode.

## Current Status

**Status:** Disabled (Future Feature)

The collaboration code is fully implemented but disabled via a feature flag. This allows the application to run efficiently in single-user mode while preserving the code for future expansion.

## Enabling Collaboration

To enable collaboration features, set the feature flag in your configuration:

### Via Config File (`~/.config/noise/noise.yaml`)

```yaml
features:
  enable_collaboration: true
```

### Via Environment Variable

```bash
export NOISE_FEATURES_ENABLE_COLLABORATION=true
```

## What Gets Enabled

When collaboration is enabled:

1. **Database Tables**: Additional tables are created for:
   - `collaboration_sessions` - Active collaboration sessions
   - `session_participants` - Users in each session
   - `document_operations` - Operational transform operations
   - `user_presence` - Real-time user status
   - `session_invitations` - Session invites
   - `conflict_resolutions` - Conflict resolution history

2. **Managers**: The following managers are initialized:
   - `CollaborationManager` - Main coordinator
   - `PresenceManager` - User presence tracking
   - `SessionManager` - Session lifecycle
   - `InvitationManager` - Invite handling
   - `ConflictResolver` - Conflict resolution strategies

3. **UI Components**:
   - Presence indicators
   - Collaboration status bar
   - Conflict resolution dialogs

## Running Collaboration Tests

Collaboration tests are excluded from the default test run. To run them:

```bash
# Run all collaboration tests
go test -tags collaboration ./internal/collaboration/...
go test -tags collaboration ./internal/ui/collaboration/...
go test -tags collaboration ./test/...
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    CollaborationManager                      │
│  ┌─────────────┐ ┌─────────────┐ ┌──────────────────────┐   │
│  │ Session     │ │ Presence    │ │ Conflict             │   │
│  │ Manager     │ │ Manager     │ │ Resolver             │   │
│  └─────────────┘ └─────────────┘ └──────────────────────┘   │
│  ┌─────────────┐                                             │
│  │ Invitation  │                                             │
│  │ Manager     │                                             │
│  └─────────────┘                                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Database (SQLite)                         │
│  collaboration_sessions, session_participants,               │
│  document_operations, user_presence, session_invitations,    │
│  conflict_resolutions                                        │
└─────────────────────────────────────────────────────────────┘
```

## Key Components

### Session Management
- Create and manage collaboration sessions
- Handle user joining and leaving
- Track session state and settings

### Presence System
- Real-time user status (online, away, busy, offline)
- Cursor position tracking
- Heartbeat mechanism for connection health

### Conflict Resolution
- Operational Transform (OT) based synchronization
- Multiple merge strategies available
- Manual and automatic resolution options

### Invitation System
- Invite users to sessions
- Role-based permissions (owner, editor, viewer)
- Expiring invitations

## Future Development

When ready to enable collaboration:

1. Set the feature flag to `true`
2. Run the application - collaboration tables will be created automatically
3. Run collaboration tests to verify functionality
4. Consider adding network transport layer for remote collaboration

## Files

### Core Implementation
- `internal/collaboration/` - Core collaboration logic
  - `manager.go` - Main collaboration manager
  - `session_manager.go` - Session lifecycle
  - `presence.go` - User presence
  - `conflict_resolution.go` - Conflict handling
  - `constants.go` - Shared constants

### UI Components
- `internal/ui/collaboration/` - Collaboration UI
  - `presence_indicator.go` - Shows online users
  - `status_bar.go` - Collaboration status

### Database Schema
- `internal/infra/db/schema.go` - `CollaborationSchema` constant

### Configuration
- `internal/config/config.go` - `FeaturesConfig.EnableCollaboration`
