# mermaid-lint

一个用于检查和修复 Markdown 文件中 Mermaid 图表常见问题的 Go 静态分析工具。

## 功能特性

- 递归扫描目录中的 Markdown 文件
- 检测 8 种常见的 Mermaid 问题
- 自动修复可修复的问题
- 支持文本和 JSON 两种输出格式
- 可通过严格模式集成到 CI/CD 流水线

## 安装

```bash
go build -o mermaid-lint main.go
```

## 使用方法

```bash
./mermaid-lint <目录> [选项]

选项:
  --fix       自动修复可修复的问题
  --dry-run   仅显示更改，不修改文件
  --json      以 JSON 格式输出结果
  --strict    如果发现问题，以非零码退出
```

## 检测问题

| 问题类型 | 可修复 | 描述 |
|---------|-------|------|
| newline_literal | ✓ | 使用 `\n` 而不是 `<br>` 换行 |
| unquoted_text | ✓ | 包含特殊字符 `():,` 的文本未加引号 |
| html_literal | | 图表中使用了 HTML 标签 |
| duplicate_node | | 重复定义了节点 ID |
| undefined_class | | 使用了未定义的样式类 |
| invalid_style | | 无效的样式属性名 |
| isolated_node | | 存在无连接的孤立节点 |
| duplicate_subgraph | | 重复定义了子图名称 |

## 示例

```bash
# 检查当前目录
./mermaid-lint .

# 检查并自动修复
./mermaid-lint ./docs --fix

# CI 用法 - 发现问题时构建失败
./mermaid-lint ./docs --strict
```

## 许可证

MIT
