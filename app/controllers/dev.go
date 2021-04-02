package controllers

import (

	//"strconv"
	//"time"

	//"github.com/najamsk/eventvisor/eventvisorHQ/repositories"

	"fmt"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/revel/revel"
)

// Dev controller
type Dev struct {
	Admin
}

// Index action: GET
func (c Dev) Index() revel.Result {
	return c.Render()
}

// Index action: GET
func (c Dev) Show(id string) revel.Result {

	f, _ := os.Open(path.Join("uploads/thumbnail", id))

	return c.RenderFile(f, revel.Inline)
}

// Upload action: GET
func (c Dev) Upload() revel.Result {

	return c.Render()
}

//WriteFileToDisk func
func WriteFileToDisk(filename string, file []byte) error {

	targetpath := revel.Config.StringDefault("hq.image.folderbasepath", "")
	targetmode := revel.Config.StringDefault("mode.dev", "mode not found")

	fmt.Printf("targetmode = %v \n", targetmode)
	fmt.Printf("targetpath = %v \n", targetpath)

	dpath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dpath) // for example /home/user

	localfile, err := os.Create(path.Join(targetpath, filename))
	if err != nil {
		fmt.Printf("os.create throw error %v \n", err)
		return err
	}
	fmt.Println("localfile", localfile)
	_, err = localfile.Write(file)
	if err != nil {
		fmt.Printf("localfile.write throw error %v \n", err)

		return err
	}

	return nil
}

//UploadPost  action
func (c Dev) UploadPost(name string, poster []byte) revel.Result {
	fmt.Printf("name is = %v\n", name)
	// fmt.Printf("posterImage is = %v\n", posterImage)
	if c.Params.Files["poster"] == nil {
		fmt.Println("No file was found, please try again")
		// c.Flash.Error("No file was found, please try again")
		// return c.Redirect(Dev.Upload)
	}
	filename := c.Params.Files["poster"][0].Filename
	contenttype := mime.TypeByExtension(path.Ext(filename))
	if contenttype == "" {
		// Try to figure out the content type from the data
		contenttype = http.DetectContentType(poster)
	}

	fmt.Printf("filename= %v \n", filename)
	fmt.Printf("contenttype= %v \n", contenttype)

	// filename, err := SaveFileToDb(c, contenttype, filename)
	// if err != nil {
	// 	c.Flash.Error(err.Error())
	// 	return c.Redirect(App.Index)
	// }
	if err := WriteFileToDisk(filename, poster); err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}
	c.Flash.Success("Success!")
	return c.RenderText(name)
}
