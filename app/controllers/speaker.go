package controllers

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

// Conference controller
type Speakers struct {
	Admin
}

// Index action: GET
func (c Speakers) Index() revel.Result {
	return c.Render()
}

// List action: GET
func (c Speakers) ListByConference(id string) revel.Result {

	var serv = services.SpeakerService{}
	speakerVM, servError := serv.GetSpeakerByConferenceID(id)
	if servError != nil {
		fmt.Printf("conf service returns error %v \n", servError)
		return c.RenderError(servError)
	}
	fmt.Printf("hye121", speakerVM)
	// return c.RenderText(id)
	return c.Render(speakerVM)
}

func (c Speakers) ListBySession(confid string, id string) revel.Result {

	var serv = services.SpeakerService{}

	speakerVM, servError := serv.GetSpeakerBySession(id, confid)
	if servError != nil {
		fmt.Printf("conf service returns error %v \n", servError)
		return c.RenderError(servError)
	}
	fmt.Printf("hye121", speakerVM)
	// return c.RenderText(id)
	return c.Render(speakerVM)
}
func (c Speakers) Edit(id string, sessionId string) revel.Result {

	//its working using c.Params.Get for post values as well
	// name := c.Params.Get("name")
	var srv = services.SpeakerService{}
	speaker, srvError := srv.GetSpeakerByID(id, sessionId)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
	}
	fmt.Printf("confernce return from service is %v \n ", speaker)
	return c.Render(speaker)
	// return c.RenderText("mil gia")
}

