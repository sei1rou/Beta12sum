package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type k27804 struct {
	head string
	kd   [5]kdata
}

type kdata struct {
	kCode  string
	kValue string
	kFlag  string
	kCmt   string
}

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	//ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Print("Start\r\n")

	// ファイルの読み込み準備
	infile, err := os.Open(flag.Arg(0))
	failOnError(err)
	defer infile.Close()

	// ファイルの書き込み準備
	writeFileDir, _ := filepath.Split(flag.Arg(0))
	writeFileDir = writeFileDir + "K27804B.DAT"
	outfile, err := os.Create(writeFileDir)
	failOnError(err)
	defer outfile.Close()

	var ik27804 k27804
	//ik27804の初期化
	ik27804.head = ""
	for i := 0; i < 5; i++ {
		ik27804.kd[i].kCode = ""
		ik27804.kd[i].kValue = ""
		ik27804.kd[i].kFlag = ""
		ik27804.kd[i].kCmt = ""
	}

	// ファイルの読み込み
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		ik27804.head = line[:78]
		pos := 78
		B1 := 0.0
		B2 := 0.0
		for i := 0; i < 5; i++ {
			ik27804.kd[i].kCode = line[pos+(i*32) : pos+(i*32)+5]
			ik27804.kd[i].kFlag = line[pos+(i*32)+5 : pos+(i*32)+5+1]
			ik27804.kd[i].kValue = line[pos+(i*32)+5+1 : pos+(i*32)+5+1+19]
			ik27804.kd[i].kCmt = line[pos+(i*32)+5+1+19 : pos+(i*32)+5+1+19+7]
			if ik27804.kd[i].kCode == "01935" && ik27804.kd[i].kFlag == "A" {
				B1, _ = strconv.ParseFloat(strings.Replace(ik27804.kd[i].kValue, " ", "", -1), 64)
				//log.Printf("H:%s Code:%s Value:%s Flag:%s Cmt:%s", ik27804.head, ik27804.kd[i].kCode, ik27804.kd[i].kValue, ik27804.kd[i].kFlag, ik27804.kd[i].kCmt)
			} else if ik27804.kd[i].kCode == "01935" && ik27804.kd[i].kFlag == "B" {
				B2, _ = strconv.ParseFloat(strings.Replace(ik27804.kd[i].kValue, " ", "", -1), 64)
				//log.Printf("H:%s Code:%s Value:%s Flag:%s Cmt:%s", ik27804.head, ik27804.kd[i].kCode, ik27804.kd[i].kValue, ik27804.kd[i].kFlag, ik27804.kd[i].kCmt)
			}
		}

		if B1 != 0.0 && B2 == 0.0 {
			if scanner.Scan() == false {
				break
			}
			line = scanner.Text()

			if ik27804.head == line[:78] {
				for i := 0; i < 5; i++ {
					ik27804.kd[i].kCode = line[pos+(i*32) : pos+(i*32)+5]
					ik27804.kd[i].kFlag = line[pos+(i*32)+5 : pos+(i*32)+5+1]
					ik27804.kd[i].kValue = line[pos+(i*32)+5+1 : pos+(i*32)+5+1+19]
					ik27804.kd[i].kCmt = line[pos+(i*32)+5+1+19 : pos+(i*32)+5+1+19+7]
					//log.Printf("*H:%s Code:%s Value:%s Flag:%s Cmt:%s", ik27804.head, ik27804.kd[i].kCode, ik27804.kd[i].kValue, ik27804.kd[i].kFlag, ik27804.kd[i].kCmt)
					if ik27804.kd[i].kCode == "01935" && ik27804.kd[i].kFlag == "B" {
						B2, _ = strconv.ParseFloat(strings.Replace(ik27804.kd[i].kValue, " ", "", -1), 64)
					}
				}
			}
		}

		if B1 != 0.0 {
			// ファイルの書き込み
			B := strconv.FormatFloat(B1+B2, 'f', 1, 64)
			//log.Printf("head:%s B:%s B1:%f B2:%f", ik27804.head, B, B1, B2)
			wline := ik27804.head
			wline = wline + "01935"
			wline = wline + strings.Repeat(" ", 20-len(B)) + B
			wline = wline + strings.Repeat(" ", 153)
			_, err = outfile.WriteString(wline + "\r\n")
			failOnError(err)

		}

		//ik27804の初期化
		ik27804.head = ""
		for i := 0; i < 5; i++ {
			ik27804.kd[i].kCode = ""
			ik27804.kd[i].kValue = ""
			ik27804.kd[i].kFlag = ""
			ik27804.kd[i].kCmt = ""
		}

	}

	log.Print("Finesh !\r\n")

}
