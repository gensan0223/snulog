```mermaid
graph TD
  CLI -->|gRPC| Server
  Server --> Usecase
  Usecase --> Repository
  Repository -->|SQL| PostgreSQL
```