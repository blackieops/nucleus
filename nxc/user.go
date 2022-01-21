package nxc

// {"ocs":{"meta":{"status":"ok","statuscode":100,"message":"OK","totalitems":"","itemsperpage":""},"data":{"enabled":true,"storageLocation":"\/var\/www\/html\/data\/alex","id":"alex","lastLogin":1642731498000,"backend":"Database","subadmin":[],"quota":{"free":41555959808,"used":23087106,"total":41579046914,"relative":0.06,"quota":-3},"avatarScope":"v2-federated","email":null,"emailScope":"v2-federated","additional_mail":[],"additional_mailScope":[],"displaynameScope":"v2-federated","phone":"","phoneScope":"v2-local","address":"","addressScope":"v2-local","website":"","websiteScope":"v2-local","twitter":"","twitterScope":"v2-local","organisation":"","organisationScope":"v2-local","role":"","roleScope":"v2-local","headline":"","headlineScope":"v2-local","biography":"","biographyScope":"v2-local","profile_enabled":"1","profile_enabledScope":"v2-local","groups":["admin"],"language":"en","locale":"","notify_email":null,"backendCapabilities":{"setDisplayName":true,"setPassword":true},"display-name":"alex"}}}%

type UserResponse struct {
	Enabled     bool   `json:"enabled"`
	UserId      string `json:"id"`
	DisplayName string `json:"display-name"` // just *why*
}
