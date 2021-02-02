package helper

import (
	"fmt"
	"os"
)

func AddPagesdb(name string, pages int) []Paper {

	var pagesArr []Paper
	imageBase := getIMAGEDomain()
	pagesArr = make([]Paper, pages+1)

	for i := 0; i <= pages; i++ {
		page := new(Paper)
		page.Number = i + 1
		page.Link = fmt.Sprintf("%s/pages/%v/%v.jpg", imageBase, name, i+1)
		pagesArr[i] = *page
	}

	return pagesArr
}

func getIMAGEDomain() string {
	if root := os.Getenv("ROOT_IMAGES"); root == "" {
		return "http://167.172.41.222:3000"
	} else {
		return root
	}

}
