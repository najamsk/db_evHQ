package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/najamsk/eventvisor/eventvisorHQ/repositories"
	"github.com/najamsk/eventvisor/eventvisorHQ/services"
	"github.com/najamsk/eventvisor/eventvisorHQ/utils"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Token JWT claims struct
type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
	// UserRoles []*models.Role
	// MaxRoleWeight int32
	// MaxRoleName   string
	// MaxRoleID     uuid.UUID
}

// Account controller
type Account struct {
	*revel.Controller
}

// MyData receving curl request body json data
type MyData struct {
	Email string `json:"email"`
}

// Index action: GET
func (c Account) Index() revel.Result {
	// fmt.Printf("c.request = %+v \n", c.Request)

	// cookie, err := c.Request.Cookie("jtoken")

	// if err != nil {
	// 	fmt.Printf("cant read cookie for jtoken with error = %v \n", err)
	// }

	// fmt.Printf("cookie is = %v", cookie.GetValue())

	return c.Redirect(Account.Login)
}

func (c Account) UpdatePassword() revel.Result {
	return c.Render()
}

// Login action: GET
func (c Account) Login() revel.Result {
	// fmt.Printf("c.request = %+v \n", c.Request)

	cookie, err := c.Request.Cookie("jtoken")

	// if cookie got expired it will be nil

	if err != nil || cookie.GetValue() == "" {
		fmt.Printf("cant read cookie for jtoken with error = %v \n", err)
		return c.Render()
	}

	fmt.Printf("cookie is = %v", cookie.GetValue())

	passedcookie := cookie.GetValue()
	claims := &Token{}

	fmt.Printf("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&\n")
	fmt.Printf("jtoken = %v \n", passedcookie)

	envtokenpass := revel.Config.StringDefault("hq.jwtsecret", "")

	if envtokenpass == "" {
		fmt.Printf("empty jwt token pass \n")
		c.Flash.Error("Login service failed.")
		return c.Redirect(Account.Login)
	}

	jwtTokenSecret := []byte(envtokenpass)

	tkn, err := jwt.ParseWithClaims(passedcookie, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error\n")
		}

		return jwtTokenSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Printf("invalid jwt\n")
		}
		fmt.Printf("bad request\n")
		fmt.Printf("bad request eror = %v\n", err.Error())
		return c.Render()
	}
	if !tkn.Valid {
		fmt.Printf("invalid jwt token.valid false\n")
		return c.Render()
	}

	userId := claims.UserID // logedin user id
	fmt.Printf("jtoken userid = %v\n", userId)
	//parse jwt token out from cookie

	//send user to admin/clients list page

	return c.Redirect(Clients.Index)
}

// ResetPassword action: GET
func (c Account) ResetPassword() revel.Result {
	return c.Render()
}

