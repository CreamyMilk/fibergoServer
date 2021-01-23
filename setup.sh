#!/usr/bin/bash

go get -u github.com/gofiber/fiber/v2
go get -u github.com/gofiber/fiber/v2/middleware/cors
go get -u github.com/gofiber/fiber/v2/middleware/logger
go get -u github.com/gofiber/fiber/v2/middleware/recover

mkdir uploads && cd uploads && mkdir testfolder
sudo apt install ghostscript
sudo apt install exiftool