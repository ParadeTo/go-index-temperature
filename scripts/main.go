package main

import (
	. "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-pe-pb/model"
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func insert(conn *DB, pbdata, pedata map[string]model.Asset) map[string]model.Asset {
	for k, v := range pedata {
		if _pbdata, ok := pbdata[k]; ok {
			v.Pb = _pbdata.Pb
		}
		conn.Create(v)
	}

	return pedata
}

func readData(filename string, getAsset func([]string) model.Asset) map[string]model.Asset {
	m := map[string]model.Asset{}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	b := bufio.NewReader(f)
	line, _, err := b.ReadLine()
	for line, _, err = b.ReadLine(); err == nil; line, _, err = b.ReadLine() {
		tokens := strings.Split(string(line), ",")
		asset := getAsset(tokens)
		m[tokens[0]] = asset
	}

	return m
}

func parseFloat(s string, size int) float64 {
	f, err := strconv.ParseFloat(s, size)
	if err != nil {
		f = -1
	}
	return f
}

func main() {
	conn, err := Open("mysql", "root:123456@tcp(localhost:3306)/pepb?charset=utf8")
	//conn.LogMode(true)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.DB().SetMaxIdleConns(10)
	conn.DB().SetMaxOpenConns(100)

	if !conn.HasTable(&model.Asset{}) {
		conn.Debug().AutoMigrate(&model.Asset{})
	}

	code := "000015.sh"
	pbfile := "./data/000015.sh/pb.csv"
	pbdata := readData(pbfile, func(i []string) model.Asset {
		return model.Asset{
			Code: code,
			Date: i[0],
			Price: float32(parseFloat(i[1], 32)),
			Cap: parseFloat(i[2], 64),
			Pb: float32(parseFloat(i[3], 32)),
		}
	})

	fefile := "./data/000015.sh/pe.csv"
	pedata := readData(fefile, func(i []string) model.Asset {
		return model.Asset{
			Code: code,
			Date: i[0],
			Price: float32(parseFloat(i[1], 32)),
			Cap: parseFloat(i[2], 64),
			Pe: float32(parseFloat(i[3], 32)),
		}
	})

	insert(conn, pbdata, pedata)

	return
}
