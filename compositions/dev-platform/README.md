# DetectViz 開發預設平台 (dev-platform)

> DetectViz 預設 scaffold 組合平台，適用於本地開發、plugin 驗證與平台整合測試。

## 概述

dev-platform 是 DetectViz 的預設開發組合，提供一個最小但完整的平台環境。這個組合包含了：

- **JWT 認證系統**：基於 core-auth-jwt 插件的安全認證
- **OtelZap 結構化日誌**：使用 OpenTelemetry 整合的日誌系統
- **記憶體式基礎設施**：快取和事件匯流排的記憶體實作
- **自動插件管理**：支援插件註冊與健康檢查
- **開發友好配置**：適合本地開發的預設值

## 檔案結構

```
dev-platform/
├── composition.yaml    # 主要組合配置文件
├── meta.yaml          # 組合元資料和描述
└── README.md          # 使用說明文件
```

## 配置檔案說明

### composition.yaml

主要的組合配置文件，定義了：

- **Core Plugins**：核心插件配置
  - `core-auth-jwt`：JWT 認證插件
  - `otelzap-logger`：日誌插件
- **Applications**：應用程式配置（預設只啟用 server）
- **Infrastructure**：基礎設施配置（記憶體快取、事件匯流排）
- **Health Check**：健康檢查配置，支援 `plugins.include: all`

#### 重要配置項目

```yaml
health:
  enabled: true
  interval: "30s"
  timeout: "5s"
  plugins:
    include: all  # 包含所有插件的健康檢查
```

### meta.yaml

組合的元資料文件，包含：

- 組合基本資訊（名稱、版本、描述）
- 功能特色列表
- 系統需求
- 使用場景
- 維護者資訊

## 使用方式

### 本地開發

1. 確保滿足系統需求：
   - Go 1.21+
   - 至少 256MB 記憶體

2. 啟動 DetectViz：
   ```bash
   go run apps/server/main.go --composition=dev-platform
   ```

3. 預設配置：
   - 伺服器埠號：8080
   - JWT Secret：`dev-secret`
   - 日誌等級：info
   - 日誌格式：json

### 插件開發

dev-platform 提供了完整的插件開發環境：

1. **插件註冊**：支援自動插件發現和註冊
2. **健康檢查**：自動包含所有已註冊插件的健康檢查
3. **日誌整合**：使用 OtelZap 提供結構化日誌和追蹤

### 測試環境

適用於：
- 框架功能測試
- 插件整合測試
- 平台架構驗證

## 插件配置

### 核心插件

#### core-auth-jwt
```yaml
- name: core-auth-jwt
  type: auth
  enabled: true
  config:
    secret: dev-secret
    issuer: "detectviz-dev"
    expiry_time: "24h"
```

#### otelzap-logger
```yaml
- name: otelzap-logger
  type: logger
  enabled: true
  config:
    level: info
    format: json
    output: stdout
```

### 社群插件

預設禁用，可根據需要啟用：

```yaml
community_plugins:
  - name: prometheus-importer
    type: importer
    enabled: false  # 設為 true 來啟用
    config:
      endpoint: "http://localhost:9090"
```

## 健康檢查

dev-platform 支援全面的健康檢查：

- **系統健康檢查**：每 30 秒檢查一次
- **插件健康檢查**：包含所有已註冊插件
- **超時設定**：5 秒超時限制

健康檢查端點：
```
GET /health
```

## 故障排除

### 常見問題

1. **埠號衝突**
   - 修改 `composition.yaml` 中的 `applications.server.config.port`

2. **JWT 認證失敗**
   - 檢查 `core-auth-jwt.config.secret` 設定
   - 確認 token 格式正確

3. **日誌輸出問題**
   - 檢查 `otelzap-logger.config.level` 設定
   - 確認 `output` 設定為 `stdout` 或有效檔案路徑

### 除錯模式

啟用除錯日誌：
```yaml
logging:
  level: debug
```

## 參考資料

- [DetectViz 開發指南](../../docs/develop-guide.md)
- [插件開發文件](../../docs/interfaces/)
- [組合平台架構](../../docs/composition.md) 