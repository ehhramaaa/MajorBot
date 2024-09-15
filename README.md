[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/major/start?startapp=5024522783)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/bansos_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/bansos_code_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/MajorBot/main/assets/Screenshot_1.png)

## Recommendation before use

# 🔥🔥 Go Version Tested 1.23.1 🔥🔥

## Features

|        Feature         | Supported |
| :--------------------: | :-------: |
|     Multithreading     |    ✅     |
|     Use Query Data     |    ✅     |
|  Auto Completing Task  |    ✅     |
| Auto Play Swipe Coins  |    ✅     |
|  Auto Play Hold Coins  |    ✅     |
| Auto Play Durov Puzzle |    ❌     |
|  Auto Connect Wallet   |    ✅     |
|   Random User Agent    |    ✅     |

## [Settings](https://github.com/ehhramaaa/MajorBot/blob/main/config.yml)

|     Settings     |                          Description                          |
| :--------------: | :-----------------------------------------------------------: |
| **SWIPE_COINS**  |      Amount coins in Swipe Coin (e.g. MIN:500, MAX:600)       |
|  **HOLD_COINS**  |      Amount coins in Hold Coin (e.g. MIN:1500, MAX:2000)      |
|  **MAX_THREAD**  |       Max Thread Worker Run Parallel Recommend 10 - 100       |
| **RANDOM_SLEEP** | Delay before the next lap (e.g. MIN:3600, MAX:7200) in second |

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) **version > 1.22**

## Installation

You can download the [**repository**](https://github.com/ehhramaaa/agent301.git) by cloning it to your system and installing the necessary dependencies:

```shell
git clone https://github.com/ehhramaaa/MajorBot.git
cd MajorBot
go mod tidy
go run .
```

Then you can do build application by typing:

Windows:

```shell
go build -o MajorBot.exe
```

Linux:

```shell
go build -o MajorBot
```

## Usage

```shell
go run .
```

Or

```shell
go run main.go
```

**If You Want Auto Select Choice In Terminal**

For Option 1

```shell
go run . -c 1
```

For Option 2

```shell
go run . -c 2
```
