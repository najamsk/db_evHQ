package controllers

import (
	"fmt"
	//"strconv"
	//"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	uuid "github.com/satori/go.uuid"

	//"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	//"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"github.com/revel/revel"
	//uuid "github.com/satori/go.uuid"
)

// Conference controller
type Attendies struct {
	Admin
}

// Index action: GET
func (c Attendies) Index() revel.Result {
	return c.Render()
}

// List action: GET
func (c Attendies) ListByConference(id string) revel.Result {

	var serv = services.AttendiesService{}
	attendiesVM, servError := serv.GetAttendiesByConferenceID(id)
	if servError != nil {
		fmt.Printf("conf service returns error %v \n", servError)
		return c.RenderError(servError)
	}

	// return c.RenderText(id)

	return c.Render(attendiesVM)
}

//Edit action
func (c Attendies) Edit(id string, confid string) revel.Result {
	var srv = services.AttendiesService{}
	attendies, srvError := srv.GetAttendiesByID(id, confid)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		// return c.RenderText(srvError.Error())
	}
	fmt.Printf("confernce return from service is %v \n ", attendies)

	return c.Render(attendies)

}

//EditPost func
func (c Attendies) EditPost(confid string, id string, firstname string, lastname string, email string, designation string, phonenumber string, organization string, bio string) revel.Result {
	attendeeID, _ := uuid.FromString(id)

	fmt.Printf("form id = %v \n", id)
	fmt.Printf("form confiD = %v \n", confid)
	fmt.Printf("form firstname = %v \n", firstname)
	fmt.Printf("form lastname = %v \n", lastname)
	fmt.Printf("form email = %v \n", email)
	fmt.Printf("form designation = %v \n", designation)
	fmt.Printf("form phonenumber = %v \n", phonenumber)
	fmt.Printf("form organization = %v \n", organization)
	fmt.Printf("form bio = %v \n", bio)

	vm := viewmodels.AttendiesEditVMWrite{
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		Designation:  designation,
		PhoneNumber:  phonenumber,
		Organization: organization,
		Bio:          bio,
	}
	vm.ID = attendeeID

	var srv = services.AttendiesService{}
	srvError := srv.UpdateAttendee(vm)

	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		url := "/admin/conferences/" + confid + "/attendies/details/" + id + "?flashstatus=false&flashmsg=Error while updating. Sorry!"

		c.Flash.Error("Sorry found errors while updating.")
		// return c.RenderText(srvError.Error())
		return c.Redirect(url)

	}
	url := "/admin/conferences/" + confid + "/attendies/details/" + id + "?flashstatus=true&flashmsg=Updated Successfully."
	c.Flash.Success("Updated successfully.")

	return c.Redirect(url)
	// return c.Render(vm)

}
