package controllers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/viewmodels"
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"

	//"github.com/najamsk/eventvisor/eventvisorHQ/app/controllers"

	//"encoding/base64"
	//"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

	//"strconv"
	"time"

	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	//"github.com/najamsk/eventvisor/eventvisorHQ/repositories"

	"github.com/revel/revel"
)

// Conference controller
type Users struct {
	Admin
}

// Index action: GET
func (c Users) Index() revel.Result {
	return c.Render()
}
func (c Users) GetSearch() revel.Result {
	return c.Redirect("/admin/users/list")
}

// List action: GET
func (c Users) List(id string) revel.Result {

	var repo = repositories.Users{}
	users := repo.GetAll(0)
	for _, user := range users {

		fmt.Printf("router client loop name = %v\n", user.Roles)
	}
	// return c.RenderText("list

	return c.Render(users)
}
func (c Users) Search(email string) revel.Result {
	fmt.Println("yeh email hye", email)

	var repo = repositories.Users{}
	users, Usrerror := repo.GetAllByEmail(email)
	if Usrerror != nil {
		fmt.Println("user not found by email", Usrerror)
		return c.Redirect("/admin/users/list")
	}
	fmt.Println("yeh users hye", users)
	return c.Render(users)

}
func (c Users) Edit(id string) revel.Result {
	// geting logedin user from jwt start
	logineduser := fmt.Sprint(c.ViewArgs["LogedinUserID"])
	UsrID, _ := uuid.FromString(id)
	rolerepo := repositories.Roles{}
	
	userId,_ :=uuid.FromString(logineduser)
	fmt.Printf("jtoken userid = %v\n", userId)
	// geting logedin user from jwt end
	var srv = services.UserService{}
	vm, srvErr := srv.GetByID(id, userId)
	if srvErr != nil {
		fmt.Println("user srv return an error", srvErr)
		return c.Redirect(Account.Login)

	}
	if vm.LoginRoleweight == 1000 { //admin can edit any user
		fmt.Println("Admin can edit any user")
	} else if vm.UserRoleweight >= vm.LoginRoleweight {
		fmt.Println("user role is greater then logined user role",vm.LoginRoleweight)
		c.Flash.Error("You cant edit this user")
		return c.Redirect("/admin/users/list")

	}
	roleList := rolerepo.GetbyWeight(UsrID, vm.LoginRoleweight)
	if len(roleList) < 1 {
		fmt.Println("role listing returning empty[]")
		c.Flash.Error("You cant edit users")
		return c.Redirect("/admin/users/list")

	}
	vm.Roles = roleList
	fmt.Println("user roles", roleList)
	return c.Render(vm)
}
func (c Users) EditPost(id string,
	firstname string,
	lastname string,
	email string,
	phonenumber string,
	organization string,
	designation string,
	isActive bool,
	bio string,
	poster []byte,
	profile []byte,
	facebook string,
	twitter string,
	youtube string,
	linkedin string,) revel.Result {
	roles := c.Params.Values["roles"]
	fmt.Println("hye roles",c.Params.Values["roles"])

	if len(roles)<1{
		c.Flash.Error("Please assign a role")
		return c.Redirect("/admin/users/detail/" + id)
	}
	userID, _ := uuid.FromString(id)
	UserData := viewmodels.UserEditVMWrite{
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		Organization: organization,
		Designation:  designation,
		PhoneNumber:  phonenumber,
		Bio:          bio,
		Roles:        roles,
		IsActive:     isActive,
		Facebook:facebook,
		Twitter:twitter,
		Youtube:youtube,
		Linkedin:linkedin,
	}
	UserData.ID = userID
	srv := services.UserService{}
	fmt.Printf("user data sending to service:%+v\n", UserData)
	UserID, srvError := srv.UpdateUser(UserData)
	if srvError != nil {
		fmt.Printf("conf service returns error %v \n", srvError)
		return c.RenderText(srvError.Error())
	}
	fmt.Print("hye userid",UserID)
	// getting poster Extention

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
			return c.Redirect("/admin/users/detail/" + id)
		}

		posterData := viewmodels.ImageVMWrite{
			Name:           UserID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(posterfile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.userposter.folderbasepath", ""),
			ImageURLPrefix: "user/poster",
			EntityID:       UserID,
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
				return c.Redirect("/admin/users/detail/" + id)
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
			return c.Redirect("/admin/users/detail/" + id)
		}
		ProfileData := viewmodels.ImageVMWrite{
			Name:           UserID.String() + "_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(Profilefile),
			BasicURL:       revel.Config.StringDefault("hq.image.basicurl", ""),
			FolderPath:     revel.Config.StringDefault("hq.userprofile.folderbasepath", ""),
			ImageURLPrefix: "user/thumbnail",
			EntityID:       UserID,
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
				return c.Redirect("/admin/users/detail/" + id)
			}
		}
	}

	c.Flash.Success("Users Successfully Updated")

	return c.Redirect("/admin/users/detail/" + id)
	// return c.RenderText("mil gia")
}
func (c Users) ChangePassword() revel.Result {
	fmt.Println("ser id",c.ViewArgs["LogedinUserID"])
	return c.Render()
}

func (c Users) PostChangePassword(newpass string, confirmpass string, oldpassword string) revel.Result {
	userrepo := repositories.Users{}
	logineduser := fmt.Sprint(c.ViewArgs["LogedinUserID"])

	if len(newpass) < 6 {
		fmt.Println("Passwords must be 6 or more characters without space.")
		c.Flash.Error("Passwords must be 6 or more characters")
		return c.Redirect("/admin/users/changepassword")

	}
	if newpass != confirmpass {
		fmt.Println("Password not matched with confirm password.")
		c.Flash.Error("Password not matched with confirm password.")
		return c.Redirect("/admin/users/changepassword")
	}
	if len(oldpassword) < 1 {
		fmt.Println("Please provide valid old password.")
		c.Flash.Error("Please provide valid old password.")
		return c.Redirect("/admin/users/changepassword")
	}
	userid,_:=uuid.FromString(logineduser)
	var userdb, err = userrepo.GetByID(userid)
	if err != nil {
		if gorm.IsRecordNotFoundError(err){
			fmt.Println("User doesnot exist")
			c.Flash.Error("User not exist")
			return c.Redirect("/admin/users/changepassword")
	

		}
		fmt.Println("User repo shows error")
		c.Flash.Error("System error")
		return c.Redirect("/admin/users/changepassword")

	}
	if !userdb.IsActive {

		fmt.Println("changePassword:Error: This account is currently not active.")
		c.Flash.Error("This account is currently not active.")
		return c.Redirect("/admin/users/changepassword")

	}
	passerr := bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(oldpassword))
	if passerr != nil {
		fmt.Println("Please provide valid old password.")
		c.Flash.Error("Please provide valid old password.")
		return c.Redirect("/admin/users/changepassword")
	}
	var userepo = userrepo.UpdatePassword(userdb.Email, newpass)
	if userepo != nil {
		fmt.Println("UpdatePassword repo shows error")
		c.Flash.Error("Password cannot update")
		return c.Redirect("/admin/users/changepassword")
	}
	c.Flash.Success("Password updated successfully.")
	return c.Redirect("/admin/users/list")

}

// func IsValidUUID(u string) bool {
// 	var invalidUUID, err1 = uuid.FromString("dummy")
// 	fmt.Println(err1)
// 	fmt.Println(invalidUUID)
// 	id, err := uuid.FromString(u)

// 	return err == nil && id != invalidUUID
// }
