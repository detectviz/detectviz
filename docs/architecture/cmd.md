

# CLI Architecture

本文件說明 detectviz 專案中 CLI（命令列工具）的設計原則、結構規劃與擴充建議。  
本專案採用 [Cobra](https://github.com/spf13/cobra) 作為 CLI 框架，支援多層次指令、模組化註冊與外掛注入。

---

## 設計目標

- 對齊模組架構（rule, plugin, store 等皆可掛 CLI 指令）
- 提供工具型操作介面（如匯入/匯出/驗證/初始化）
- 支援多層指令與自動補全
- 便於 plugin 動態註冊新指令

---

## 專案結構建議

```
pkg/cmd/
├── root.go         # detectviz 主命令，包含 global flags
├── rule/
│   ├── apply.go    # detectviz rule apply
│   └── list.go     # detectviz rule list
├── plugin/
│   ├── list.go     # detectviz plugin list
│   └── enable.go   # detectviz plugin enable
├── config/
│   └── show.go     # detectviz config show
```

---

## 起手範例

`root.go`

```go
var rootCmd = &cobra.Command{
    Use:   "detectviz",
    Short: "Detectviz CLI 工具",
    Long:  "提供資料規則匯入、plugin 操作、系統診斷等命令",
}

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}
```

`rule/apply.go`

```go
var ruleApplyCmd = &cobra.Command{
    Use:   "apply",
    Short: "套用資料規則",
    Run: func(cmd *cobra.Command, args []string) {
        // 呼叫 service 或檔案讀取
    },
}

func init() {
    ruleCmd.AddCommand(ruleApplyCmd)
}
```

---

## plugin 註冊介面建議

```go
type CLIPlugin interface {
    Command() *cobra.Command
}
```

plugin 可透過以下註冊：

```go
manager.RegisterCLI(func() *cobra.Command {
    return &cobra.Command{
        Use: "custom",
        Run: func(cmd *cobra.Command, args []string) { ... },
    }
})
```

---

## 設定整合建議

搭配 `viper` 可支援：

- `--config` 指定組態檔
- `--verbose` 控制 log 層級
- 支援 env/config/flag 優先序

---

## 測試建議

- 每個 command 應具備獨立測試（使用 `pkg/cmd.Execute()` 模擬輸入）
- 建議使用 `os.Args` mock 套件或直接注入 `args []string`

---

## 延伸建議

- 支援 `detectviz completion bash/zsh/fish` 等補全指令
- 將 CLI 模組對應到 service 呼叫流程（避免重複實作邏輯）
- 可透過 lifecycle plugin 注入 CLI 指令模組，對齊 plugin 架構

---

> 💡 detectviz 採用 library-first 設計，將 CLI 模組放置於 `pkg/cmd/`，方便主程式與 plugin 共用指令建構器，並維持一致性。