package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const Version = "0.0.1-beta"
const QtdExecutions = 5
const DelaySeconds = 10
const PathFileSites = "sites.txt"
const PathLogs = "logs/"

func main() {
	for {
		fmt.Println("Sites Monitor v" + Version)
		fmt.Println("1 - Start monitoring")
		fmt.Println("2 - Show logs files")
		fmt.Println("3 - Show log")
		fmt.Println("0 - Exit")
		fmt.Println("Entre com a opção:")

		comand := 0
		fmt.Scan(&comand)
		fmt.Println("\n\n===================================")

		switch comand {
		case 0:
			os.Exit(0)
		case 1:
			monitoring()
		case 2:
			showLogsFiles()
		case 3:
			fmt.Println("Enter the name of the log file: <year_month_day.log>")
			var logFile = "nil.log"
			fmt.Scan(&logFile)
			showLog(logFile)
		default:
			fmt.Println("Invalid command")
		}
		fmt.Println("\n===================================\n\n")
	}
}

func showLog(fileName string) {
	data, err := getDataFromFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data)
}

func monitoring() {
	sites, err := getSitesFromFile(PathFileSites)

	if nil != err {
		fmt.Println(err)
		return
	}

	if len(sites) <= 0 {
		fmt.Println("Lista de sites vazia")
		return
	}

	for i := 0; i < QtdExecutions; i++ {
		for _, site := range sites {
			siteTest(site)
		}

		time.Sleep(DelaySeconds * time.Second)
	}
}

func getSitesFromFile(path string) ([]string, error) {
	var sites []string
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	for {
		linha, err := reader.ReadString('\n')
		linha = strings.TrimSpace(linha)
		if len(linha) > 0 {
			sites = append(sites, linha)
		}
		if err == io.EOF {
			break
		}
	}

	_ = file.Close()

	return sites, err
}

func siteTest(site string) {
	resp, err := http.Get(site)
	statusCode := -1
	if resp != nil {
		statusCode = resp.StatusCode
	}
	registerLog(statusCode, site, err)
}

func registerLog(statusCode int, site string, err error) {
	var nameFileLog = time.Now().Format("2006_01_02") + ".log"
	fileLog, errLog := getFileLog(nameFileLog)
	if errLog != nil {
		fmt.Sprintf("Erro ao abrir o log {", errLog.Error(), "}")
	}

	var stringToWrite = "{'date_time':%s,'status_code':%d, 'site':'%s'}"
	formatedDate := time.Now().Format("2006-01-02T15:04:05.999-07:00")
	if err == nil {
		stringToWrite = fmt.Sprintf(stringToWrite, formatedDate, statusCode, site)
	} else {
		stringToWrite = "{'date_time':" + formatedDate + "','error':'" + err.Error() + "'}"
	}
	fmt.Println(stringToWrite)
	if fileLog != nil {
		fileLog.WriteString(stringToWrite)
		fileLog.WriteString("\n")
		fileLog.Close()
	}

}

func getFileLog(nameFile string) (arquivo *os.File, err error) {
	var pathFile = fmt.Sprintf("%s/%s", PathLogs, nameFile)

	err = os.MkdirAll(PathLogs, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(pathFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
}

func getDataFromFile(fileLog string) (data string, err error) {
	fileReadBytes, err := ioutil.ReadFile(PathLogs + fileLog)
	if err != nil {
		return "", err
	}
	return string(fileReadBytes), nil
}

func showLogsFiles() {
	files, err := ioutil.ReadDir(PathLogs)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
