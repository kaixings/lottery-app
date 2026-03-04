# 🎱 lottery-app

> 双色球 & 大乐透号码预测工具 —— 基于历史数据的多算法选号辅助程序

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)]()

---

## 功能特性

- **双彩种支持**：双色球（SSQ）和大乐透（DLT）独立数据管理
- **实时抓取**：从 500彩票网自动拉取最新历史开奖数据
- **四种预测算法**：频率、冷热号、随机、加权随机
- **交互式 TUI**：基于 Bubble Tea 的全键盘操作终端界面
- **CLI 模式**：支持脚本化调用的命令行子命令

---

## 截图预览

```
┌─────────────────────────────┐
│       彩票预测工具           │
│                             │
│  请选择彩种：               │
│                             │
│  > 双色球                   │
│    大乐透                   │
│                             │
│  ↑↓ 移动  Enter 确认  q 退出 │
└─────────────────────────────┘
```

```
┌─────────────────────────────┐
│     双色球预测工具           │
│                             │
│  > 抓取数据                 │
│    查看历史                 │
│    预测号码                 │
│    返回                     │
│    退出                     │
│                             │
│  ↑↓ 移动  Enter 确认  q 退出 │
└─────────────────────────────┘
```

---

## 快速开始

### 前置要求

- Go 1.24+

### 安装

```bash
git clone https://github.com/yourname/lottery-app.git
cd lottery-app
go build -o lottery .
```

### 启动 TUI（推荐）

```bash
./lottery
```

程序启动后先选择彩种（双色球 / 大乐透），进入对应功能菜单。

---

## CLI 命令参考

> 以下命令均以双色球（SSQ）为例，操作大乐透数据请通过 TUI 进行。

### 抓取历史数据

```bash
# 抓取最近 30 期（默认）
lottery fetch

# 抓取最近 100 期
lottery fetch -n 100
```

### 查看历史记录

```bash
# 显示最近 10 期（默认）
lottery list

# 显示最近 20 期
lottery list -n 20
```

输出示例：

```
期号            日期            红球                       蓝球
-----------------------------------------------------------------
2024001         2024-01-02      03 07 12 18 24 31          08
2023153         2023-12-30      05 11 14 21 27 33          12
...
共 100 期数据
```

### 预测号码

```bash
# 使用所有算法各预测 1 注
lottery predict

# 指定算法
lottery predict -a frequency
lottery predict -a hot-cold
lottery predict -a random
lottery predict -a weighted

# 生成 5 注
lottery predict -a weighted -t 5
```

输出示例：

```
基于 100 期历史数据预测

【frequency】频率最高的号码
  红球: 03 07 12 18 24 31  蓝球: 08

【hot-cold】冷热号组合
  红球: 05 11 17 22 28 33  蓝球: 11

【random】纯随机
  红球: 01 09 15 20 26 32  蓝球: 04

【weighted】加权随机
  红球: 06 13 19 23 27 30  蓝球: 09
```

### 全局选项

```bash
# 指定数据目录（默认 ./data）
lottery --data /path/to/data predict
```

---

## 算法说明

| 算法 | 名称 | 说明 |
|------|------|------|
| `frequency` | 频率算法 | 统计全量历史数据，选取出现次数最多的号码 |
| `hot-cold` | 冷热号算法 | 近期高频号码（热号）+ 长期未出现号码（冷号）组合 |
| `random` | 随机算法 | 在合法号码范围内完全随机抽取 |
| `weighted` | 加权随机 | 按历史出现频率赋权，频率越高被选中概率越大 |

> ⚠️ 所有预测结果均基于统计规律，不构成任何购彩建议。彩票开奖结果为随机事件，历史数据对未来开奖无预测效力。

---

## 彩种规则

| | 双色球 | 大乐透 |
|---|---|---|
| 号码1 | 红球：6个，范围 1–33 | 前区：5个，范围 1–35 |
| 号码2 | 蓝球：1个，范围 1–16 | 后区：2个，范围 1–12 |
| 数据来源 | 500彩票网 SSQ 历史 | 500彩票网 DLT 历史 |
| 本地存储 | `data/lottery.json` | `data/dlt.json` |

---

## 项目结构

```
lottery-app/
├── main.go
├── cmd/
│   ├── root.go          # CLI 入口，TUI 启动
│   ├── fetch.go         # fetch 子命令（双色球）
│   ├── list.go          # list 子命令（双色球）
│   └── predict.go       # predict 子命令（双色球）
├── internal/
│   ├── storage/
│   │   ├── storage.go   # 双色球存储（LotteryRecord / Store）
│   │   └── dlt.go       # 大乐透存储（DLTRecord / DLTStore）
│   ├── fetcher/
│   │   ├── fetcher.go   # 双色球抓取器
│   │   └── dlt.go       # 大乐透抓取器
│   ├── algorithm/
│   │   ├── algorithm.go       # SSQ 算法接口 & 注册
│   │   ├── frequency.go       # SSQ 频率算法
│   │   ├── hotcold.go         # SSQ 冷热号算法
│   │   ├── random.go          # SSQ 随机 & 加权算法
│   │   ├── dlt.go             # DLT 算法接口 & 注册
│   │   ├── dlt_frequency.go   # DLT 频率算法
│   │   ├── dlt_hotcold.go     # DLT 冷热号算法
│   │   └── dlt_random.go      # DLT 随机 & 加权算法
│   └── tui/
│       ├── model.go           # 顶层模型 & 视图路由
│       ├── menu.go            # 功能菜单视图
│       ├── lottery_select.go  # 彩种选择视图
│       ├── fetch.go           # SSQ 抓取视图
│       ├── history.go         # SSQ 历史视图
│       ├── predict.go         # SSQ 预测视图
│       ├── dlt_fetch.go       # DLT 抓取视图
│       ├── dlt_history.go     # DLT 历史视图
│       ├── dlt_predict.go     # DLT 预测视图
│       └── styles.go          # 全局样式定义
└── data/
    ├── lottery.json     # 双色球本地历史数据
    └── dlt.json         # 大乐透本地历史数据
```

---

## 依赖

| 包 | 用途 |
|---|---|
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | TUI 框架 |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | TUI 样式 |
| [spf13/cobra](https://github.com/spf13/cobra) | CLI 框架 |

---

## License

MIT
