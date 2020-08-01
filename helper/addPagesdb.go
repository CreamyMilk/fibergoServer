package helper

import "fmt"

func AddPagesdb(name string, pages int) []string {
	var pagesArr []string
	for i := 1; i <= pages; i++ {
		pagesArr = append(pagesArr, fmt.Sprintf("http://167.172.41.222:3000/pages/%v/%v.jpg", name, i))
	}
	return pagesArr
}
