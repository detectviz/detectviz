# DetectViz Docs OverviewAdd commentMore actions

此文件目錄為 DetectViz 專案的文件與架構說明集合，供開發者與自動化 scaffold 工具（如 Cursor）參考。


---


## 目錄



```bash
.
├── deprecated-architecture/
├── deprecated-detectviz-sourcecode/
├── deprecated-interfaces/
├── reference-sourcecode/
│   ├── grafana
│   ├── grafana-plugin-sdk-go
│   ├── oncall
│   └── telegraf
├── coding-style-guide.md
├── develop-guide.md
├── foundation.md
├── test-guide.md
└── README.md
```

## 文件分類與用途

| 類別 | 路徑 | 說明 |
|------|------|------|
| 開發與程式規範 | `develop-guide.md`, `coding-style-guide.md` | 說明目錄架構、命名原則、程式風格 |
| 架構原則 | `foundation.md`, `test-guide.md` | 描述模組組成邏輯、測試分層、平台核心觀念 |
| 架構規劃文件 | `detectviz-source-code.md` | DetectViz scaffold 建立與模組化規則藍圖 |
| 舊版歷史檔案 | `deprecated-*` | 舊 detectviz 結構與 interface 歷史文件 |
| 參考實作 | `reference-sourcecode/grafana`, `.../telegraf` | 外部系統實作結構簡化版，供映射與比較使用 |

---

## Scaffold 工具（如 Cursor）參考優先順序

當進行 scaffold 或 plugin 自動產生時，建議依照下列優先讀取順序：

1. `develop-guide.md` – 目錄與分類原則、命名慣例
2. `coding-style-guide.md` – 命名、interface、實作撰寫風格
3. `test-guide.md` – 對應測試目錄與資料夾分層
4. `foundation.md` – 平台與模組可組合性核心觀念補充

---

## 注意

- 所有 `deprecated-*` 目錄為歷史遷移保留，Cursor 不應參考。
- 所有 `reference-sourcecode/` 僅供對照映射，不作 scaffold 依據。