# Encryption Interface

本文件定義 detectviz 中的加密模組 interface 設計，提供統一的 Encryptor 抽象與可替換實作方式。支援密碼、API 金鑰、設定參數等敏感資訊加解密功能。

---

## 設計目標

- 支援多種加密實作（如 AES、Vault、KMS）
- 封裝為抽象 interface，可由 plugin 註冊與注入
- 與 service 或 store 層整合，保護儲存與傳輸中的敏感資料

---

## Interface 定義

```go
type Encryptor interface {
    ID() string
    Encrypt(plain []byte) ([]byte, error)
    Decrypt(cipher []byte) ([]byte, error)
}
```

---

## 實作範例

### AES Encryptor

```go
type aesEncryptor struct {
    key []byte
}

func (e *aesEncryptor) Encrypt(plain []byte) ([]byte, error) {
    // 使用 AES-CFB 加密邏輯
}

func (e *aesEncryptor) Decrypt(cipher []byte) ([]byte, error) {
    // 解密並驗證長度
}
```

---

## Plugin 註冊建議

```go
func init() {
    plugins.RegisterEncryptor("aes", func() encryption.Encryptor {
        return NewAESEncryptor(defaultKey)
    })
}
```

---

## 使用情境

- 儲存密碼、token、webhook URL 等設定前先加密
- 讀取時由 Injected Encryptor 解密後交給 service 層處理
- 不直接與 config / database 耦合，僅負責 encode/decode

---

## 安全性建議

- Encryptor 實作應具備 IV 與 padding 處理
- Key 應由 config 模組注入（非硬編碼）
- 不建議回傳 base64 字串，保持純 bytes 輸出
- 加密錯誤應封裝為自定錯誤型別避免資訊洩漏

---

## 延伸方向

- 支援 encryptor metadata，例如 version、provider name
- 支援 rotate key 機制
- CLI 工具支援 Encrypt/Decrypt 測試與驗證

---