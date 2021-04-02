package controllers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Conference controller
type Conferences struct {
	Admin
}

// Index action: GET
func (c Conferences) Index() revel.Result {

	return c.Render()
}

// List action: GET
func (c Conferences) List() revel.Result {
	return c.Render()
}

// List action: GET
func (c Conferences) ListByClientID(id string) revel.Result {

	var srv = services.ConferenceService{}
	conferencesVM, srvError := srv.GetConferencesByClientID(id)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}

	// return c.RenderText(id)
	return c.Render(conferencesVM)
}

// SessionsList action: GET
func (c Conferences) SessionsList() revel.Result {
	return c.Render()
}

// TODO: needs to redirect to listing or thank you page. with flash ? also use reverse routes inside html templates

// Create action: Create
func (c Conferences) Create(id string) revel.Result {
	//its working using c.Params.Get for post values as well
	// name := c.Params.Get("name")
	var srv = services.ConferenceService{}
	// confDB, srvError := srv.GetConferencesByClientID(id)
	// if srvError != nil {
	// 	fmt.Printf("conf service returns error %v \n", srvError)
	// }

	vm := viewmodels.ConferenceEditVMRead{}
	clientid, _ := uuid.FromString(id)
	vm.ClientID = clientid
	vm.ClientName = srv.GetClientNameByClientID(clientid)
	return c.Render(vm)
}

// Create action: Edit
func (c Conferences) Edit(id string) revel.Result {

	//its working using c.Params.Get for post values as well
	// name := c.Params.Get("name")
	var srv = services.ConferenceService{}
	conference, srvError := srv.GetConferenceByID(id)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
	}
	fmt.Printf("confernce return from service is %v \n ", conference)

	return c.Render(conference)
	// return c.RenderText("mil gia")
}

