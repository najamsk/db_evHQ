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

type Sponsors struct {
	Admin
}

func (c Sponsors) Create(id string) revel.Result {

	var serv = services.SponsorService{}
	vm, servError := serv.GetByconference(id)
	if servError != nil {
		fmt.Printf("SponsorService returns error %v \n", servError)
		return c.RenderError(servError)
	}
	fmt.Printf("hye121", vm)
	// return c.RenderText(id)
	return c.Render(vm)
}

func (c Sponsors) ListByConference(id string) revel.Result {

	var serv = services.SponsorService{}
	sponsorVM, servError := serv.GetByconference(id)
	if servError != nil {
		fmt.Printf("SponsorService returns error %v \n", servError)
		return c.RenderError(servError)
	}
	fmt.Printf("hye121", sponsorVM)
	// return c.RenderText(id)
	return c.Render(sponsorVM)
}

func (c Sponsors) Detail(confid string,id string) revel.Result {

	//its working using c.Params.Get for post values as well
	// name := c.Params.Get("name")
	var srv = services.SponsorService{}
	vm, srvError := srv.GetByID(id, confid)
	if srvError != nil {
		fmt.Printf("SponsorService returns error %v \n", srvError)
	}
	fmt.Printf("Sponsor return from service is %v \n ", vm)
	return c.Render(vm)
	// return c.RenderText("mil gia")
}
func (c Sponsors) Update(
	id string,
	confiD string,
	name string,
	details string,
	clientid string,
	sortorder string,
	logo []byte,
	facebook string,
	twitter string,
	youtube string,
	linkedin string,
	sponlevel string,
	IsActive bool,

) revel.Result {

	clientID, _ := uuid.FromString(clientid)
	confID, _ := uuid.FromString(confiD)
	spnID, _ := uuid.FromString(id)
	sponlevelID, _ := uuid.FromString(sponlevel)
	sortOrder, _ := strconv.Atoi(sortorder)
	fmt.Println("name11", name)
	spnData := viewmodels.SponsorEditVmWrite{
		Name:         name,
		Bio:          details,
		ClientID:     clientID,
		SortOrder:    sortOrder,
		ConferenceID: confID,
		Facebook:     facebook,
		Twitter:      twitter,
		Linkedin:     linkedin,
		Youtube:      youtube,
		Sponlevel:    sponlevelID,
		IsActive:     IsActive,
	}
	spnData.ID = spnID
	srv := services.SponsorService{}
	speakerdata, srvError := srv.Update(spnData)
	if srvError != nil {
		fmt.Printf("SponsorService returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}

	imgSrv := services.ImageService{}

	if len(c.Params.Files["logo"]) > 0 {
		// for profile
		logofile := c.Params.Files["logo"][0].Filename

		// getting thumbnail Extention
		logotype := mime.TypeByExtension(path.Ext(logofile))
		if logotype == "" {
			// Try to figure out the content type from the data
			logotype = http.DetectContentType(logo)
		}

		fmt.Println("thumbnailtype1", logotype)
		// checking valid thumbnail Extention
		Validlogo := imgSrv.Validation(logotype)
		if Validlogo == false {
			fmt.Println("profile type is not valid")
			c.Flash.Error("profile type is not valid")
			return c.Redirect("/admin/conferences/" + confiD + "/sponsors/" + spnID.String())
		}
		logoData := viewmodels.ImageVMWrite{
			Name:           spnID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(logofile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.sponsor.folderbasepath", ""),
			ImageURLPrefix: "sponsor/logo",
			EntityID:       spnID,
			EntityType:     "sponsor",
			ImageCategory:  "logo",
			IsActive:       true,
		}

		profileErr := imgSrv.Update(logoData)

		if profileErr == nil {
			profErr := imgSrv.WriteprofileToDisk(logoData.FolderPath, logoData.Name, logo)
			if profErr != nil {
				c.Flash.Error("logo size must be less than 1mb")
				fmt.Println("logo size must be less than 1mb")
				errtmb := os.Remove(logoData.FolderPath + "/" + logoData.Name)
				if errtmb != nil {
					fmt.Println("logo file not deleted", errtmb)

				}
				return c.Redirect("/admin/conferences/" + confiD + "/sponsors/" + spnID.String())
			}
		}
	}
	fmt.Print(speakerdata.String())
	c.Flash.Success("sponsor Successfully Updated")

	return c.Redirect("/admin/conferences/" + confiD + "/sponsors/" + spnID.String())

}

func (c Sponsors) CreatePost(
	confiD string,
	name string,
	details string,
	clientid string,
	sortorder string,
	logo []byte,
	facebook string,
	twitter string,
	youtube string,
	linkedin string,
	sponlevel string,
	IsActive bool,

) revel.Result {

	clientID, _ := uuid.FromString(clientid)
	confID, _ := uuid.FromString(confiD)
	sponlevelID, _ := uuid.FromString(sponlevel)
	sortOrder, _ := strconv.Atoi(sortorder)
	fmt.Println("name11", name)
	spnData := viewmodels.SponsorEditVmWrite{
		Name:         name,
		Bio:          details,
		ClientID:     clientID,
		SortOrder:    sortOrder,
		ConferenceID: confID,
		Facebook:     facebook,
		Twitter:      twitter,
		Linkedin:     linkedin,
		Youtube:      youtube,
		Sponlevel:    sponlevelID,
		IsActive:     IsActive,
	}
	srv := services.SponsorService{}
	spnID, srvError := srv.Create(spnData)
	if srvError != nil {
		fmt.Printf("SponsorService returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}

	imgSrv := services.ImageService{}

	if len(c.Params.Files["logo"]) > 0 {
		// for profile
		logofile := c.Params.Files["logo"][0].Filename

		// getting thumbnail Extention
		logotype := mime.TypeByExtension(path.Ext(logofile))
		if logotype == "" {
			// Try to figure out the content type from the data
			logotype = http.DetectContentType(logo)
		}

		fmt.Println("thumbnailtype1", logotype)
		// checking valid thumbnail Extention
		Validlogo := imgSrv.Validation(logotype)
		if Validlogo == false {
			fmt.Println("profile type is not valid")
			c.Flash.Error("profile type is not valid")
			return c.Redirect("/admin/conferences/" + confiD + "/sponsors/" + spnID.String())
		}
		logoData := viewmodels.ImageVMWrite{
			Name:           spnID.String() + filepath.Ext(logofile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.sponsor.folderbasepath", ""),
			ImageURLPrefix: "sponsor/logo",
			EntityID:       spnID,
			EntityType:     "sponsor",
			ImageCategory:  "logo",
			IsActive:       true,
		}

		profileErr := imgSrv.CreateImage(logoData)

		if profileErr == nil {
			profErr := imgSrv.WriteprofileToDisk(logoData.FolderPath, logoData.Name, logo)
			if profErr != nil {
				c.Flash.Error("logo size must be less than 1mb")
				fmt.Println("logo size must be less than 1mb")
				errtmb := os.Remove(logoData.FolderPath + "/" + logoData.Name)
				if errtmb != nil {
					fmt.Println("logo file not deleted", errtmb)

				}
				return c.Redirect("/admin/conferences/" + confiD + "/sponsors/" + spnID.String())
			}
		}
	}
	c.Flash.Success("sponsor Successfully Updated")

	return c.Redirect("/admin/conferences/" + confiD + "/sponsors/" + spnID.String())

}
