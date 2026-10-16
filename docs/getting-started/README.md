# Getting Started

Crux runs on your own server. You add an AI model and a trading account in the web app, then create a strategy and a trader.

## 1. Install

See the [README](../../README.md#install). The web app runs at `http://127.0.0.1:3000` by default.

## 2. Register

The first account created becomes the owner of the instance.

## 3. Add an AI model

Any model with an OpenAI-compatible chat completions endpoint works. Enter the base URL, model name and API key in the AI model settings.

## 4. Add a trading account

| Account type | Market | What you need |
| :--- | :--- | :--- |
| `alpaca_paper` | US stocks | An Alpaca paper trading key and secret |
| `cn_paper` | China A-shares | A Tushare Pro token |

Both markets run in paper trading mode. See the [Markets](../../README.md#markets) section of the README for the rules that apply to each market.

## 5. Create a strategy and a trader

Create a strategy in **Strategy Studio** with the symbols and the risk limits. Then create a trader that uses the strategy, the AI model and the account, and start it.

A trader only accepts symbols that match its market. For example, `AAPL` is valid for `us_equity` and `600519.SH` is valid for `cn_equity`.
