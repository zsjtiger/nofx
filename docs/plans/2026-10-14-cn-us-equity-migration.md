# 去除加密货币，支持 A 股与美股（实施计划）

## 已确认的决策

| 问题 | 决定 |
|---|---|
| 加密货币 | 彻底删除，不做隐藏 |
| A 股券商 | 大智慧（待确认：其个人量化交易接口是否可用，见"阻塞项"） |
| A 股数据源 | Tushare |
| 交易模式 | 先纸上交易，后实盘 |
| 大模型 | 用户自配 API Key，兼容 OpenAI 协议（自定义 base URL 与模型名）；仅删除 Claw402 支付链路（`mcp/payment`、x402、USDC 计费） |
| 货币与时区 | A 股：CNY / 北京时间；美股：USD / 美东时间 |

## 阻塞项

- **大智慧交易接口未确认**。公开资料里未找到大智慧面向个人的量化下单 API；大智慧 App 的交易功能是手动下单。个人量化交易常见的接口是 QMT 和 Ptrade，需要开户券商开通。
- 在确认接入路径前，A 股的"交易适配器"只能先做接口占位，不影响美股纸上交易的推进。
- 需要用户确认：是否有大智慧的开发者权限或券商侧的 QMT / Ptrade 权限。

## 阶段划分

### 阶段 1：市场抽象（基础）

- 在交易员/交易所配置中增加 `market` 字段：`us_equity` / `cn_equity`（删除 `crypto` 后不再保留该值）。
- 定义市场能力集：是否支持做空、杠杆、保证金模式、最小交易单位（A 股 100 股）、交易时段、交易日历、涨跌停、结算规则（A 股 T+1，美股 T+1 结算但可当日回转）。
- 标的代码规范化：A 股 `600519.SH` / `000001.SZ`，美股 `AAPL`。
- 风控中的杠杆 clamp、做空限制改为读取能力集。
- 涉及文件：`trader/types/interface.go`、`kernel/schema.go`、`store/exchange.go`、`store/visibility.go`。

### 阶段 2：数据层

- 美股：`provider/alpaca`（已有 K 线）为主，补充报价；或沿用 TwelveData。
- A 股：新增 `provider/tushare`，覆盖日 K、分钟线、涨跌停价、停牌状态、复权因子。复权方式需统一为前复权或后复权。
- 股票池：Strategy Studio 的"币种宇宙"改为指数成分（沪深 300、中证 500、标普 500、纳斯达克 100）。
- 涉及文件：`provider/`、`market/`、`store/strategy.go`（`CoinSource`、`VergexChain` 相关字段需清理）。

### 阶段 3：交易适配（纸上）

- 美股：实现 Alpaca Paper Trading 适配 `trader/types.Trader`。
- A 股：纸上交易用本地撮合模拟器（基于 Tushare 行情），实现同一 `Trader` 接口，先不接真实券商。
- 实盘适配（大智慧 / QMT / Ptrade）留到阻塞项解决后再做。
- 纸上撮合需考虑：A 股 T+1、100 股整手、涨跌停不可成交、停牌不可成交、集合竞价时段。

### 阶段 4：风控与交易规则

- A 股：T+1（当日买入不可当日卖出）、100 股整数倍、涨跌停（主板 ±10%，创业板/科创板 ±20%，ST ±5%）、停牌、交易时段（09:30–11:30、13:00–15:00）。
- 美股：PDT 规则（账户低于 25k USD 时的日内交易限制）、盘前盘后时段、做空需借券可用。
- 保留现有节流逻辑（最短持仓、冷却、每小时入场上限）。
- 对 A 股，`OpenShort` 与杠杆返回不支持。

### 阶段 5：决策引擎与提示词

- `kernel/engine_prompt.go`、`kernel/prompt_builder.go` 删除资金费率、持仓量、清算热力图、Vergex 方向等加密专属输入。
- 替换为股票相关输入：成交量、换手率、涨跌停状态、指数走势、行业。
- 提示词按市场分别维护（中文 / 英文）。
- 输出 schema 按市场裁剪（A 股无做空、无杠杆）。

### 阶段 6：删除加密代码

- 删除 `trader/` 下的 9 个加密交易所适配（binance、bybit、okx、bitget、hyperliquid、aster、lighter、gate、kucoin、indodax 等）。
- 删除 `provider/coinank`、`provider/hyperliquid`、`provider/vergex`、`provider/nofxos`（视 nofxos 是否仅服务于加密而定，需逐一核对）。
- 删除 `mcp/payment`（x402 / Claw402 支付）、`wallet/`、Hyperliquid 相关 onboarding 与燃料费流程。
- 删除 `api/handler_hyperliquid_wallet.go`、`api/handler_wallet.go`、`api/handler_ai_cost.go` 中的 USDC 相关逻辑（需逐一核对）。
- 模型配置保留自定义 Key 模式（`api/handler_ai_model.go`），移除 Claw402 选项；确认 OpenAI 兼容的自定义 base URL 与模型名路径保持可用（`mcp/providers.go` 中 `ProviderOpenAI`）。
- 注意：Claw402 不只是支付层。`provider/nofxos/claw402.go` 以及 `store/ai_charge.go`、`api/launch_preflight.go`、`api/handler_onboarding.go`、`manager/trader_manager.go`、`telegram/bot.go` 等多处引用 claw402，删除前要先区分"支付"与"数据/模型路由"两类引用，逐一改掉。
- 数据库：`store/exchange.go` 中硬编码的交易所 ID 列表需更新；新增迁移处理已有加密交易员记录（建议直接停用，不迁移）。
- 前端：交易所图标、钱包 / 燃料费页面、Hyperliquid 文案、竞赛排行榜中的加密字段。
- 文档：README（中英文）、`docs/architecture/`、`docs/getting-started/`、`install.sh`、`.env.example`。
- 测试：删除依赖加密 API 的集成测试（`provider/coinank`、`provider/hyperliquid` 等）。

### 阶段 7：验证

- Go 单元测试：`go test ./...`，按市场补充 T+1、涨跌停、PDT 的针对性测试。
- 前端：`npm run test`、`npm run lint`、`npm run build`。
- 纸上交易：美股和 A 股各跑一段，记录决策与成交。

## 建议的执行顺序

1. 阶段 1（市场抽象）和阶段 6 中的"删除加密交易所适配"可以并行准备，但先做阶段 1，避免删除后无处接入。
2. 美股纸上交易（阶段 2 美股部分 + 阶段 3 美股部分）作为第一条完整链路。
3. A 股数据（Tushare）和纸上撮合（阶段 2 A 股部分 + 阶段 3 A 股纸上）。
4. 阻塞项解决后，再做 A 股实盘适配。
