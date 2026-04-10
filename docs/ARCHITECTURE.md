# Architecture Overview

This project implements a secure file storage system with fine-grained access control using **Kratos** for gRPC services, **SpiceDB** (Authzed) for permission management, and **PostgreSQL** for metadata persistence.

## System Architecture

The system follows a layered architecture:
- **API**: gRPC and HTTP definitions (Protobuf).
- **Service**: Implements handlers for File, Share, Plan, Subscription, Dashboard, and Greeter.
- **Biz (Business Logic)**: Core domain logic and permission enforcement using `AuthRepo`.
- **Data (Persistence)**: Implements `Repo` interfaces for SQL (PostgreSQL/TimescaleDB) and Auth (SpiceDB).

> [!NOTE]
> Elasticsearch-based search is currently disabled in favor of database-native queries for the Dashboard.

## File Sharing Flows

The sharing system allows users to create shareable links with specific permission levels (`viewer`, `editor`).

### 1. Create Share Link
This flow generates a secure token and records the intent in both the database and SpiceDB.

```mermaid
sequenceDiagram
    participant Client
    participant Service as ShareService
    participant Biz as ShareUseCase
    participant SpiceDB
    participant DB as PostgreSQL

    Client->>Service: CreateShareLink(ResourceID, Level)
    Service->>Biz: CreateShareLink(input)
    Biz->>SpiceDB: CheckPermission(Actor, "share", Resource)
    alt Denied
        SpiceDB-->>Biz: Unauthorized
        Biz-->>Client: Error: Permission Denied
    else Allowed
        SpiceDB-->>Biz: Authorized
        Biz->>Biz: Generate Secure Token (NanoID)
        Biz->>DB: Begin Transaction
        DB->>DB: Insert Share Link Record
        Biz->>SpiceDB: WriteRelationship(Resource, "share", Actor)
        Note over Biz, SpiceDB: Grants the creator permanent sharing rights on the file
        DB->>DB: Commit Transaction
        Biz-->>Client: Return Token
    end
```

### 2. Revoke Share Link
Invalidates a share link and removes the associated permissions.

```mermaid
sequenceDiagram
    participant Client
    participant Service as ShareService
    participant Biz as ShareUseCase
    participant SpiceDB
    participant DB as PostgreSQL

    Client->>Service: RevokeShareLink(Token)
    Service->>Biz: RevokeShareLink(input)
    Biz->>DB: Get Share Link Info
    Biz->>SpiceDB: CheckPermission(Actor, "share", Resource)
    alt Denied
        SpiceDB-->>Biz: Unauthorized
        Biz-->>Client: Error: Permission Denied
    else Allowed
        SpiceDB-->>Biz: Authorized
        Biz->>DB: Begin Transaction
        DB->>DB: Delete Share Link Record
        Biz->>SpiceDB: DeleteRelationship(Resource, "share", Actor)
        DB-->>DB: Commit Transaction
        Biz-->>Client: Success
    end
```

### 3. Update Share Link Permissions
Changes the access level of an existing share link.

```mermaid
sequenceDiagram
    participant Client
    participant Service as ShareService
    participant Biz as ShareUseCase
    participant SpiceDB
    participant DB as PostgreSQL

    Client->>Service: UpdateShareLink(Token, NewLevel)
    Service->>Biz: UpdateUserPermission(input)
    Biz->>DB: Get Share Link Info
    Biz->>SpiceDB: CheckPermission(Actor, "share", Resource)
    
    alt Allowed
        Biz->>DB: Update Permission Level
        Biz->>SpiceDB: SwapRelationship(Resource, Actor, OldRelation, NewRelation)
        Note right of SpiceDB: Atomically updates the relation in SpiceDB
        Biz-->>Client: Success
    end
```

### 4. Resolve Share Link (Visitor Flow)
Allows an anonymous visitor to exchange a token for information about the shared resource.

```mermaid
sequenceDiagram
    participant Visitor
    participant Service as ShareService
    participant Biz as ShareUseCase
    participant DB as PostgreSQL

    Visitor->>Service: ResolveShareLink(Token)
    Service->>Biz: ResolveShareLink(input)
    Biz->>DB: Get Share Link Info
    alt Not Found / Expired
        DB-->>Biz: nil / expired
        Biz-->>Visitor: Error: Link Invalid
    else Valid
        DB-->>Biz: Share Link Data (ResourceID, Level)
        Biz-->>Visitor: Return Metadata
        Note over Visitor, Service: UI then uses ResourceID to fetch content
    end
```

## Permission Model (SpiceDB)

Permissions are modeled based on the following resources:
- `user`: System users.
- `service`: Machine-to-machine actors.
- `file` / `folder`: Resources being protected.

Common relations:
- `owner`: Full control (Read, Write, Delete, Share).
- `writer`: Can edit and share.
- `viewer`: Read-only access.
- `share`: Explicit permission to manage share links for a resource.
