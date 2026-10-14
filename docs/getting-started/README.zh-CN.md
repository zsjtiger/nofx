# 快速开始

NOFX 运行在您自己的服务器上。在网页端添加 AI 模型和交易账户，然后创建策略和交易员即可。

## 1. 安装

参见 [README](../../README.md#install)。默认访问 `http://127.0.0.1:3000`。

## 2. 注册

第一个创建的账户将成为实例的所有者。

## 3. 添加 AI 模型

任何兼容 OpenAI Chat Completions 接口的模型均可使用。在 AI 模型设置中填写 Base URL、模型名称和 API Key。

## 4. 添加交易账户

| 账户类型 | 市场 | 需要准备 |
| :--- | :--- | :--- |
| `alpaca_paper` | 美股 | Alpaca 纸盘的 Key 和 Secret |
| `cn_paper` | A 股 | Tushare Pro 的 Token |

两个市场目前都是纸上交易。各市场适用的规则见[迁移方案](../plans/2026-10-14-cn-us-equity-migration.md)。

## 5. 创建策略和交易员

在 **Strategy Studio** 中创建策略，设置标的和风控参数。然后创建交易员，选择策略、AI 模型和账户，并启动。

交易员只接受与其市场匹配的标的。例如 `AAPL` 适用于 `us_equity`，`600519.SH` 适用于 `cn_equity`。
