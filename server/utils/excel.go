package utils

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

// ExcelResult 导出的excel结果
type ExcelResult struct {
	Page          bool                   // 是否需要分页导出 TODO: 注意 暂时只支持单sheet
	SheetNameList []string               // sheet name 集合
	SheetMap      map[string]interface{} // sheet 数据map，sheetName 为键
	Headers       []string               //自定义表头
	ExportLog     interface{}            //导出log
}

func parseData(data interface{}) []interface{} {
	arr := make([]interface{}, 0)
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// 如果是切片集合
	if val.Kind() == reflect.Slice {
		for j := 0; j < val.Len(); j++ {
			item := val.Index(j).Interface()
			arr = append(arr, item)
		}
	} else {
		// 如果只有一条数据
		arr = append(arr, data)
	}
	return arr
}

// WriteExcelSheet
// sheetName 写sheet页
// data 要写的数据
// rowNum 从该sheet的第几行开始写
// headerParams 动态表头参数
func WriteExcelSheet(f *excelize.File, sheetName string, headers []string, data interface{},
	rowNum int, headerParams map[string]string) int {
	arr := parseData(data)

	if len(arr) == 0 {

		return rowNum
	}

	firstObj := arr[0]

	// 写header
	headers = queryHeaders(firstObj, headers)

	colorList := ListTagValue(arr[0], "color")
	if headerParams != nil {
		for i := range headers {
			headers[i] = formatDynamicHeader(headers[i], headerParams)
		}
	}

	// 设置表头
	rowNum = setHeader(f, sheetName, headers, colorList, rowNum)

	//保存每列文字的最大长度
	colMapMaxLength := make(map[int]int)

	// 收集表头文本长度 Todo:多行表头
	for i, header := range headers {
		colMapMaxLength[i] = len(header)
	}

	// 获取单元格样式
	bodyStyle := createDefaultBodyStyle(f)
	lastCol := convertNumToChars(len(headers))
	headerSet := make(map[string]bool)
	for _, item := range headers {
		headerSet[item] = true
	}
	//保存内容
	for _, obj := range arr {
		//读取所有标注xlsx的字段
		xlsxValues := ListFieldWithTagByCustom(obj, "xlsx", headers)

		length := len(xlsxValues)
		//读取值
		rowValues := make([]interface{}, 0)
		for i := 0; i < length; i++ {

			value := cast.ToString(xlsxValues[i].Interface())
			// 枚举 可能读不到  so  取下 String()
			if value == "" {
				value = xlsxValues[i].String()
			}
			l := len(value)
			// 覆盖最大长度
			if colMapMaxLength[i] < l {
				colMapMaxLength[i] = l
			}
			rowValues = append(rowValues, xlsxValues[i].Interface())
		}
		f.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &rowValues)
		f.SetCellStyle(sheetName, convertCell("A", rowNum), convertCell(lastCol, rowNum), bodyStyle)
		rowNum++
	}

	// 设置计算默认列宽
	setColWidth(f, sheetName, colMapMaxLength)

	return rowNum
}

func queryHeaders(obj interface{}, customHeaders []string) []string {
	result := make([]string, 0)
	if obj == nil {
		return result
	}
	//如果没传自定义表头   则用默认表头
	if len(customHeaders) == 0 {
		return ListTagValue(obj, "xlsx")
	}

	//如果有传自定义表头
	customHeaders = ListTagValueByCustom(obj, "xlsx", customHeaders)

	return customHeaders
}

func formatDynamicHeader(format string, headerParams map[string]string) string {
	tmpl, _ := template.New("").Parse(format)
	s := new(strings.Builder)
	tmpl.Execute(s, headerParams)
	return s.String()
}

func setColWidth(f *excelize.File, sheetName string, colMapMaxLength map[int]int) {
	for col, length := range colMapMaxLength {
		startCol := convertNumToChars(col + 1)
		newLength := convertLength(length)
		f.SetColWidth(sheetName, startCol, startCol, newLength)
	}
}

func convertLength(length int) float64 {
	// 字长+2
	num := float64(length + 2)
	// 最长100
	if num > 100 {
		return 100
	}
	return num
}