// LoginPost action: POST
func (c Account) LoginPost(email string, password string) revel.Result {

	fmt.Printf("email = %v \n", email)
	fmt.Printf("password = %v \n", password)

	email = strings.TrimSpace(email)
	//validating inputs - email, password
	if len(email) <= 0 {
		c.Flash.Error("Email required")
		return c.Redirect(Account.Login)
	}

	if len(password) <= 0 {
		c.Flash.Error("password required")
		return c.Redirect(Account.Login)
	}

	var serv = services.AccountService{}
	userdb, servError := serv.GetByEmail(email)
	if servError != nil {
		if servError.Error() == "record not found" {
			c.Flash.Error("User credentials are not valid.")
			return c.Redirect(Account.Login)
		}
		c.Flash.Error("sorry cant process your request")
		return c.Redirect(Account.Login)
	}

	fmt.Printf("user return from service inside controller ation is = %+v \n", userdb)

	data := MyData{}
	c.Params.BindJSON(&data)

	//n@j@m123
	dbpass := userdb.Password
	err := bcrypt.CompareHashAndPassword([]byte(dbpass), []byte(password))
	if err != nil {
		c.Flash.Error("wrong email or password.")
		return c.Redirect(Account.Login)
	}

	fmt.Printf("password match for passed in email address \n")

	// lets create jwt token now

	//TODO: check if any role exist that allows HQ access if yes set hqaccess to true

	hqroles, hqroleserr := serv.HQRoles()

	if hqroleserr != nil {
		fmt.Printf("acount service hqroles returns error = %v\n", hqroleserr)
		c.Flash.Error("can't authorize.")

		return c.Redirect(Account.Login)
	}

	// userdb.roles
	counter := len(userdb.Roles)
	fmt.Printf("number or roles for user %v\n", counter)

	counter = len(hqroles)
	fmt.Printf("number or hqroles %v\n", counter)
	for hqKey, hqValue := range hqroles {
		fmt.Println("Key:", hqKey, ",Value:", hqValue)
	}

	//checking if user roles have any role that can match with HQ Roles list
	hqRolePresent := false
	for _, role := range userdb.Roles {

		fmt.Printf(" role:name = %v, weight=%v\n", role.Name, role.Weight)
		if _, hqRolePresent = hqroles[role.Name]; hqRolePresent {
			fmt.Println("found one hq level role")
			break
		}
	}

	if hqRolePresent == false {
		fmt.Println("You don't have sufficient roles.")
		c.Flash.Error("You don't have sufficient roles.")
		return c.Redirect(Account.Login)
	}

	//TODO: jwt token expiry  - expired on ? or created on

	expirationTime := time.Now().Add(60 * time.Minute)
	tk := &Token{UserID: userdb.ID} // UserRoles: userdb.Roles,
	tk.StandardClaims.ExpiresAt = expirationTime.Unix()
	// MaxRoleName:   userdb.Roles[0].Name,
	// MaxRoleWeight: userdb.Roles[0].Weight,
	// MaxRoleID:     userdb.Roles[0].ID,

	fmt.Printf("token we build is =  %+v \n", tk)

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	envtokenpass := revel.Config.StringDefault("hq.jwtsecret", "")

	if envtokenpass == "" {
		fmt.Printf("empty jwt token pass \n")
		c.Flash.Error("Login service failed.")
		return c.Redirect(Account.Login)
	}

	// envtokenpass := os.Getenv("token_password")
	fmt.Printf("token pass from env file = %v \n", envtokenpass)

	tokenString, errTokenString := token.SignedString([]byte(envtokenpass))

	if errTokenString != nil {
		fmt.Printf("Token string error: %v\n", errTokenString.Error())
	}

	fmt.Printf("jwt token string is = %+v \n", tokenString)

	// c.Response.Out.Header().Set("WWW-Authenticate", `Basic realm="revel"`)
	imageService := services.ImageService{}
	var imgUrl string
	isImage := true

	var imgDB, Imgerr = imageService.GetImage(userdb.ID.String(), "user", "user_profile")
	if Imgerr != nil {
		if Imgerr.Error() == "record not found" {
			isImage = false
			fmt.Println("Image not exist")
		} else {
			c.Flash.Error("sorry cant process your request")
			return c.Redirect(Account.Login)
		}
	}
	if userdb.FirstName == "" {
		fmt.Println("jhghjgjhghjgjg")
	}

	if isImage == true {
		imgUrl = imgDB.BasicURL + imgDB.ImageURLPrefix + "/" + imgDB.Name
	} else {
		imgUrl = "noimage"
	}
	newcookie := &http.Cookie{
		Name:    "jtoken",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	}
	username := strings.Title(userdb.FirstName + " " + userdb.LastName)
	userInfo := &http.Cookie{
		Name:    "userinfo",
		Value:   username + "&&&" + imgUrl,
		Expires: expirationTime,
		Path:    "/",
	}

	c.SetCookie(newcookie)
	c.SetCookie(userInfo)
	return c.Redirect(Clients.Index)
	// return c.RenderText(userdb.FirstName)
}

// LoginPost action: POST
// func (c Account) LoginPost(email string, password string) revel.Result {

// 	data := MyData{}
// 	c.Params.BindJSON(&data)

// 	// myval := "hello mera bahi "

// 	fmt.Println("email is :", data.Email)

// 	return c.RenderText(data.Email + " ... oye")
// }

