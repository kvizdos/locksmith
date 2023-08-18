package logger

import (
	"fmt"
	"strings"
)

func LogFormatToXML(lf LogFormats) string {
	format := logs[lf]

	format.WazuhName = "locksmith"
	format.Parent = "locksmith"

	xmlStr := fmt.Sprintf("<decoder name=\"locksmith_%s\">\n", lf)
	xmlStr += fmt.Sprintf("    <parent>%s</parent>\n", format.Parent)
	xmlStr += fmt.Sprintf("    <regex>%s</regex>\n", format.Regex)
	xmlStr += fmt.Sprintf("    <order>%s</order>\n", strings.Join(format.RegexOrder, ","))
	xmlStr += "</decoder>"

	fmt.Println(xmlStr)
	return xmlStr
}