func setHeader(f *excelize.File, sheetName string, tagValueList, colorList []string, rowNum int) int {
	// 获取最后一列值
	lastCol := convertNumToChars(len(tagValueList))

	hasColor := len(colorList) > 0

	tempList := make([][]string, len(tagValueList))

	// 按列分组的header 二维数组
	row := 1
	for index, tag := range tagValueList {
		tempList[index] = strings.Split(tag, " > ")
		if len(tempList[index]) > row {
			row = len(tempList[index])
		}
	}

	col := len(tagValueList)
	headerList := make([][]string, row)

	//TODO 校验header 合法性
	//转换成按行分组的header 二维数组
	for i := 0; i < row; i++ {

		headerList[i] = make([]string, col)

		for j := 0; j < col; j++ {
			length := len(tempList[j])
			// 同一个表中 可能部分字段二级表头，部分字段三级表头，存在数组越界可能
			if i > length-1 {
				headerList[i][j] = tempList[j][length-1]
				continue
			}
			headerList[i][j] = tempList[j][i]
		}

	}

	for i, header := range headerList {
		//写表头
		f.SetSheetRow(sheetName, fmt.Sprintf("A%d", rowNum), &header)

		//合并单元格
		lastHeader := header[0]
		start, end := 0, 0
		for index, h := range header {
			if h != lastHeader {
				end = index

				startCol := convertNumToChars(start + 1)
				endCol := convertNumToChars(end)

				if hasColor {
					// 设置合并单元格的样式
					setHeaderStyle(f, sheetName, colorList[start], startCol, endCol, rowNum)
				}

				//合并
				//判断上下列是否同表头
				mergeCell(i, rowNum, start, end, header, headerList, f, sheetName)

				start = index
				lastHeader = h
			}
		}

		//合并最后一个
		if end != 0 {
			end = len(header)
			startCol := convertNumToChars(start + 1)
			endCol := convertNumToChars(end)
			if hasColor {
				// 设置合并单元格的样式
				setHeaderStyle(f, sheetName, colorList[start], startCol, endCol, rowNum)
			}
			//判断上下列是否同表头
			mergeCell(i, rowNum, start, end, header, headerList, f, sheetName)
		}

		//写默认样式
		if !hasColor {
			setDefaultHeaderStyle(f, sheetName, lastCol, rowNum)
		}

		rowNum++
	}
	return rowNum
}

func mergeCell(index, rowNum, start, end int, header []string, headerList [][]string, f *excelize.File, sheetName string) {
	startCol := convertNumToChars(start + 1)
	endCol := convertNumToChars(end)
	//合并
	//判断上下行是否同表头
	merged := false
	if index != 0 {
		// 计算有多少行需要合并
		mergerRowNum := 0
		for j := 1; j <= index; j++ {
			// 比较不同行表头是否相同
			if header[end-1] == headerList[index-j][end-1] && header[start] == headerList[index-j][start] {
				mergerRowNum++
			}
		}
		if mergerRowNum != 0 {
			// 合并上下行+列
			f.MergeCell(sheetName, fmt.Sprintf(startCol+"%d", rowNum-mergerRowNum), fmt.Sprintf(endCol+"%d", rowNum))
			merged = true
		}

	}
	// 否则合并同行
	if !merged {
		f.MergeCell(sheetName, fmt.Sprintf(startCol+"%d", rowNum), fmt.Sprintf(endCol+"%d", rowNum))
	}
}

func WriteDefaultExcelSheet(f *excelize.File, data interface{}) int {
	return WriteExcelSheet(f, "Sheet1", nil, data, 1, nil)
}

// setHeaderStyle  设置header 样式
func setHeaderStyle(f *excelize.File, sheetName, fillColor, startCol, endCol string, rowNum int) {
	// 创建表头样式
	titleStyle := createHeaderStyle(f, fillColor)

	// 设置表头样式
	f.SetCellStyle(sheetName, convertCell(startCol, rowNum), convertCell(endCol, rowNum), titleStyle)
}

// setHeaderStyle  设置header 样式
func setDefaultHeaderStyle(f *excelize.File, sheetName, colStr string, rowNum int) {
	// 创建表头样式
	titleStyle := createHeaderStyle(f, "E6F4EA")

	// 设置表头样式
	f.SetCellStyle(sheetName, convertCell("A", rowNum), convertCell(colStr, rowNum), titleStyle)
}

func convertCell(col string, rowNum int) string {
	return fmt.Sprintf("%s%d", col, rowNum)
}

