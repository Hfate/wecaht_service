package utils

var (
	IdVerify               = Rules{"ID": []string{NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	LoginVerify            = Rules{"CaptchaId": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify         = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}

	MediaVerify            = Rules{"TargetAccountId": {NotEmpty()}}
	TopicVerify            = Rules{"Topic": {NotEmpty()}}
	PortalVerify           = Rules{"PortalName": {NotEmpty()}, "PortalKey": {NotEmpty()}, "ArticleKey": {NotEmpty()}, "Link": {NotEmpty()}, "GraphQuery": {NotEmpty()}, "TargetNum": {NotEmpty()}}
	OfficialAccountVerify  = Rules{"AccountName": {NotEmpty()}, "UserEmail": {NotEmpty()}, "Topic": {NotEmpty()}}
	BenchmarkAccountVerify = Rules{"AccountName": {NotEmpty()}, "Topic": {NotEmpty()}, "ArticleLink": {NotEmpty()}}
	WxTokenVerify          = Rules{"SlaveSid": {NotEmpty()}, "BizUin": {NotEmpty()}, "DataTicket": {NotEmpty()}, "RandInfo": {NotEmpty()}, "Token": {NotEmpty()}}
	PromptVerify           = Rules{"Topic": {NotEmpty()}, "PromptType": {NotEmpty()}, "Prompt": {NotEmpty()}, "Language": {NotEmpty()}}
)
