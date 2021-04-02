package app

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/najamsk/eventvisor/eventvisorHQ/app/controllers"
	"github.com/revel/revel"
	uuid "github.com/satori/go.uuid"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {

	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	//interceptoers
	revel.InterceptFunc(CheckJWT, revel.BEFORE, &controllers.Admin{})

	revel.TemplateFuncs["reverseUrlUuid"] = func(param1 string, param2 uuid.UUID) (template.URL, error) {
		uuidString := param2.String()
		// fmt.Printf("uuid string = %v \n", uuidString)
		reversePath, err := revel.ReverseURL(param1, uuidString)
		return reversePath, err
	}
	revel.TemplateFuncs["reverseUrlUuidTwoparams"] = func(param1 string, param2 uuid.UUID, param3 uuid.UUID) (template.URL, error) {
		uuidString := param2.String()
		uuidString1 := param3.String()
		// fmt.Printf("uuid string = %v \n", uuidString)
		reversePath, err := revel.ReverseURL(param1, uuidString, uuidString1)
		return reversePath, err
	}
	revel.TemplateFuncs["converttostring"] = func(param2 uuid.UUID) string {
		uuidString := param2.String()

		// fmt.Printf("uuid string = %v \n", uuidString)
		return uuidString
	}

	revel.TemplateFuncs["TimeZoneConfig"] = func(sourceTime time.Time) time.Time {

		hqzone := revel.Config.StringDefault("hq.timezone", "Asia/Karachi")
		// fmt.Printf("hqzone = %v\n", hqzone)
		loc, _ := time.LoadLocation(hqzone)
		// fmt.Printf("hazone location  = %v\n", loc)
		// fmt.Printf("source time is   = %v\n", sourceTime)

		//set timezone,
		now := sourceTime.In(loc)
		// fmt.Printf("hazone translated time  = %v\n", now)

		return now
	}

	revel.TemplateFuncs["Appversion"] = func() string {

		appversion := revel.Config.StringDefault("hq.app.version", "")
		fmt.Printf("hqzone = %v\n", appversion)

		return appversion
	}
	revel.TemplateFuncs["YearNow"] = func() int{

		currentTime := time.Now()
		year:=currentTime.Year()
		return year
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// Token JWT claims struct
type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

//CheckJWT interceptor
func CheckJWT(c *revel.Controller) revel.Result {
	fmt.Printf("called *************************************************************** \n")
	fmt.Printf("action name = %v \n", c.Action)

	cookie, err := c.Request.Cookie("jtoken")

	// if cookie got expired it will be nil

	if err != nil || cookie.GetValue() == "" {
		fmt.Printf("cant read cookie for jtoken with error = %v \n", err)
		// c.Flash.Error("password required")
		return c.Redirect(controllers.Account.Login)

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
		return c.Redirect(controllers.Account.Login)
	}

	jwtTokenSecret := []byte(envtokenpass)

	tkn, err := jwt.ParseWithClaims(passedcookie, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("CheckJWT func: There was an error\n")
		}

		return jwtTokenSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Printf("invalid jwt\n")
		}
		fmt.Printf("bad request\n")
		fmt.Printf("bad request eror = %v\n", err.Error())
		return c.Redirect(controllers.Account.Login)
	}
	if !tkn.Valid {
		fmt.Printf("invalid jwt token.valid false\n")
		return c.Redirect(controllers.Account.Login)

	}
	userId := claims.UserID 
	userCookie, err := c.Request.Cookie("userinfo")
	if err != nil || userCookie.GetValue() == "" {
		fmt.Printf("cant read cookie for userinfo at"+c.Action+" with error = %v \n", err)
		// c.Flash.Error("password required")
		return c.Redirect(controllers.Account.Login)

	}
	userinfo := userCookie.GetValue()
	s := strings.Split(userinfo, "&&&")
	fmt.Println("yeh lo cookies", s[0], s[len(s)-1])
	c.ViewArgs["username"] = s[0]
	c.ViewArgs["ImgUrl"] = s[len(s)-1]
	c.ViewArgs["profilesize"]=revel.Config.IntDefault("hq.ProfileImage.size",1)
	c.ViewArgs["imagewidth"]=revel.Config.IntDefault("hq.image.width",764)
	c.ViewArgs["imageheight"]=revel.Config.IntDefault("hq.image.height",1024)
	c.ViewArgs["postersize"]=revel.Config.IntDefault("hq.image.size",2)
	c.ViewArgs["LogedinUserID"] = userId

	return nil
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