func createHeaderStyle(f *excelize.File, fillColor string) int {

	style := &excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#37424e", Style: 2},
			{Type: "top", Color: "#37424e", Style: 2},
			{Type: "bottom", Color: "#37424e", Style: 2},
			{Type: "right", Color: "#37424e", Style: 2},
		},
		Fill: excelize.Fill{
			// gradient： 渐变色    pattern   填充图案
			Pattern: 1, // 填充样式  当类型是 pattern 0-18 填充图案  1 实体填充
			Type:    "pattern",
			Color:   []string{fillColor}, //"E6F4EA"
		}, Font: &excelize.Font{
			Bold: true,
		}, Alignment: &excelize.Alignment{
			// 水平对齐方式 center left right fill(填充) justify(两端对齐)  centerContinuous(跨列居中) distributed(分散对齐)
			Horizontal: "left",
			// 垂直对齐方式 center top  justify distributed
			Vertical: "center",
			// Indent:     1,        // 缩进  只要有值就变成了左对齐 + 缩进
			// TextRotation: 30, // 旋转
			// RelativeIndent:  10,   // 好像没啥用
			// ReadingOrder:    0,    // 不知道怎么设置
			// JustifyLastLine: true, // 两端分散对齐，只有 水平对齐 为 distributed 时 设置true 才有效
			// WrapText:        true, // 自动换行
			// ShrinkToFit:     true, // 缩小字体以填充单元格
		},
	}

	titleStyle, err := f.NewStyle(style)
	if err != nil {
		fmt.Println(err)
	}
	return titleStyle
}

func createDefaultBodyStyle(f *excelize.File) int {
	style := &excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#37424e", Style: 1},
			{Type: "top", Color: "#37424e", Style: 1},
			{Type: "bottom", Color: "#37424e", Style: 1},
			{Type: "right", Color: "#37424e", Style: 1},
		},
		Alignment: &excelize.Alignment{
			// 水平对齐方式 center left right fill(填充) justify(两端对齐)  centerContinuous(跨列居中) distributed(分散对齐)
			Horizontal: "left",
			// 垂直对齐方式 center top  justify distributed
			Vertical: "center",
		},
	}

	bodyStyle, err := f.NewStyle(style)
	if err != nil {
		fmt.Println(err)
	}
	return bodyStyle
}

func ReadExcelByReader(reader io.Reader, holder interface{}) ([]string, []string, error) {
	ef, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, nil, err
	}
	return ReadExcel(ef, holder)
}

var ExcelChar = []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func convertNumToChars(num int) string {
	var cols string
	v := num
	for v > 0 {
		k := v % 26
		if k == 0 {
			k = 26
		}
		v = (v - k) / 26
		cols = ExcelChar[k] + cols
	}
	return cols
}

func ReadExcel(ef *excelize.File, holder interface{}) ([]string, []string, error) {
	//默认读Sheet1
	rows, _ := ef.GetRows("Sheet1") // 所有行
	var err error
	// 没有则读第一个sheet
	if len(rows) == 0 {
		sheetName := ef.GetSheetName(0)
		rows, err = ef.GetRows(sheetName) // 所有行
		if err != nil {
			return nil, nil, err
		}
	}
	// 获取第一行表头
	header := rows[0]
	objHeaders := make([]string, 0)
	// 记录每一个表头对应的列
	mapIndex := make(map[string]int)
	for index, h := range header {
		mapIndex[h] = index
	}

	for index, row := range rows {
		if row == nil {
			// 跳过空行
			continue
		}
		tp := reflect.TypeOf(holder).Elem().Elem() // 结构体的类型
		val := reflect.New(tp)                     // 创建一个新的结构体对象
		valType := val.Type()
		if valType.Kind() == reflect.Ptr {
			valType = valType.Elem()
		}
		hasData := false
		for i := 0; i < val.Elem().NumField(); i++ {
			field := val.Elem().Field(i) // 第idx个字段的反射Value
			fieldStruct := valType.Field(i)
			xlsxTag := fieldStruct.Tag.Get("xlsx")
			if xlsxTag != "" {
				objHeaders = append(objHeaders, xlsxTag)
				if rowIndex, ok := mapIndex[xlsxTag]; ok {
					// 假设excel模板有10列，用户只填前5列，这里读到的该行数据将只有5列
					// 不能数组越界
					if rowIndex < len(row) {
						cellValue := strings.TrimSpace(row[rowIndex]) // 第idx个字段对应的Excel数据
						field.SetString(cellValue)                    // 将Excel数据保存到结构体对象的对应字段中
						hasData = true
					}
				}
			}
		}
		// 有数据才添加进list
		if hasData && index != 0 {
			listV := reflect.ValueOf(holder)
			listV.Elem().Set(reflect.Append(listV.Elem(), val.Elem())) // 将结构体对象添加到holder中
		}
	}
	return objHeaders, header, nil
}
