package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"

	"jotham/helper"
)

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
	generatedURL := "http://167.172.41.222/?pdf_name=" + filename + "&pages=" + pages
	return generatedURL, pages
}
func newconvert(filename, quality string) (string, string, []string) {
	//Add quality string
	var imageurls []string
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
	cmdArgs := []string{"-sDEVICE=jpeg", "-dBATCH", "-dNOPAUSE", outputfile, "-r144", tempfile}
	//the _ corresponds to cmdOut
	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Error in conversion the conversion: ", err)
	}
	pages := getpdfpages("./uploads/" + filename)
	generatedURL := "http://167.172.41.222/?pdf_name=" + filename + "&pages=" + pages

	if i, _ := strconv.Atoi(pages); i != 0 {
		imageurls = helper.AddPagesdb(filename, i)
	}
	return generatedURL, pages, imageurls
}

func main() {
	if err := helper.Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New(&fiber.Settings{
		BodyLimit:    52428800, //50mb
		ServerHeader: "Fiber",
	
	app.Use(middleware.Recover())
	app.Use(cors.New())
	app.Use(middleware.Logger())
	
	//app.Use(middleware.Favicon("./favicon.ico"))
	//Use nginx for server side caching https://ww.nginx.com/resources/wiki/start/topics/examples/reverseproxycachingexample/
	app.Static("/", "./public")
	// app.Static("/pages", "./pages")
	app.Static("/pages", "./uploads/testfolder", fiber.Static{
		Compress:  false,
		ByteRange: false,
		Browse:    true,
	})
	app.Post("/savefiles", func(c *fiber.Ctx) {
		file, err := c.FormFile("myFile")

		if err == nil {
			c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
			filename, pages, imageurls := newconvert(file.Filename, "200")
			c.JSON(&fiber.Map{
				"pages":     pages,
				"filename":  filename,
				"imagesUrl": imageurls,
			})
			//c.Redirect(filename)
		}
	})
	// app.Get("/pages/:id",func(c *fiber.Ctx){

	// }
	app.Get("/redirect", func(c *fiber.Ctx) {
		c.Send("Hello World")
		//ridirects /rick rolls
		c.Redirect("https://google.com")
	})

	app.Get("/sourcecodedownload", func(c *fiber.Ctx) {
		if err := c.Download("./main.go", "sourcecode"); err != nil {
			c.Next(err) // Pass err to fiber
		}
	})

	app.Use(func(c *fiber.Ctx) {
		c.SendFile("./public/404.html")
		c.Status(404).JSON(&fiber.Map{
			"success": false,
			"error":   "There are no posts!",
		})
	})
	app.Listen(3000)
}
