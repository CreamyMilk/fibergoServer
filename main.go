package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"


	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"jotham/helper"
)

//Uses the ls command to find the number of files
//foo=$(strings < pdffile.pdf | sed -n 's|.*/Count -\{0,1\}\([0-9]\{1,\}\).*|\1|p' | sort -rn | head -n 1)
func directpages(path string) string {
	template := "strings < '%s' | sed -n 's|.*/Count -\\{0,1\\}\\([0-9]\\{1,\\}\\).*|\\1|p' | sort -rn | head -n 1"
	extractscript := fmt.Sprintf(template,path)
	var (
		err    error
		cmdOut []byte
	)
	cmdName := "bash"
	cmdArgs := []string{"-c", extractscript}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return getpdfsls(path)
	}
	pages := string(cmdOut)
	if pages ==""{
		return getpdfpages(path)
	}
	fmt.Printf("The pdf %s has %s done",path,pages)
	return pages
}
func getpdfpages(path string) string {
	extractscript := "exiftool '" + path + "'| grep Page" + "| grep -o '[[:digit:]]*'"
	var (
		err    error
		cmdOut []byte
	)
	cmdName := "bash"
	cmdArgs := []string{"-c", extractscript}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return getpdfsls(path)
	}
	fmt.Fprintln(os.Stderr, "Pages counter exiftool error: ", err)

	pages := string(cmdOut)
	if pages != "" {
		return pages[:len(pages)-1]
	}
	//Count using the ls command
	return getpdfsls(path)
}
func getpdfsls(path string) string {
	//Rethink this implementaion
	extractscript := "ls " + "./uploads/testfolder/'" + path + "'/| wc -l"
	var (
		err    error
		cmdOut []byte
	)
	cmdName := "bash"
	cmdArgs := []string{"-c", extractscript}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Pages counter ls error: ", err)
	}

	pages := string(cmdOut)
	if pages != "" {
		return pages[:len(pages)-1]
	}
	return ""
}
func oldconvert(filename, quality string) (string, string) {
	// f, _ := os.Create(filename)
	basepath := "./uploads/testfolder/" + filename
	//outputlocal := "-sOutputFile=" + basepath + "/cool%03d.jpg"
	outputlocal := "-sOutputFile=" + basepath + "/%d.jpg"
	_ = os.Mkdir(basepath, 0755)
	tempfile := "./uploads/" + filename
	var (
		// cmdOut []byte
		err error
	)
	cmdName := "gs"
	cmdArgs := []string{"-sDEVICE=jpeg", "-dBATCH", "-dNOPAUSE", outputlocal, "-dJPEGQ=100", "-r72x72", "-g595x842", tempfile}
	//the _ corresponds to cmdOut
	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Error in conversion the conversion: ", err)
	}
	pages := getpdfpages("./uploads/" + filename)
	generatedURL := "https://google.com/?pdf_name=" + filename + "&pages=" + pages
	return generatedURL, pages
}
func newconvert(filename, quality string) (string, string, []helper.Paper) {
	//Add quality string
	var imageurls []helper.Paper
	storeagepath := "./uploads/testfolder/" + filename
	outputfile := "-sOutputFile=" + storeagepath + "/%d.jpg"
	_ = os.Mkdir(storeagepath, 0755)
	tempfile := "./uploads/" + filename
	var (
		// cmdOut []byte
		err error
	)
	//gs -dNOPAUSE -sDEVICE=jpeg -r144 -sOutputFile=p%03d.jpg foodmag.pdf
	cmdName := "gs"
	cmdArgs := []string{"-sDEVICE=jpeg","-dNOGC","-dNumRenderingThreads=4","-dBATCH", "-dNOPAUSE", outputfile, "-r144", tempfile}
	//the _ corresponds to cmdOut
	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Error in conversion the conversion: ", err)
	}
	pages := directpages("./uploads/" + filename)
	generatedURL := "https://google.com/?pdf_name=" + filename + "&pages=" + pages

	if i, _ := strconv.Atoi(pages); i != 0 {
		imageurls = helper.AddPagesdb(filename, i)
	}
	return generatedURL, pages, imageurls
}
func main() {
	if err := helper.Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{
		BodyLimit:    52428800, //50mb
		ServerHeader: "Fiber",
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())
	directpages("./uploads/greedy_Alg.pdf")
	//app.Use(middleware.Favicon("./favicon.ico"))
	//Use nginx for server side caching https://ww.nginx.com/resources/wiki/start/topics/examples/reverseproxycachingexample/
	app.Static("/", "./public")
	// app.Static("/pages", "./pages")
	app.Static("/pages", "./uploads/testfolder", fiber.Static{
		Compress:  false,
		ByteRange: false,
		Browse:    true,
	})
	app.Post("/savefiles", func(c *fiber.Ctx) error {
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
			filename, pages, imageurls := newconvert(file.Filename, "200")
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
	})
	// app.Get("/pages/:id",func(c *fiber.Ctx){

	// }
	app.Get("/redirect", func(c *fiber.Ctx) error{
		c.SendString("Hello World")
		//ridirects /rick rolls
		c.Redirect("https://google.com")
		return nil
	})

	app.Get("/sourcecodedownload", func(c *fiber.Ctx) error{
		if err := c.Download("./main.go", "sourcecode"); err != nil {
			c.Next() // Pass err to fiber
		}
		return nil
	})

	app.Use(func(c *fiber.Ctx) error{
		c.SendFile("./public/404.html")
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"error":   "There are no posts!",
		})
	})
	log.Fatal(app.Listen(":3000"))
}
