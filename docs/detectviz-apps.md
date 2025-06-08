# Detectviz Apps 應用層總覽

## 模組原型範例

- alert 模組 是流程編排者，觸發 analyzer、reporter。
- analyzer 模組 不實作模型，而是根據規則呼叫 anomaly-python 提供的多種 API。
- reporter 模組 可接收標記後的結果，呼叫 llm-python 生成摘要，再整合為報表。
- anomaly-python 是多種異常偵測模型的 REST API 服務。
- llm-python 是單一 LLM API（摘要、分析、重寫）服務。
