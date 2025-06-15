# DetectViz Docs OverviewAdd commentMore actions

此文件目錄為 DetectViz 專案的文件與架構說明集合，供開發者與自動化 scaffold 工具（如 Cursor）參考。

---

## 文件分類與用途

1. [develop-guide.md](develop-guide.md) – 目錄與分類原則、命名慣例
2. [coding-style-guide.md](coding-style-guide.md) – 命名、interface、實作撰寫風格
3. [test-guide.md](test-guide.md) – 對應測試目錄與資料夾分層
4. [foundation.md](foundation.md) – 平台與模組可組合性核心觀念補充
5. `tmp/deprecated-*`：為舊 detectviz 原始碼、結構與 interface 歷史文件。
	- [deprecated-architecture](tmp/deprecated-architecture) 
	- [deprecated-sourcecode](tmp/deprecated-sourcecode)
	- [deprecated-interfaces](tmp/deprecated-interfaces)
6. `tmp/reference-sourcecode/`：外部系統實作原始碼，供映射與比較使用，僅供對照映射，不作 scaffold 依據。
	- [grafana](../tmp/reference-sourcecode/grafana) 
		- [grafana-plugin-sdk-go](../tmp/reference-sourcecode/grafana-plugin-sdk-go)
	- [oncall](../tmp/reference-sourcecode/oncall)
	- [telegraf](../tmp/reference-sourcecode/telegraf)