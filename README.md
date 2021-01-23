# Golang base API🥛

Make the PDF creation Event Dirven
---
So on file upload
Make a temp file
Add record to sql table
  * OwnerIDs
  * EditorsIDs
  * OriginalName
  * GeneratedID
  * NoOfPages
  * Regions
  * URL

Publish to the queue a new file is available
--- 
  * Turn the pdf into images
  * Upload it original Pdf to google cloud storage (We dont use GCPs E * Events so as to prevent vendor locking
  * Move files to NGINX server so as to serve the image(
  * Use gofiber as the staic server for content with limited access
 


Reader Updates to be API based
---
domain/?pdf_name="232323"#page/1
using the pdf_name
  *Get Number of pages
  *Get clickable regions
  *Get imagesource urls



URLs Shortened verions
readzy.africa/s/i/NewYorkTimes

USE Postgress because
links
  *Long,Short,Facebook,Instagram,Whatsapp
MongoDB device bump 
  *request from which type of device,IP,Location
