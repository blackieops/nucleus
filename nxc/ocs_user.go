package nxc

// {
// "ocs": {
//   "meta": {
//     "status": "ok",
//     "statuscode": 100,
//     "message": "OK",
//     "totalitems": "",
//     "itemsperpage": ""
//   },
//   "data": {
//     "id": "admin",
//     "lastLogin": 1705158758000,
//     "backend": "Database",
//     "subadmin": [],
//     "quota": {
//       "free": 7029044543488,
//       "used": 362072964676,
//       "total": 7391117508164,
//       "relative": 4.9,
//       "quota": -3
//     },
//     "manager": "",
//     "avatarScope": "v2-federated",
//     "email": "admin@example.com",
//     "emailScope": "v2-federated",
//     "additional_mail": [],
//     "additional_mailScope": [],
//     "displayname": "Admin",
//     "display-name": "Admin",
//     "displaynameScope": "v2-federated",
//     "phone": "",
//     "phoneScope": "v2-local",
//     "address": "",
//     "addressScope": "v2-local",
//     "website": "",
//     "websiteScope": "v2-local",
//     "twitter": "",
//     "twitterScope": "v2-local",
//     "fediverse": "",
//     "fediverseScope": "v2-local",
//     "organisation": "",
//     "organisationScope": "v2-local",
//     "role": "",
//     "roleScope": "v2-local",
//     "headline": "",
//     "headlineScope": "v2-local",
//     "biography": "",
//     "biographyScope": "v2-local",
//     "profile_enabled": "1",
//     "profile_enabledScope": "v2-local",
//     "groups": [],
//     "language": "en",
//     "locale": "",
//     "notify_email": null,
//     "backendCapabilities": {
//       "setDisplayName": true,
//       "setPassword": true
//     }
//   }
// }
// }

type OCSCloudUser struct {
	XMLName struct{}         `xml:"ocs"`
	Meta    OCSCloudUserMeta `json:"meta" xml:"meta"`
	Data    OCSCloudUserData `json:"data" xml:"data"`
}

type OCSCloudUserMeta struct {
	Status       string `json:"status" xml:"status"`
	StatusCode   int    `json:"statuscode" xml:"statuscode"`
	Message      string `json:"message" xml:"message"`
	TotalItems   string `json:"totalitems" xml:"totalitems"`
	ItemsPerPage string `json:"itemsperpage" xml:"itemsperpage"`
}

type OCSCloudUserData struct {
	ID                  string                `json:"id" xml:"id"`
	LastLogin           int64                 `json:"lastLogin" xml:"lastLogin"`
	Backend             string                `json:"backend" xml:"backend"`
	SubAdmin            []string              `json:"subadmin" xml:"subadmin"`
	Quota               OCSCloudUserDataQuota `json:"quota" xml:"quota"`
	Manager             []string              `json:"manager" xml:"manager"`
	AvatarScope         string                `json:"avatarScope" xml:"avatarScope"`
	Email               string                `json:"email" xml:"email"`
	EmailScope          string                `json:"emailScope" xml:"emailScope"`
	AdditionalMail      []string              `json:"additional_mail" xml:"additional_mail"`
	AdditionalMailScope []string              `json:"additional_mailScope" xml:"additional_mailScope"`
	DisplayName         string                `json:"displayname" xml:"displayname"`
	// jesus christ, nextcloud
	DisplayName2        string                              `json:"display-name" xml:"display-name"`
	DisplayNameScope    string                              `json:"displaynameScope" xml:"displaynameScope"`
	Phone               string                              `json:"phone" xml:"phone"`
	PhoneScope          string                              `json:"phoneScope" xml:"phoneScope"`
	Address             string                              `json:"address" xml:"address"`
	AddressScope        string                              `json:"addressScope" xml:"addressScope"`
	Website             string                              `json:"website" xml:"website"`
	WebsiteScope        string                              `json:"websiteScope" xml:"websiteScope"`
	Twitter             string                              `json:"twitter" xml:"twitter"`
	TwitterScope        string                              `json:"twitterScope" xml:"twitterScope"`
	Fediverse           string                              `json:"fediverse" xml:"fediverse"`
	FediverseScope      string                              `json:"fediverseScope" xml:"fediverseScope"`
	Organisation        string                              `json:"organisation" xml:"organisation"`
	OrganisationScope   string                              `json:"organisationScope" xml:"organisationScope"`
	Role                string                              `json:"role" xml:"role"`
	RoleScope           string                              `json:"roleScope" xml:"roleScope"`
	Headline            string                              `json:"headline" xml:"headline"`
	HeadlineScope       string                              `json:"headlineScope" xml:"headlineScope"`
	Biography           string                              `json:"biography" xml:"biography"`
	BiographyScope      string                              `json:"biographyScope" xml:"biographyScope"`
	ProfileEnabled      int                                 `json:"profile_enabled" xml:"profile_enabled"`
	ProfileEnabledScope string                              `json:"profile_enabledScope" xml:"profile_enabledScope"`
	Groups              []string                            `json:"groups" xml:"groups"`
	Language            string                              `json:"language" xml:"language"`
	Locale              string                              `json:"locale" xml:"locale"`
	NotifyEmail         *string                             `json:"notify_email"`
	BackendCapabilities OCSCloudUserDataBackendCapabilities `json:"backendCapabilities" xml:"backendCapabilities"`
}

