package utils

import (
	"fmt"
	"jotham/helper"

	"os"
	"os/exec"
	"strconv"
)

func Directpages(path string) string {
	template := "strings < '%s' | sed -n 's|.*/Count -\\{0,1\\}\\([0-9]\\{1,\\}\\).*|\\1|p' | sort -rn | head -n 1"
	extractscript := fmt.Sprintf(template, path)
	var (
		err    error
		cmdOut []byte
	)
	cmdName := "bash"
	cmdArgs := []string{"-c", extractscript}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return Getpdfsls(path)
	}
	pages := string(cmdOut)
	if pages == "" {
		return Getpdfpages(path)
	}
	fmt.Printf("The pdf %s has %s done", path, pages)
	return pages
}
func Getpdfpages(path string) string {
	extractscript := "exiftool '" + path + "'| grep Page" + "| grep -o '[[:digit:]]*'"
	var (
		err    error
		cmdOut []byte
	)
	cmdName := "bash"
	cmdArgs := []string{"-c", extractscript}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return Getpdfsls(path)
	}
	fmt.Fprintln(os.Stderr, "Pages counter exiftool error: ", err)

	pages := string(cmdOut)
	if pages != "" {
		return pages[:len(pages)-1]
	}
	//Count using the ls command
	return Getpdfsls(path)
}
func Getpdfsls(path string) string {
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
	pages := Getpdfpages("./uploads/" + filename)
	generatedURL := "http://readzy.africa/?pdf_name=" + filename + "&pages=" + pages
	return generatedURL, pages
}
func Newconvert(filename, quality string) (string, string, []helper.Paper) {
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
	cmdArgs := []string{"-sDEVICE=jpeg", "-dNOGC", "-dNumRenderingThreads=4", "-dBATCH", "-dNOPAUSE", outputfile, "-r144", tempfile}
	//the _ corresponds to cmdOut
	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Error in conversion the conversion: ", err)
	}
	pages := Directpages("./uploads/" + filename)
	generatedURL := "http://readzy.africa/?pdf_name=" + filename + "&pages=" + pages

	if i, _ := strconv.Atoi(pages); i != 0 {
		imageurls = helper.AddPagesdb(filename, i)
	}
	return generatedURL, pages, imageurls
}
