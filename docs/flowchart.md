flowchart 
    %% 說明：compositions/ 定義平台組合方案，apps/ 則為其導出的應用執行項目

    %% 應用層
    subgraph App["應用層 - apps/"]
        direction TB
        Server["Web API Server<br/>📁 apps/server/"]
        CLI["CLI Tools<br/>📁 apps/cli/"]
        Agent["Monitoring Agent<br/>📁 apps/agent/"]
        TestKit["Test Kit<br/>📁 apps/testkit/"]
    end

    %% 對外接口
    subgraph Ports["對外接口層 - internal/ports/"]
        direction TB
        HTTP["HTTP<br/>📁 internal/ports/http/"]
        GRPC["gRPC<br/>📁 internal/ports/grpc/"]
        CLI["CLI<br/>📁 internal/ports/cli/"]
        Web["Web<br/>📁 internal/ports/web/"]
    end

    %% 服務實作層
    subgraph Service["服務實作層 - internal/"]
        direction TB
        CoreSvc["核心服務<br/>📁 internal/services/<br/>alerting/ | monitoring/<br/>analytics/ | security/"]
        Adapter["適配器層<br/>📁 internal/adapters/<br/>importers/ | exporters/<br/>integrations/"]
        Repo["資料存取<br/>📁 internal/repositories/<br/>alert/ | rule/ | metric/<br/>user/ | organization/"]
    end

    %% 基礎設施層
    subgraph Infra["🏭 基礎設施層 - internal/"]
        direction LR
        Runtime["運行時<br/>📁 internal/platform/<br/>runtime/bootstrap.go<br/>shutdown.go | health.go"]
        Infrastructure["基礎組件<br/>📁 internal/infrastructure/<br/>eventbus/ | cache/<br/>logging/ | tracing/"]
    end

    %% 平台核心層
    subgraph Platform["🏗️ 平台核心層 - pkg/platform/"]
        direction TB
        Contracts["服務契約<br/>📁 pkg/platform/contracts/<br/>alerting/ | monitoring/<br/>notification/ | storage/"]
        Registry["註冊中心<br/>📁 pkg/platform/registry/<br/>interface.go | discovery.go<br/>composer.go"]
        Composition["組合引擎<br/>📁 pkg/platform/composition/<br/>app.go | module.go<br/>plugin.go"]
    end

    %% 業務領域層
    subgraph Domain["📋 業務領域層 - pkg/"]
        direction LR
        Models["領域模型<br/>📁 pkg/domain/<br/>alert/ | metric/ | rule/<br/>user/ | organization/"]
        Config["配置管理<br/>📁 pkg/config/<br/>types/ | schema/<br/>composition/ | loader/"]
        Shared["共享組件<br/>📁 pkg/shared/<br/>errors/ | utils/<br/>constants/ | types/"]
    end

    %% 插件生態
    subgraph Plugin["🔌 插件生態 - plugins/"]
        direction TB
        CorePlugin["核心插件<br/>📁 plugins/core/<br/>auth/ | middleware/"]
        CommunityPlugin["社群插件<br/>📁 plugins/community/<br/>importers/ | exporters/<br/>integrations/ | tools/"]
        CustomPlugin["自定義插件<br/>📁 plugins/custom/<br/>example/"]
    end

    %% 組合配置
    subgraph Combo["🎼 預設組合 - compositions/"]
        Presets["預設組合<br/>📁 monitoring-stack/<br/>📁 alerting-platform/<br/>📁 analytics-platform/<br/>📁 security-platform/"]
    end

    %% 主要數據流
    Combo --> App
    App --> Ports
    Ports --> Service
    Service --> Infra
    Platform --> Domain
    %% Platform 僅依賴 domain 的 interface 與模型，不應耦合具體邏輯
    %% Platform 提供 registry/composition 功能給 Service 使用，並不直接依賴 service 實作
    Combo --> Platform
    Plugin --> Platform
    %% plugins 僅實作 interface 並向上註冊進 platform，不應被 internal 依賴
 
    %% 樣式定義
    classDef app fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef platform fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef domain fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    classDef service fill:#e8f5e8,stroke:#388e3c,stroke-width:2px
    classDef infra fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef plugin fill:#fff8e1,stroke:#fbc02d,stroke-width:2px
    classDef combo fill:#f1f8e9,stroke:#689f38,stroke-width:2px

    class App app
    class Ports service
    class Service service
    class Infra infra
    class Platform platform
    class Domain domain
    class Plugin plugin
    class Combo combo
