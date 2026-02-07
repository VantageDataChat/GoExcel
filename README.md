> [English](README_en.md) | [中文](README.md)

# GoExcel

纯 Go 电子表格处理库，灵感来自 [PHPOffice/PhpSpreadsheet](https://github.com/PHPOffice/PhpSpreadsheet)。

支持 XLSX 和 CSV 格式的读写，提供内存中的电子表格对象模型，包含公式计算引擎、样式系统、条件格式、数据验证等功能。

## 特性

- XLSX 读写（Open XML 格式，兼容 Excel 2007+）
- CSV 读写（自定义分隔符）
- 公式计算引擎（24 个内置函数）
- 完整的样式系统（字体、边框、填充、对齐、数字格式）
- 合并单元格、冻结窗格
- 行列插入/删除/复制
- 超链接、批注、富文本
- 条件格式、数据验证、自动筛选
- 页面设置与打印配置
- 工作表/工作簿保护
- 文档属性、命名范围
- Unicode 和特殊字符支持
- 97.7% 测试覆盖率，256 个测试用例

## 安装

```bash
go get github.com/VantageDataChat/GoExcel
```

## 快速开始

### 创建并写入 XLSX

```go
package main

import "github.com/VantageDataChat/GoExcel"

func main() {
    wb := gospreadsheet.New()
    ws := wb.GetActiveSheet()

    // 设置单元格值
    ws.SetCellValue("A1", "姓名")
    ws.SetCellValue("B1", "分数")
    ws.SetCellValue("A2", "Alice")
    ws.SetCellValue("B2", 95.5)
    ws.SetCellValue("A3", "Bob")
    ws.SetCellValue("B3", 87)

    // 添加公式
    ws.SetCellFormula("B4", "AVERAGE(B2:B3)")

    // 设置样式
    ws.SetCellStyle("A1", gospreadsheet.NewStyle().
        SetFont(&gospreadsheet.Font{Bold: true, Size: 14}))

    // 保存
    gospreadsheet.SaveFile(wb, "output.xlsx")
}
```

### 读取 XLSX

```go
wb, err := gospreadsheet.OpenFile("input.xlsx")
if err != nil {
    log.Fatal(err)
}

ws := wb.GetActiveSheet()
val, _ := ws.GetCellValue("A1")
fmt.Println(val)

// 遍历所有行
rows, _ := ws.RowIterator()
for _, row := range rows {
    for _, cell := range row {
        if cell != nil {
            fmt.Print(cell.GetStringValue(), "\t")
        }
    }
    fmt.Println()
}
```

### CSV 读写

```go
// 读取 CSV
wb, _ := gospreadsheet.OpenFile("data.csv")

// 自定义分隔符
reader := gospreadsheet.NewCSVReader()
reader.Delimiter = ';'
wb, _ = reader.Open("data.csv")

// 写入 CSV
gospreadsheet.SaveFile(wb, "output.csv")
```

### 公式计算

```go
wb := gospreadsheet.New()
ws := wb.GetActiveSheet()
ws.SetCellValue("A1", 100)
ws.SetCellValue("A2", 200)
ws.SetCellFormula("A3", "SUM(A1:A2)")

ce := gospreadsheet.NewCalculationEngine(wb)
result, _ := ce.CalculateCell(ws, "A3")
fmt.Println(result) // 300

// 计算所有公式
ce.CalculateAll()
```

### 样式

```go
style := gospreadsheet.NewStyle().
    SetFont(&gospreadsheet.Font{
        Name: "Arial", Size: 12, Bold: true, Color: "FF0000",
    }).
    SetFill(&gospreadsheet.Fill{
        Type: "solid", Color: "FFFF00",
    }).
    SetBorders(&gospreadsheet.Borders{
        Bottom: gospreadsheet.Border{
            Style: gospreadsheet.BorderThin, Color: "000000",
        },
    }).
    SetAlignment(&gospreadsheet.Alignment{
        Horizontal: gospreadsheet.AlignCenter,
        WrapText:   true,
    }).
    SetNumberFormat(&gospreadsheet.FormatPercent2Dec)

ws.SetCellStyle("A1", style)
```

### 多工作表

```go
wb := gospreadsheet.New()
ws1 := wb.GetActiveSheet() // Sheet1

ws2, _ := wb.AddSheet("数据")
ws2.SetCellValue("A1", "第二个工作表")

ws3, _ := wb.AddSheet("汇总")
wb.SetActiveSheet(2) // 切换到第三个工作表
```

### 合并单元格与冻结窗格

```go
ws.SetCellValue("A1", "合并标题")
ws.MergeCells("A1:D1")
ws.FreezePane("A2") // 冻结首行
```

### 行列操作

```go
ws.InsertRow(2)       // 在第3行前插入空行
ws.DeleteRow(5)       // 删除第6行
ws.InsertColumn(1)    // 在B列前插入空列
ws.DeleteColumn(3)    // 删除D列
ws.CopyRow(0, 10)     // 复制第1行到第11行
ws.SetColumnWidth(0, 20.0)
ws.SetRowHeight(0, 30.0)
```

### 超链接与批注

```go
ws.SetCellValue("A1", "点击访问")
ws.SetCellHyperlink("A1", "https://example.com")
ws.SetCellComment("B1", "审核员", "请检查此数据")
```

### 条件格式

```go
cf := gospreadsheet.NewConditionalFormatting("B2:B100")
cf.AddRule(gospreadsheet.CellIsRule(
    gospreadsheet.OperatorGreaterThan, "90",
    gospreadsheet.NewStyle().SetFont(&gospreadsheet.Font{Color: "00AA00"}),
))
cf.AddRule(gospreadsheet.BetweenRule("60", "90",
    gospreadsheet.NewStyle().SetFont(&gospreadsheet.Font{Color: "FF8800"}),
))
ws.AddConditionalFormatting(cf)
```

### 数据验证

```go
// 下拉列表
dv := gospreadsheet.NewDataValidation("C2:C100")
dv.SetListValues([]string{"通过", "未通过", "待定"})
dv.SetErrorMessage("错误", "请从列表中选择")
ws.AddDataValidation(dv)

// 数值范围
dv2 := gospreadsheet.NewDataValidation("D2:D100").
    SetType(gospreadsheet.ValidationWhole).
    SetOperator(gospreadsheet.ValOperatorBetween).
    SetFormula1("0").SetFormula2("100")
ws.AddDataValidation(dv2)
```

### 页面设置

```go
ps := ws.GetPageSetup()
ps.SetOrientation(gospreadsheet.OrientationLandscape)
ps.SetPaperSize(gospreadsheet.PaperA4)
ps.SetScale(85)
ps.SetPrintArea("A1:H50")
ps.SetRepeatRows("1:1") // 每页重复标题行
```

### 工作表保护

```go
sp := gospreadsheet.NewSheetProtection().SetPassword("secret")
sp.AllowSort()
sp.AllowAutoFilter()
ws.SetSheetProtection(sp)
```

## 支持的公式函数

| 类别 | 函数 |
|------|------|
| 数学 | SUM, ABS, ROUND, SQRT, POWER, MOD, INT |
| 统计 | AVERAGE, COUNT, COUNTA, MAX, MIN, MEDIAN |
| 逻辑 | IF |
| 文本 | LEN, UPPER, LOWER, TRIM, LEFT, RIGHT, MID, CONCATENATE |
| 条件 | SUMIF, COUNTIF |

支持单元格引用（A1）、范围引用（A1:A10）、嵌套公式和四则运算。

## 项目结构

```
goexcel/
├── spreadsheet.go     # Cell 类型与操作
├── workbook.go        # Workbook 工作簿管理
├── worksheet.go       # Worksheet 工作表
├── coordinate.go      # 单元格坐标解析
├── style.go           # 样式系统
├── calculation.go     # 公式计算引擎
├── functions.go       # 内置函数实现
├── xlsx_writer.go     # XLSX 写入
├── xlsx_reader.go     # XLSX 读取
├── csv.go             # CSV 读写
├── iofactory.go       # 格式自动检测
├── hyperlink.go       # 超链接
├── comment.go         # 批注
├── richtext.go        # 富文本
├── conditional.go     # 条件格式
├── validation.go      # 数据验证
├── autofilter.go      # 自动筛选
├── pagesetup.go       # 页面设置
└── protection.go      # 工作表/工作簿保护
```

## 测试

```bash
go test ./... -v
go test ./... -cover  # 覆盖率 97.7%
```

## 许可证

MIT