func (c Conferences) EditPost(id string,
	name string,
	summary string,
	durationDisplay string,
	details string,
	clientID string,
	poster []byte,
	thumbnail []byte,
	isActive bool,
	startDate string,
	endDate string,
	address string,
	latitude float64,
	longitude float64,
	radius float64) revel.Result {

	fmt.Println("edit post func activated")
	fmt.Printf("form id = %v \n", id)
	fmt.Printf("form clientid = %v \n", clientID)
	fmt.Printf("form name = %v \n", name)
	fmt.Printf("form summary = %v \n", summary)
	fmt.Printf("form durationDisplay = %v \n", durationDisplay)
	fmt.Printf("form details = %v \n", details)
	fmt.Printf("form isActive = %v \n", isActive)
	fmt.Printf("form startDate = %v \n", startDate)
	fmt.Printf("form endDate = %v \n", endDate)
	fmt.Printf("form address = %v \n", address)
	fmt.Printf("form latitude = %v \n", latitude)
	fmt.Printf("form longitude = %v \n", longitude)
	fmt.Printf("form longitude = %v \n", longitude)
	if radius < 2 {
		c.Flash.Error("Radius must b greater than 2")
		fmt.Println("Radius must b greater than 2")
		return c.Redirect(" /admin/conferences/details/" + id)
	}

	i, strerr := strconv.ParseInt(startDate, 10, 64)
	if strerr != nil {
		fmt.Println("start date err")
		panic(strerr) //maybe shourld return flash error to form page
	}

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, enderr := strconv.ParseInt(endDate, 10, 64)
	if enderr != nil {
		fmt.Println("end date err")
		panic(enderr) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()

	fmt.Printf("form startdate in go = %v \n", startDateStamp)
	fmt.Printf("form enddate in go = %v \n", endDateStamp)

	var srv = services.ConferenceService{}

	clientid, _ := uuid.FromString(clientID)

	confData := viewmodels.ConferenceEditVMWrite{ID: id,
		Title:           name,
		ClientID:        clientid,
		StartDate:       startDateStamp,
		EndDate:         endDateStamp,
		DurationDisplay: durationDisplay,
		Details:         details,
		Summary:         summary,
		IsActive:        isActive,
		Address:         address,
		Latitude:        latitude,
		Longitude:       longitude,
		LocationRadius:  radius,
	}

	confID, srvError := srv.UpdateConference(confData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}

	//poster and thumbnail updation
	imgSrv := services.ImageService{}
	if len(c.Params.Files["poster"]) > 0 {
		posterfile := c.Params.Files["poster"][0].Filename
		fmt.Println("poster file", posterfile)
		// getting poster Extention
		postertype := mime.TypeByExtension(path.Ext(posterfile))
		if postertype == "" {
			// Try to figure out the content type from the data
			postertype = http.DetectContentType(poster)
		}
		fmt.Println("postertype", postertype)

		// checking valid poster Extention
		ValidPoster := imgSrv.Validation(postertype)
		if ValidPoster == false {
			c.Flash.Error("Poster type is not valid")
			fmt.Println("Poster type is not valid")
			return c.Redirect("/admin/conferences/details/" + confID.String())
		}
		posterData := viewmodels.ImageVMWrite{
			Name:           confID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.conferenceposter.folderbasepath", ""),
			ImageURLPrefix: "conference/poster",
			EntityID:       confID,
			EntityType:     "conference",
			ImageCategory:  "poster",
			IsActive:       true,
		}

		PostrErr := imgSrv.Update(posterData)

		if PostrErr == nil {
			imgsaveerr := imgSrv.WriteFileToDisk(posterData.FolderPath, posterData.Name, poster)
			if imgsaveerr != nil {
				c.Flash.Error("Poster size must be less than 2mb")
				errRmv := os.Remove(posterData.FolderPath + "/" + posterData.Name)
				if errRmv != nil {
					fmt.Println("file not deleted", errRmv)

				}
				fmt.Println("write do disk show err", errRmv)
				return c.Redirect("/admin/conferences/details/" + confID.String())
			}
		}
	}

	// for thumbnail
	if len(c.Params.Files["thumbnail"]) > 0 {
		Thumbnailfile := c.Params.Files["thumbnail"][0].Filename

		// getting thumbnail Extention
		thumbnailtype := mime.TypeByExtension(path.Ext(Thumbnailfile))
		if thumbnailtype == "" {
			// Try to figure out the content type from the data
			thumbnailtype = http.DetectContentType(thumbnail)
		}

		fmt.Println("thumbnailtype1", thumbnailtype)
		// checking valid thumbnail Extention
		ValidThumbnail := imgSrv.Validation(thumbnailtype)
		if ValidThumbnail == false {
			c.Flash.Error("Thumbnail type is not valid")
			fmt.Println("Thumbnail type is not valid")
			return c.Redirect("/admin/conferences/details/" + confID.String())
		}

		thumbnailData := viewmodels.ImageVMWrite{
			Name:           confID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(Thumbnailfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.conferencethumbnail.folderbasepath", ""),
			ImageURLPrefix: "conference/thumbnail",
			EntityID:       confID,
			EntityType:     "conference",
			ImageCategory:  "thumbnail",
			IsActive:       true,
		}

		ThumbNLErr := imgSrv.Update(thumbnailData)

		if ThumbNLErr == nil {
			thumbnailErr := imgSrv.WriteFileToDisk(thumbnailData.FolderPath, thumbnailData.Name, thumbnail)
			if thumbnailErr != nil {
				fmt.Println("Thumbnail size must be less than 2mb")
				c.Flash.Error("Thumbnail size must be less than 2mb")
				errthmb := os.Remove(thumbnailData.FolderPath + "/" + thumbnailData.Name)
				if errthmb != nil {
					fmt.Println("thumbnail file not deleted", errthmb)

				}
				return c.Redirect("/admin/conferences/details/" + confID.String())
			}
		}
	}

	// return c.RenderText("will post data")
	c.Flash.Success("Conference Successfully Updated")

	return c.Redirect(" /admin/conferences/details/" + confID.String())

}

// TODO: needs to redirect to listing or thank you page. with flash ? also use reverse routes inside html templates
// CreatePost action: Create
func (c Conferences) CreatePost(
	id string,
	name string,
	summary string,
	durationDisplay string,
	details string,
	poster []byte,
	thumbnail []byte,
	isActive bool,
	startDate string,
	endDate string,
	address string,
	latitude float64,
	longitude float64,
	radius float64) revel.Result {
	// c.Flash.Success(greet)
	fmt.Println("edit post func activated")
	// fmt.Printf("form id = %v \n", id)
	fmt.Printf("form name = %v \n", name)
	fmt.Printf("form summary = %v \n", summary)
	fmt.Printf("form durationDisplay = %v \n", durationDisplay)
	fmt.Printf("form details = %v \n", details)
	fmt.Printf("form isActive = %v \n", isActive)
	fmt.Printf("form startDate = %v \n", startDate)
	fmt.Printf("form endDate = %v \n", endDate)
	fmt.Printf("form address = %v \n", address)
	fmt.Printf("form latitude = %v \n", latitude)
	fmt.Printf("form longitude = %v \n", longitude)
	fmt.Printf("form longitude = %v \n", longitude)
	clientid, _ := uuid.FromString(id)
	i, errEnd := strconv.ParseInt(startDate, 10, 64)
	if errEnd != nil {
		fmt.Println("date type shows error")
		panic(errEnd) //maybe shourld return flash error to form page
	}

	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, strerr := strconv.ParseInt(endDate, 10, 64)
	if strerr != nil {
		fmt.Println("date type shows error")
		panic(strerr) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()

	fmt.Printf("form startdate in go = %v \n", startDateStamp)
	fmt.Printf("form enddate in go = %v \n", endDateStamp)

	var srv = services.ConferenceService{}

	confData := viewmodels.ConferenceEditVMWrite{
		Title:           name,
		ClientID:        clientid,
		StartDate:       startDateStamp,
		EndDate:         endDateStamp,
		DurationDisplay: durationDisplay,
		Details:         details,
		Summary:         summary,
		IsActive:        isActive,
		Address:         address,
		Latitude:        latitude,
		Longitude:       longitude,
		LocationRadius:  radius,
	}

	confID, srvError := srv.CreateConference(confData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}
	imgSrv := services.ImageService{}
	if len(c.Params.Files["poster"]) > 0 {
		posterfile := c.Params.Files["poster"][0].Filename
		// getting poster Extention
		postertype := mime.TypeByExtension(path.Ext(posterfile))
		if postertype == "" {
			// Try to figure out the content type from the data
			postertype = http.DetectContentType(poster)
		}
		fmt.Println("postertype", postertype)

		// checking valid poster Extention
		ValidPoster := imgSrv.Validation(postertype)
		if ValidPoster == false {
			fmt.Println("Poster type is not valid")
			c.Flash.Error("Poster type is not valid")
			return c.Redirect("/admin/conferences/details/" + confID.String())
		}
		posterData := viewmodels.ImageVMWrite{
			Name:           confID.String() + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.conferenceposter.folderbasepath", ""),
			ImageURLPrefix: "conference/poster",
			EntityID:       confID,
			EntityType:     "conference",
			ImageCategory:  "poster",
			IsActive:       true,
		}

		PostrErr := imgSrv.CreateImage(posterData)

		if PostrErr == nil {
			imgsaveerr := imgSrv.WriteFileToDisk(posterData.FolderPath, posterData.Name, poster)
			if imgsaveerr != nil {
				fmt.Println("Poster size must be less than 2mb", imgsaveerr)
				c.Flash.Error("Poster size must be less than 2mb")
				errRmv := os.Remove(posterData.FolderPath + "/" + posterData.Name)
				if errRmv != nil {
					fmt.Println("file not deleted", errRmv)

				}
				return c.Redirect("/admin/conferences/details/" + confID.String())
			}
		}
	}

	// for thumbnail
	if len(c.Params.Files["thumbnail"]) > 0 {
		Thumbnailfile := c.Params.Files["thumbnail"][0].Filename

		// getting thumbnail Extention
		thumbnailtype := mime.TypeByExtension(path.Ext(Thumbnailfile))
		if thumbnailtype == "" {
			// Try to figure out the content type from the data
			thumbnailtype = http.DetectContentType(thumbnail)
		}

		fmt.Println("thumbnailtype1", thumbnailtype)
		// checking valid thumbnail Extention
		ValidThumbnail := imgSrv.Validation(thumbnailtype)
		if ValidThumbnail == false {
			fmt.Println("Thumbnail type is not valid")
			c.Flash.Error("Thumbnail type is not valid")
			return c.Redirect("/admin/conferences/details/" + confID.String())
		}
		thumbnailData := viewmodels.ImageVMWrite{
			Name:           confID.String() + filepath.Ext(Thumbnailfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.conferencethumbnail.folderbasepath", ""),
			ImageURLPrefix: "conference/thumbnail",
			EntityID:       confID,
			EntityType:     "conference",
			ImageCategory:  "thumbnail",
			IsActive:       true,
		}

		ThumbNLErr := imgSrv.CreateImage(thumbnailData)

		if ThumbNLErr == nil {
			thumbnailErr := imgSrv.WriteFileToDisk(thumbnailData.FolderPath, thumbnailData.Name, thumbnail)
			if thumbnailErr != nil {
				c.Flash.Error("Thumbnail size must be less than 2mb")
				fmt.Println("Thumbnail size must be less than 2mb")
				errtmb := os.Remove(thumbnailData.FolderPath + "/" + thumbnailData.Name)
				if errtmb != nil {
					fmt.Println("thumbnail file not deleted", errtmb)

				}
				return c.Redirect("/admin/conferences/details/" + confID.String())
			}
		}
	}
	c.Flash.Success("Conference Successfully Created")
	return c.Redirect(" /admin/conferences/details/" + confID.String())
}

// Payments Index action: Payments
func (c Conferences) Payments() revel.Result {
	return c.Render()
}