type OCSCloudUserDataQuota struct {
	Free     int64   `json:"free" xml:"free"`
	Used     int64   `json:"used" xml:"used"`
	Total    int64   `json:"total" xml:"total"`
	Relative float64 `json:"relative" xml:"relative"`
	Quota    int     `json:"quota" xml:"quota"`
}

type OCSCloudUserDataBackendCapabilities struct {
	SetDisplayName bool `json:"setDisplayName" xml:"setDisplayName"`
	SetPassword    bool `json:"setPassword" xml:"setPassword"`
}

type OCSCloudUserInput struct {
	QuotaTotal    int64
	QuotaUsed     int64
	QuotaFree     int64
	DisplayName   string
	Username      string
	EmailAddress  string
	LastLoginTime int64
}

func BuildOCSCloudUserResponse(input OCSCloudUserInput) OCSCloudUser {
	return OCSCloudUser{
		Meta: OCSCloudUserMeta{
			Status:       "ok",
			StatusCode:   200,
			Message:      "OK",
			TotalItems:   "",
			ItemsPerPage: "",
		},
		Data: OCSCloudUserData{
			ID:        input.Username,
			LastLogin: input.LastLoginTime,
			Backend:   "Database",
			SubAdmin:  []string{},
			Quota: OCSCloudUserDataQuota{
				Total: input.QuotaTotal,
				Used:  input.QuotaUsed,
				Free:  input.QuotaFree,
				// TODO:
				Relative: 0.1,
				Quota:    -3,
			},
			Manager:             []string{},
			AvatarScope:         "v2-local",
			Email:               input.EmailAddress,
			EmailScope:          "v2-local",
			AdditionalMail:      []string{},
			AdditionalMailScope: []string{},
			DisplayName:         input.DisplayName,
			DisplayName2:        input.DisplayName,
			DisplayNameScope:    "v2-local",
			Phone:               "",
			PhoneScope:          "v2-local",
			Address:             "",
			AddressScope:        "v2-local",
			Website:             "",
			WebsiteScope:        "v2-local",
			Twitter:             "",
			TwitterScope:        "v2-local",
			Fediverse:           "",
			FediverseScope:      "v2-local",
			Organisation:        "",
			OrganisationScope:   "v2-local",
			Role:                "",
			RoleScope:           "v2-local",
			Headline:            "",
			HeadlineScope:       "v2-local",
			Biography:           "",
			BiographyScope:      "v2-local",
			ProfileEnabled:      1,
			ProfileEnabledScope: "v2-local",
			Groups:              []string{},
			Language:            "en",
			Locale:              "",
			NotifyEmail:         nil,
			BackendCapabilities: OCSCloudUserDataBackendCapabilities{
				SetDisplayName: false,
				SetPassword:    false,
			},
		},
	}
}

// <?xml version="1.0"?>
// <ocs>
//  <meta>
//   <status>ok</status>
//   <statuscode>200</statuscode>
//   <message>OK</message>
//  </meta>
//  <data>
//   <id>alexblackie</id>
//   <lastLogin>1705158758000</lastLogin>
//   <backend>Database</backend>
//   <subadmin/>
//   <quota>
//    <free>7029044543488</free>
//    <used>362072964676</used>
//    <total>7391117508164</total>
//    <relative>4.9</relative>
//    <quota>-3</quota>
//   </quota>
//   <manager></manager>
//   <avatarScope>v2-federated</avatarScope>
//   <email>alex@blackie.me</email>
//   <emailScope>v2-federated</emailScope>
//   <additional_mail/>
//   <additional_mailScope/>
//   <displayname>Alex Blackie</displayname>
//   <display-name>Alex Blackie</display-name>
//   <displaynameScope>v2-federated</displaynameScope>
//   <phone></phone>
//   <phoneScope>v2-local</phoneScope>
//   <address></address>
//   <addressScope>v2-local</addressScope>
//   <website></website>
//   <websiteScope>v2-local</websiteScope>
//   <twitter></twitter>
//   <twitterScope>v2-local</twitterScope>
//   <fediverse></fediverse>
//   <fediverseScope>v2-local</fediverseScope>
//   <organisation></organisation>
//   <organisationScope>v2-local</organisationScope>
//   <role></role>
//   <roleScope>v2-local</roleScope>
//   <headline></headline>
//   <headlineScope>v2-local</headlineScope>
//   <biography></biography>
//   <biographyScope>v2-local</biographyScope>
//   <profile_enabled>1</profile_enabled>
//   <profile_enabledScope>v2-local</profile_enabledScope>
//   <groups/>
//   <language>en</language>
//   <locale></locale>
//   <notify_email/>
//   <backendCapabilities>
//    <setDisplayName>1</setDisplayName>
//    <setPassword>1</setPassword>
//   </backendCapabilities>
//  </data>
// </ocs>
