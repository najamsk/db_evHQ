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
	//"github.com/najamsk/eventvisor/eventvisorHQ/config"
)

// Conference controller
type Sessions struct {
	Admin
}

//var HqConfig *config.Config
// Index action: GET
func (c Sessions) Index() revel.Result {
	return c.Render()
}

// List action: GET
func (c Sessions) List(id string) revel.Result {
	fmt.Printf("conferece id is = %v\n", id)
	var srv = services.SessionService{}
	vm, srvError := srv.GetSessionsByConferenceID(id)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	fmt.Printf("sessionhye", vm)
	// return c.RenderText(sessionsListVM.ClientName)
	return c.Render(vm)
}

// Details action: GET
func (c Sessions) Details(id string) revel.Result {
	fmt.Printf("conferece id is = %v\n", id)
	var srv = services.SessionService{}
	vm, srvError := srv.GetCByID(id)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	fmt.Println("session Data", vm)

	// return c.RenderText(sessionsListVM.ClientName)
	// c.ViewArgs["vm"] = vm
	// return c.RenderTemplate("/views/Sessions/Edit.html")
	fmt.Printf("confernce return from service is %v \n ", vm)
	return c.Render(vm)
}

// List action: GET
func (c Sessions) ListByClientID(id string) revel.Result {

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
func (c Sessions) SessionsList() revel.Result {
	return c.Render()
}

// Create action: GET
func (c Sessions) Create(id string) revel.Result {

	fmt.Printf("conferece id is = %v\n", id)
	var srv = services.SessionService{}
	vm, srvError := srv.GetConferenceAndClient(id)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	return c.Render(vm)
}

// Create action: GET
func (c Sessions) CreatePost(id string,
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
	venue string,
	weight string,
	IsFeatured bool,
	radius float64) revel.Result {

	sortOrder, _ := strconv.Atoi(weight)

	// id passed in maybe conference id which is not getting used,
	// so i dont think session will be binded with conference on creation. its a bug.
	fmt.Println("create session  post func activated")
	// fmt.Printf("form id = %v \n", id)
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

	clientid, _ := uuid.FromString(clientID)
	confid, _ := uuid.FromString(id)

	i, strerr := strconv.ParseInt(startDate, 10, 64)
	if strerr != nil {
		fmt.Println("stat date type show err", strerr)
		panic(strerr) //maybe shourld return flash error to form page
	}

	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, Enderr := strconv.ParseInt(endDate, 10, 64)
	if Enderr != nil {
		fmt.Println("stat date type show err", strerr)
		panic(Enderr) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()

	fmt.Printf("form startdate in go = %v \n", startDateStamp)
	fmt.Printf("form enddate in go = %v \n", endDateStamp)

	inputData := viewmodels.SessionCreateVMWrite{
		Title:           name,
		ClientID:        clientid,
		ConfID:          confid,
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
		Venue:           venue,
		SortOrder:       sortOrder,
		IsFeatured:      IsFeatured,
	}
	//fmt.Println(inputData)
	srv := services.SessionService{}
	sessionID, srvError := srv.CreateSession(inputData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
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
			return c.Redirect("/admin/sessions/" + sessionID.String())
		}
		posterData := viewmodels.ImageVMWrite{
			Name:           sessionID.String() + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.sessionposter.folderbasepath", ""),
			ImageURLPrefix: "session/poster",
			EntityID:       sessionID,
			EntityType:     "session",
			ImageCategory:  "poster",
			IsActive:       true,
		}

		PostrErr := imgSrv.CreateImage(posterData)

		if PostrErr == nil {
			imgsaveerr := imgSrv.WriteFileToDisk(posterData.FolderPath, posterData.Name, poster)
			if imgsaveerr != nil {
				fmt.Println("Poster size must be less than 2mb")
				c.Flash.Error("Poster size must be less than 2mb")
				Rmverr := os.Remove(posterData.FolderPath + "/" + posterData.Name)
				if Rmverr != nil {
					fmt.Println("file not deleted", Rmverr)

				}
				return c.Redirect("/admin/sessions/" + sessionID.String())
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
			return c.Redirect("/admin/sessions/" + sessionID.String())
		}

		thumbnailData := viewmodels.ImageVMWrite{
			Name:           sessionID.String() + filepath.Ext(Thumbnailfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.sessionthumbnail.folderbasepath", ""),
			ImageURLPrefix: "session/thumbnail",
			EntityID:       sessionID,
			EntityType:     "session",
			ImageCategory:  "thumbnail",
			IsActive:       true,
		}

		ThumbNLErr := imgSrv.CreateImage(thumbnailData)

		if ThumbNLErr == nil {
			thumbnailErr := imgSrv.WriteFileToDisk(thumbnailData.FolderPath, thumbnailData.Name, thumbnail)
			if thumbnailErr != nil {
				fmt.Println("Thumbnail size must be less than 2mb")
				c.Flash.Error("Thumbnail size must be less than 2mb")
				thmberr := os.Remove(thumbnailData.FolderPath + "/" + thumbnailData.Name)
				if thmberr != nil {
					fmt.Println("thumbnail file not deleted", thmberr)

				}
				return c.Redirect("/admin/sessions/" + sessionID.String())
			}
		}
	}
	c.Flash.Success("Session Successfully Created")

	return c.Redirect("/admin/sessions/" + sessionID.String())
	//return nil
}

func (c Sessions) EditPost(id string,
	name string,
	summary string,
	durationDisplay string,
	details string,
	clientID string,
	confID string,
	poster []byte,
	thumbnail []byte,
	isActive bool,
	venue string,
	weight string,
	IsFeatured bool,
	startDate string,
	endDate string,
	address string,
	latitude float64,
	longitude float64,
	radius float64) revel.Result {

	sortOrder, _ := strconv.Atoi(weight)

	fmt.Println("create session  post func activated")
	// fmt.Printf("form id = %v \n", id)
	fmt.Printf("form clientid = %v \n", clientID)
	fmt.Printf("form confid = %v \n", confID)
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
	fmt.Printf("form venue = %v \n", sortOrder)

	var srv = services.SessionService{}

	clientid, _ := uuid.FromString(clientID)
	confid, _ := uuid.FromString(confID)
	sessionid, _ := uuid.FromString(id)

	i, strerr := strconv.ParseInt(startDate, 10, 64)
	if strerr != nil {
		fmt.Println("strt date type shows err at edit session,err")
		panic(strerr) //maybe shourld return flash error to form page
	}

	startDateStamp := time.Unix(i/1000, 0).UTC()

	i, enderr := strconv.ParseInt(endDate, 10, 64)
	if enderr != nil {
		fmt.Println("Date type shows err at edit session,err", enderr)
		panic(enderr) //maybe shourld return flash error to form page
	}
	endDateStamp := time.Unix(i/1000, 0).UTC()

	fmt.Printf("form startdate in go = %v \n", startDateStamp)
	fmt.Printf("form enddate in go = %v \n", endDateStamp)

	inputData := viewmodels.SessionCreateVMWrite{
		Title:           name,
		ClientID:        clientid,
		ConfID:          confid,
		StartDate:       startDateStamp,
		EndDate:         endDateStamp,
		DurationDisplay: durationDisplay,
		Details:         details,
		Summary:         summary,
		IsActive:        isActive,
		Address:         address,
		Latitude:        latitude,
		Longitude:       longitude,
		Venue:           venue,
		LocationRadius:  radius,
		SortOrder:       sortOrder,
		IsFeatured:      IsFeatured,
	}

	inputData.ID = sessionid

	sessionID, srvError := srv.UpdateSession(inputData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	imgSrv := services.ImageService{}
	//poster and thumbnail updation
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
			return c.Redirect("/admin/sessions/" + sessionID.String())
		}

		posterData := viewmodels.ImageVMWrite{
			Name:           sessionID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.sessionposter.folderbasepath", ""),
			ImageURLPrefix: "session/poster",
			EntityID:       sessionID,
			EntityType:     "session",
			ImageCategory:  "poster",
			IsActive:       true,
		}

		PostrErr := imgSrv.Update(posterData)

		if PostrErr == nil {
			imgsaveerr := imgSrv.WriteFileToDisk(posterData.FolderPath, posterData.Name, poster)
			if imgsaveerr != nil {
				c.Flash.Error("Poster size must be less than 2mb")
				fmt.Println("Poster size must be less than 2mb")
				Rmverr := os.Remove(posterData.FolderPath + "/" + posterData.Name)
				if Rmverr != nil {
					fmt.Println("file not deleted", Rmverr)

				}
				return c.Redirect("/admin/sessions/" + sessionID.String())
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
			return c.Redirect("/admin/sessions/" + sessionID.String())
		}
		thumbnailData := viewmodels.ImageVMWrite{
			Name:           sessionID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(Thumbnailfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.sessionthumbnail.folderbasepath", ""),
			ImageURLPrefix: "session/thumbnail",
			EntityID:       sessionID,
			EntityType:     "session",
			ImageCategory:  "thumbnail",
			IsActive:       true,
		}

		ThumbNLErr := imgSrv.Update(thumbnailData)

		if ThumbNLErr == nil {
			thumbnailErr := imgSrv.WriteFileToDisk(thumbnailData.FolderPath, thumbnailData.Name, thumbnail)
			if thumbnailErr != nil {
				c.Flash.Error("Thumbnail size must be less than 2mb")
				println("Thumbnail size must be less than 2mb")
				errThb := os.Remove(thumbnailData.FolderPath + "/" + thumbnailData.Name)
				if errThb != nil {
					fmt.Println("thumbnail file not deleted", errThb)

				}
				return c.Redirect("/admin/sessions/" + sessionID.String())
			}
		}
	}

	// fmt.Printf("vm conf name is = %v \n", vm.ConfName)
	c.Flash.Success("Session Successfully Updated")
	return c.Redirect("/admin/sessions/" + sessionID.String())
}
func (c Sessions) ListBySpeaker(confid string, id string) revel.Result {
	fmt.Printf("conferece id is = %v\n", id)
	var srv = services.SessionService{}
	vm, srvError := srv.GetSessionsBySpeakerID(id, confid)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	fmt.Printf("sessionhye", vm)
	// return c.RenderText(sessionsListVM.ClientName)
	return c.Render(vm)

}