func (c Speakers) EditPost(
	id string,
	confiD string,
	firstname string,
	lastname string,
	email string,
	phonenumber string,
	organization string,
	designation string,
	bio string,
	clientid string,
	profileimg string,
	sessionWeight string,
	conferenceWeight string,
	sessionid string,
	poster []byte,
	profile []byte,
	facebook string,
	twitter string,
	youtube string,
	linkedin string,

) revel.Result {
	speakerID, _ := uuid.FromString(id)

	Session_sortOrder, _ := strconv.Atoi(sessionWeight)
	conference_sortOrder, _ := strconv.Atoi(conferenceWeight)

	fmt.Println("edit post func activated")
	fmt.Printf("form id = %v \n", id)
	fmt.Printf("form name = %v \n", firstname)
	fmt.Printf("form lastname = %v \n", profileimg)
	fmt.Printf("form email = %v \n", email)
	fmt.Printf("form phonenumber = %v \n", phonenumber)
	fmt.Printf("form organization = %v \n", organization)
	fmt.Printf("form designation = %v \n", designation)
	fmt.Printf("form sessionid = %v \n", sessionid)
	clientID, _ := uuid.FromString(clientid)
	confID, _ := uuid.FromString(confiD)
	sessionID, _ := uuid.FromString(sessionid)
	speakerData := viewmodels.SpeakerEditVMWrite{
		FirstName:        firstname,
		LastName:         lastname,
		Email:            email,
		Organization:     organization,
		Designation:      designation,
		PhoneNumber:      phonenumber,
		Bio:              bio,
		ClientID:         clientID,
		SessionWeight:    Session_sortOrder,
		ConferenceWeight: conference_sortOrder,
		SessionID:        sessionID,
		ConferenceID:     confID,
		Facebook:         facebook,
		Twitter:          twitter,
		Linkedin:         linkedin,
		Youtube:          youtube,
	}
	speakerData.ID = speakerID
	srv := services.SpeakerService{}
	speakerdata, srvError := srv.UpdateSpeaker(speakerData)
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
			return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + speakerID.String())
		}

		posterData := viewmodels.ImageVMWrite{
			Name:           speakerID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.userposter.folderbasepath", ""),
			ImageURLPrefix: "user/poster",
			EntityID:       speakerID,
			EntityType:     "user",
			ImageCategory:  "poster",
			IsActive:       true,
		}

		PostrErr := imgSrv.Update(posterData)

		if PostrErr == nil {
			imgsaveerr := imgSrv.WriteFileToDisk(posterData.FolderPath, posterData.Name, poster)
			if imgsaveerr != nil {
				fmt.Println("Poster size must be less than 2mb", imgsaveerr)
				c.Flash.Error("Poster size must be less than 2mb")
				errRmv := os.Remove(posterData.FolderPath + "/" + posterData.Name)
				if errRmv != nil {
					fmt.Println("file not deleted", errRmv)

				}
				return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + speakerID.String())
			}
		}
	}
	if len(c.Params.Files["profile"]) > 0 {
		// for profile
		Profilefile := c.Params.Files["profile"][0].Filename

		// getting thumbnail Extention
		profiletype := mime.TypeByExtension(path.Ext(Profilefile))
		if profiletype == "" {
			// Try to figure out the content type from the data
			profiletype = http.DetectContentType(profile)
		}

		fmt.Println("thumbnailtype1", profiletype)
		// checking valid thumbnail Extention
		ValidProfile := imgSrv.Validation(profiletype)
		if ValidProfile == false {
			fmt.Println("profile type is not valid")
			c.Flash.Error("profile type is not valid")
			return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + speakerID.String())
		}
		ProfileData := viewmodels.ImageVMWrite{
			Name:           speakerID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(Profilefile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.userprofile.folderbasepath", ""),
			ImageURLPrefix: "user/thumbnail",
			EntityID:       speakerID,
			EntityType:     "user",
			ImageCategory:  "user_profile",
			IsActive:       true,
		}

		profileErr := imgSrv.Update(ProfileData)

		if profileErr == nil {
			profErr := imgSrv.WriteprofileToDisk(ProfileData.FolderPath, ProfileData.Name, profile)
			if profErr != nil {
				c.Flash.Error("Profile size must be less than 1mb")
				fmt.Println("Profile size must be less than 1mb")
				errtmb := os.Remove(ProfileData.FolderPath + "/" + ProfileData.Name)
				if errtmb != nil {
					fmt.Println("Profile file not deleted", errtmb)

				}
				return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + speakerID.String())
			}
		}
	}

	fmt.Print(speakerdata.String())
	c.Flash.Success("Speaker Successfully Updated")

	return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + id)

}
func (c Speakers) Create(confid string, id string) revel.Result {

	fmt.Printf("conferece id is = %v\n", confid)
	var srv = services.SessionService{}
	sessionID, _ := uuid.FromString(id)
	vm, srvError := srv.GetConferenceAndClient(confid)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderError(srvError)
	}
	fmt.Println("sessions", id)
	vm.SessionID = sessionID
	return c.Render(vm)
}
func (c Speakers) Search(Email string, SessionID string, ConferenceID string) revel.Result {
	usrRepo := repositories.Users{}
	imgrepo := repositories.Images{}
	var sessionWeight int = -1
	var confWeight int = -1
	fmt.Println("ccccc", ConferenceID, "hhhhh", SessionID)

	sessionid, _ := uuid.FromString(SessionID)
	confid, _ := uuid.FromString(ConferenceID)
	speakerRepo := repositories.Speaker{}
	vm := viewmodels.UserEditVMRead{}
	user, err := usrRepo.GetByEmail(Email)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("nai mila")
			return c.RenderJSON(map[string]interface{}{
				"Status": "succeeded",
				"data":   "not exist",
			})
		}
		//c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]interface{}{
			"Status": "failed",
			"data":   "system error",
		})
	}
	sessionSpkr, sessionErr := speakerRepo.SessionSpeakerByid(user.ID, sessionid)
	if sessionErr != nil {
		if gorm.IsRecordNotFoundError(sessionErr) {
			fmt.Println("nai mila")
		} else {
			fmt.Println("error in sessionspkr", sessionErr)
			return c.RenderJSON(map[string]interface{}{
				"Status": "failed",
				"data":   "system error",
			})
		}

	} else {
		sessionWeight = sessionSpkr.SortOrder

	}
	confSpkr, sessionErr := speakerRepo.ConferenceSpeakerByid(user.ID, confid)
	if sessionErr != nil {
		if gorm.IsRecordNotFoundError(sessionErr) {
			fmt.Println("nai mila")
		} else {
			fmt.Println("error in sessionspkr", sessionErr)
			return c.RenderJSON(map[string]interface{}{
				"Status": "failed",
				"data":   "system error",
			})
		}

	} else {
		confWeight = confSpkr.SortOrder

	}

	Profile, prerr := imgrepo.GetImage(user.ID, "user", "user_profile")
	if prerr == nil {
		vm.ProfileURL = Profile.BasicURL + Profile.ImageURLPrefix + "/" + Profile.Name

	}
	poster, poserr := imgrepo.GetImage(user.ID, "user", "poster")
	if poserr == nil {
		vm.PosterURL = poster.BasicURL + poster.ImageURLPrefix + "/" + poster.Name

	}
	vm.ID = user.ID
	vm.FirstName = user.FirstName
	vm.LastName = user.LastName
	vm.Email = user.Email
	vm.Organization = user.Organization
	vm.Designation = user.Designation
	vm.PhoneNumber = user.PhoneNumber
	vm.Bio = user.Bio
	vm.Facebook = user.SocialMedia.Facebook
	vm.Youtube = user.SocialMedia.Youtube
	vm.Linkedin = user.SocialMedia.LinkedIn
	vm.Twitter = user.SocialMedia.Twitter
	response := map[string]interface{}{
		"Status":        "succeeded",
		"user":          vm,
		"confWeight":    confWeight,
		"sessionWeight": sessionWeight,
		"data":          "exist",
	}
	fmt.Println("speaker detail", vm)
	return c.RenderJSON(response)
}
func (c Speakers) AddSessionspeaker(
	confid string,
	sessionid string,
	FirstName string,
	LastName string, Email string,
	bio string, PhoneNumber string,
	Designation string,
	Organization string,
	userid string,
	sessionWeight string,
	conferenceWeight string,
	poster []byte,
	profile []byte,
	facebook string,
	twitter string,
	youtube string,
	linkedin string,
) revel.Result {
	Session_sortOrder, _ := strconv.Atoi(sessionWeight)
	conference_sortOrder, _ := strconv.Atoi(conferenceWeight)
	var Speakrsrv = services.SpeakerService{}
	confID, _ := uuid.FromString(confid)
	sessionID, _ := uuid.FromString(sessionid)
	vm := viewmodels.AddSessionSpeakerVMWrite{}
	vm.ConfID = confID
	vm.SessionID = sessionID
	vm.FirstName = FirstName
	vm.LastName = LastName
	vm.Bio = bio
	vm.Designation = Designation
	vm.Organization = Organization
	vm.PhoneNumber = PhoneNumber
	vm.Email = Email
	vm.UserId = userid
	vm.SessionWeight = Session_sortOrder
	vm.ConferenceWeight = conference_sortOrder
	vm.Facebook = facebook
	vm.Twitter = twitter
	vm.Linkedin = linkedin
	vm.Youtube = youtube

	fmt.Println("hyeeid", userid)
	UsrID, DbErr := Speakrsrv.AddsessionSpeaker(vm)
	if DbErr != nil {
		fmt.Printf("conf service returns error %v \n", DbErr)
		return c.RenderText(DbErr.Error())
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
			return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + UsrID.String())
		}
		posterData := viewmodels.ImageVMWrite{
			Name:           UsrID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.userposter.folderbasepath", ""),
			ImageURLPrefix: "user/poster",
			EntityID:       UsrID,
			EntityType:     "user",
			ImageCategory:  "poster",
			IsActive:       true,
		}

		PostrErr := imgSrv.Update(posterData)

		if PostrErr == nil {
			imgsaveerr := imgSrv.WriteFileToDisk(posterData.FolderPath, posterData.Name, poster)
			if imgsaveerr != nil {
				fmt.Println("Poster size must be less than 2mb", imgsaveerr)
				c.Flash.Error("Poster size must be less than 2mb")
				errRmv := os.Remove(posterData.FolderPath + "/" + posterData.Name)
				if errRmv != nil {
					fmt.Println("file not deleted", errRmv)

				}
				return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + UsrID.String())
			}
		}
	}
	// for profile
	if len(c.Params.Files["profile"]) > 0 {
		Profilefile := c.Params.Files["profile"][0].Filename

		// getting thumbnail Extention
		profiletype := mime.TypeByExtension(path.Ext(Profilefile))
		if profiletype == "" {
			// Try to figure out the content type from the data
			profiletype = http.DetectContentType(profile)
		}

		fmt.Println("thumbnailtype1", profiletype)
		// checking valid thumbnail Extention
		ValidProfile := imgSrv.Validation(profiletype)
		if ValidProfile == false {
			fmt.Println("profile type is not valid")
			c.Flash.Error("profile type is not valid")
			return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + UsrID.String())
		}

		ProfileData := viewmodels.ImageVMWrite{
			Name:           UsrID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) +filepath.Ext(Profilefile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.userprofile.folderbasepath", ""),
			ImageURLPrefix: "user/thumbnail",
			EntityID:       UsrID,
			EntityType:     "user",
			ImageCategory:  "user_profile",
			IsActive:       true,
		}

		profileErr := imgSrv.Update(ProfileData)

		if profileErr == nil {
			profErr := imgSrv.WriteprofileToDisk(ProfileData.FolderPath, ProfileData.Name, profile)
			if profErr != nil {
				c.Flash.Error("Profile size must be less than 1mb")
				fmt.Println("Profile size must be less than 1mb")
				errtmb := os.Remove(ProfileData.FolderPath + "/" + ProfileData.Name)
				if errtmb != nil {
					fmt.Println("Profile file not deleted", errtmb)

				}
				return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + UsrID.String())
			}
		}

	}
	c.Flash.Success("Speaker Successfully Updated")

	return c.Redirect("/admin/conferences/" + sessionid + "/speakers/details/" + UsrID.String())
}

func (c Speakers) DeleteSessionSpeaker(sessionId string, id string) revel.Result {
	speakerID, _ := uuid.FromString(id)
	sessioidID, _ := uuid.FromString(sessionId)

	speakerrepo := repositories.Speaker{}
	sessionrepo := repositories.Sessions{}
	DBerr := speakerrepo.DeleteSessionSpeaker(speakerID, sessioidID)
	if DBerr != nil {
		fmt.Println("cant delete speaker from session", DBerr)

	}
	session, sessionerr := sessionrepo.GetByID(sessioidID)
	if DBerr != nil {
		fmt.Println("cant session byid", sessionerr)

	}

	return c.Redirect("/admin/conferences/" + session.ConferenceID.String() + "/sessions/" + sessionId + "/speakers")
}
