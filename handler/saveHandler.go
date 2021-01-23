package handler

import (
	"fmt"
	"jotham/helper"
	"jotham/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SaveHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("myFile")

	if err == nil {
		err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
		if err != nil {
			fmt.Println(err)
			c.Status(500)
			//panic(err)
			fmt.Printf("%T", err)
		}
		fmt.Println("Doing the db stuff")
		filename, pages, imageurls := utils.Newconvert(file.Filename, "200")
		collection := helper.Mg.Db.Collection("pdfs")
		page := new(helper.Record)
		page.ID = ""
		page.PdfName = filename
		page.NPages, _ = strconv.Atoi(pages)
		page.Pages = imageurls

		insertionResult, err := collection.InsertOne(c.Context(), page)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T %v", insertionResult, insertionResult)
		fmt.Printf("%v", imageurls)
		c.JSON(&fiber.Map{
			"pages":     pages,
			"filename":  filename,
			"imagesUrl": imageurls,
		})
		c.Redirect(filename)
	}
	return err
}
