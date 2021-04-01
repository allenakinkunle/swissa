package main

import (
	"fmt"
	"os"

	"github.com/allenakinkunle/swissa/converter"
)

func main() {
	file, _ := os.Open("/Users/allen/Downloads/biostats.csv")
	defer file.Close()

	converter := converter.NewCSVConverter(file)
	_, err2 := converter.Convert("json", os.Stdout)
	fmt.Println(err2)
	//cmd.Run()

}
