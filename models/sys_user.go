package models

type SysUser struct {
	BaseModel
	Name       string       // user name
	NameCn     string       // user cn name
	Email      string       // user email address
	Active     bool         // user is banned, active=0
	LastLogin  JsonTime     // user last login at
	LoginTimes int          // user login times
	RoleIds    mySlice[int] // user roles
}