// Logout action:GET
func (c Account) Logout() revel.Result {

	// newcookie2 := http.Cookie{
	// 	Name:   "ithinkidroppedacookie",
	// 	MaxAge: -1}

	newcookie := &http.Cookie{
		Name:    "jtoken",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	userInfo := &http.Cookie{
		Name:    "userinfo",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	c.SetCookie(userInfo)
	c.SetCookie(newcookie)
	// return c.RenderText("shit")
	return c.Redirect(Account.Login)
}
func (c Account) ForgotPassword(email string) revel.Result {
	var serv = services.AccountService{}
	mailer := utils.Mailer{}

	//fmt.Println("id:",userobj.Email)
	fmt.Println("aye meri mail hye", email)
	Email := strings.TrimSpace(email)

	if len(Email) < 1 {
		fmt.Println("Please provide valid email address.")
		c.Flash.Error("Please provide valid email address.")
		return c.Redirect("/account/reset")

	}

	var userdb, usrerr = serv.GetByEmail(Email)
	fmt.Println(usrerr)
	if usrerr != nil {
		if gorm.IsRecordNotFoundError(usrerr) {
			fmt.Println("Wrong email. Please try again.")
			c.Flash.Error("Please provide valid email address.")
			return c.Redirect("/account/reset")
		}
		fmt.Println("Get user by email shows error", usrerr)
		c.Flash.Error("system err")
		return c.Redirect("/account/reset")

	}
	if !userdb.IsActive {

		fmt.Println("forgotPassword:Error: This account is currently not active.")
		c.Flash.Error("This account is currently not active.")
		return c.Redirect("/account/reset")

	}

	myrand := rand.Intn(999999-100000) + 100000
	fmt.Println(myrand)
	code := strconv.FormatInt(int64(myrand), 10)
	upderr := serv.UpdateResetPassword(Email, code)
	if upderr != nil {
		fmt.Println("UpdateResetPassword shows err", upderr)
		c.Flash.Error("Password reset request not completed. Please try again")
		return c.Redirect("/account/reset")
	}
	fmt.Println("mycode", code)
	type TemplateData struct {
		FirstName string
		LastName  string
		Code      string
	}

	data := TemplateData{userdb.FirstName, userdb.LastName, code}
	var body string
	body, err := mailer.ParseEmailTemplate(data, "forgotPassword.html")
	if err != nil {
		fmt.Println("mailer.ParseEmailTemplate show err", err)
		c.Flash.Error("Password reset request not completed. Please try again")
		return c.Redirect("/account/reset")

	}
	mailsndEror := mailer.SendEmail([]string{Email}, []string{}, "Forgot Password", body)
	if mailsndEror != nil {
		fmt.Println("mailer.SendEmail show err")
		c.Flash.Error("Password reset request not completed. Please try again")
		return c.Redirect("/account/reset")
	}
	c.Flash.Success("Code successfully sent to your email address")
	return c.Redirect("/account/Update")

}

func (c Account) ResetPasswordPost(email string, pass string, confirmpass string, passtkn string) revel.Result {

	//fmt.Println("id:",userobj.Email)
	accountRepo := repositories.Accounts{}
	UserRepo:=repositories.Users{}
	Email := strings.TrimSpace(email)
	Password := strings.TrimSpace(pass)
	ConfirmPassword := strings.TrimSpace(confirmpass)

	if len(Email) < 1 {
		fmt.Println("resetPassword:Error: Please provide valid email address.")
		c.Flash.Error("Please provide valid email address.")
		return c.Redirect("/account/Update")
	}
	if len(Password) < 6 {
		fmt.Println("resetPassword:Error: Please provide valid password.")
		c.Flash.Error("Passwords must be 6 or more characters.")
		return c.Redirect("/account/Update")

	}

	if Password != ConfirmPassword {
		fmt.Println("resetPassword:Error: Password not matched with confirm password.")
		c.Flash.Error("Password not matched with confirm password")
		return c.Redirect("/account/Update")

	}
	if len(passtkn) < 1 {
		fmt.Println("resetPassword:Error: Please provide valid token.")
		c.Flash.Error("Please provide valid token")
		return c.Redirect("/account/Update")
	}

	expiryMin, passErr := accountRepo.GetPasswordTokenExpiryInMin(Email, passtkn) //email string, code string
	fmt.Println("expiryMin:", expiryMin)

	if passErr != nil {
		if gorm.IsRecordNotFoundError(passErr) {
			fmt.Println("resetPassword:Error: record not found.",passErr)
			c.Flash.Error("Record not found")
			return c.Redirect("/account/Update")
		}
		fmt.Println("Get user by email shows error", passErr)
		c.Flash.Error("System Error.")
		return c.Redirect("/account/Update")
	}
	if expiryMin > revel.Config.IntDefault("hq.expiryminute", 15) {

		fmt.Println("resetPassword:Error: forgot password token expired.")
		c.Flash.Error("Forgot password token expired")
		return c.Redirect("/account/Update")
	}
	var userepo = UserRepo.UpdatePassword(Email,Password)
	if userepo != nil {
		fmt.Println("UpdatePassword repo shows error")
		c.Flash.Error("Password cannot update")
		return c.Redirect("/account/Update")
	}
	c.Flash.Success("Password Successfully Updated")
	return c.Redirect("/account/login")
}
