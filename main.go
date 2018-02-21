package main

import (
	"flag"
	"strings"
	"fmt"
	"path/filepath"
	"os"
	"log"
	"sort"
	"strconv"
)

const (
	kb = 1024
	mb = 1024 * 1024
	gb = 1024 * 1024 * 1024
)

type extSta struct {
	Ext string
	Qtd int
	Tam int64
}

var (
	strignoradas string
	ignoradas    []string
	exts         = make(map[string]extSta)
	tot          = 0
	totTam       int64
)

func main() {
	flag.StringVar(&strignoradas, "ig", "", "lista de pastas a ignorar separadas por vírgula")
	flag.Parse()
	ignoradas = strings.Split(strignoradas, ",")
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		for _, ig := range ignoradas {
			if info.Name() == ig && info.IsDir() {
				return filepath.SkipDir
			}
		}
		if info.Name() == "." {
			return nil
		}
		ext := filepath.Ext(info.Name())
		es := exts[ext]
		es.Ext = ext
		es.Qtd += 1
		es.Tam += info.Size()
		exts[ext] = es
		tot += 1
		totTam += info.Size()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	extStas := make([]extSta, 0)
	for _, es := range exts {
		extStas = append(extStas, es)
	}
	sort.Slice(extStas, func(i, j int) bool {
		return extStas[i].Qtd > extStas[j].Qtd
	})
	fmt.Println("total quantidade:", tot)
	for _, es := range extStas {
		fmt.Printf("%v: %v\n", es.Ext, es.Qtd)
	}
	fmt.Println("---")
	fmt.Println("total tamanho:", tam(totTam))
	sort.Slice(extStas, func(i, j int) bool {
		return extStas[i].Tam > extStas[j].Tam
	})
	for _, es := range extStas {
		fmt.Printf("%v: %v\n", es.Ext, tam(es.Tam))
	}
}

func tam(size int64) string {
	if size < kb {
		return fmt.Sprintf("%v B", size)
	}
	if size < mb {
		return fmt.Sprintf("%v KB", strconv.FormatFloat(float64(size)/kb, 'f', 2, 32))
	}
	if size < gb {
		return fmt.Sprintf("%v MB", strconv.FormatFloat(float64(size)/mb, 'f', 2, 32))
	}
	return fmt.Sprintf("%vGB", strconv.FormatFloat(float64(size)/gb, 'f', 2, 32))
}
