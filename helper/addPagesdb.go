package helper

import "fmt"

func AddPagesdb(name string, pages int) []Paper {

	var pagesArr []Paper
	pagesArr = make([]Paper, pages+1)

	for i := 0; i <= pages; i++ {
		page := new(Paper)
		page.Number = i + 1
		page.Link = fmt.Sprintf("http://167.172.41.222:3000/pages/%v/%v.jpg", name, i+1)
		pagesArr[i] = *page
	}

	return pagesArr
}
