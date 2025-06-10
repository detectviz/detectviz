# CLI Interface (`pkg/cmd/`)

本文件定義 detectviz 專案中 CLI 模組的核心介面與模組設計方式，採用 [Cobra](https://github.com/spf13/cobra) 建構命令列工具，提供模組化、可注入的指令結構。

---

## 設計原則

- 採用 **library-first** 架構（指令定義與執行分離）
- 所有 CLI command 模組皆放置於 `pkg/cmd/` 下，可獨立測試與註冊
- 所有應用程式入口（如 `apps/cli/`）僅需呼叫 `cmd.Execute()` 作為進入點

---

## 類別與範型說明

### `root.go`

```go
var rootCmd = &cobra.Command{
    Use: "detectviz",
    Short: "Detectviz CLI 工具",
}

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}
```

- 提供主指令與 global flags
- 所有子指令將由其他模組動態註冊進來

---

### 子指令模組範例（如 rule/apply.go）

```go
var ruleApplyCmd = &cobra.Command{
    Use: "apply",
    Short: "套用資料規則",
    RunE: func(cmd *cobra.Command, args []string) error {
        return rule.ApplyFromFile(args[0])
    },
}

func init() {
    ruleCmd.AddCommand(ruleApplyCmd)
}
```

---

## 擴充接口建議

若需支援 plugin CLI 註冊，建議定義以下介面：

```go
type CLIPlugin interface {
    Command() *cobra.Command
}
```

註冊方式：

```go
plugins.RegisterCLI(func() *cobra.Command {
    return &cobra.Command{
        Use: "plugin:hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello from plugin")
        },
    }
})
```

---

## 結構建議

```
pkg/cmd/
├── root.go             # 主命令定義與 Execute()
├── rule/
│   ├── apply.go        # detectviz rule apply
│   └── list.go         # detectviz rule list
├── plugin/
│   ├── list.go
│   └── enable.go
├── config/
│   └── show.go
```

---

## 測試建議

- 所有 command 應可透過注入 `args []string` 呼叫測試
- 可使用 `cmd.SetArgs()` 與 `cmd.Execute()` 模擬 CLI 輸入
- plugin 註冊 CLI 可透過 lifecycle 註冊與測試

---

## 應用程式進入點

CLI 工具主程式建議放於：

```
apps/cli/main.go
```

執行方式：

```go
func main() {
    cmd.Execute()
}
```

---