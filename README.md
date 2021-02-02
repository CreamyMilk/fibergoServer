# Golang base API🥛

Make the PDF creation Event Dirven
---
So on file upload
Make a temp file

Then add record to sql table
  * OwnerIDs
  * EditorsIDs
  * OriginalName
  * GeneratedID
  * NoOfPages
  * Regions
  * URL


Publish to the queue a file created event
--- 
  * Turn the pdf into images
  * Upload the original PDF to google cloud storage (_We dont use GCPs Events so as to prevent vendor locking_)
  * Move images to NGINX server
  * Use gofiber as the staic server for content with limited access
 


Reader Updates to be API based
---

#### Reader Page => __domain/?pdf_name="232323"#page/1__
using the *pdf_name*
  * Get Number of pages
  * Get clickable regions
  * Get imagesource urls



### URLs Shortener service
domain/s/i/NewYorkTimes

USE Postgress because
links TBL
  *Long,Short,Facebook,Instagram,Whatsapp

MongoDB device bump 
  *request from which type of device,IP,Location
