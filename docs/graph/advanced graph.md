graph TD
    A[Workflow Start] --> B[Parallel Branch 1]
    A --> C[Parallel Branch 2]
    B --> E[API Gateway 2]
    E --> D[API Gateway 3]

    C --> F[API Gateway 3]
    C --> G[API Gateway 4]
    D --> H[Sync Point]

    F --> H
    G --> H
    H --> I[Final Processing]
    I --> J[Workflow Complete]
    
    style A fill:#51cf66
    style H fill:#91e5a3
    style J fill:#51cf66
