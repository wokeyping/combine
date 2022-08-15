package main

import (
	"bufio"
	"fmt"
	"github.com/Luxurioust/excelize"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	f, err := excelize.OpenFile("./combine/combine.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	//读取某个单元格的值
	/*value, err := f.GetCellValue("combine", "D2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)*/
	//读取某个表单的所有数据
	rows, err := f.GetRows("combine")
	if err != nil {
		log.Fatal(err)

	}
	var sevs = make(map[string][]string)
	for kk, row := range rows {
		//fmt.Printf("\t%s", row)
		if kk == 0 {
			continue
		}
		for k, _ := range row {
			if k > 0 {
				continue
			}
			row[k] = strings.Trim(row[k], "区")
			sevs[row[k+1]] = append(sevs[row[k+1]], row[k])
			//fmt.Printf("\t%s", row[k+1])

		}
		//fmt.Println()
	}
	sort.Strings(sevs)
	filePath := "F:/Goproject/src/combine/combine.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	//及时关闭file句柄
	defer file.Close()
	write := bufio.NewWriter(file)

	for _, value := range sevs {
		write.WriteString("'")
		write.WriteString(strings.Join(value, ","))
		write.WriteString("',\r\n")
		//fmt.Printf("index:%s value:%v\r\n", index, value)
	}
	write.WriteString("\r\n\r\n\r\n===============export redis===================\r\n")
	export := []string{}
	for _, value := range sevs {
		for _, v := range value {
			export = append(export, v)
		}
	}
	write.WriteString("php exportHeRedis.php " + strings.Join(export, ",") + " > /tmp/he20220429 2>&1" + "\r\n")

	var buidetable []string
	write.WriteString("\r\n\r\n\r\n===============built table===================\r\n")
	for k, _ := range sevs {
		buidetable = append(buidetable, k)
	}
	write.WriteString("php buidetable.php " + strings.Join(buidetable, ",") + "\r\n")

	write.WriteString("\r\n\r\n\r\n===============merge name===================\r\n")
	for key, value := range sevs {
		for _, v := range value {
			if key == v {
				continue
			}
			write.WriteString("php mergeNameYn.php " + v + "\r\n")
		}
	}

	write.WriteString("\r\n\r\n\r\n===============php start===================\r\n")
	for _, value := range sevs {
		for _, v := range value {
			write.WriteString("php start.php " + v + "\r\n")
		}
	}

	write.Flush()
}
