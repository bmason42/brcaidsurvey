package model

type PiInfo struct {
	//a generated uuid used for associates
	ID               string `json:"id"`
	PlayaName        string `json:"playaName"`
	DefaultWorldName string `json:"defaultWorldName"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
}

type EncryptionVersion int

const (
	ENCRYPT_VERSION_NONE EncryptionVersion = iota
)

type ZipCode struct {
	ZipCode string `gorm:"type:varchar(10);primary_key"`
	State   string
	City    string
	Lat     float64
	Long    float64
	Xaxis   float64
	Yaxis   float64
	Zaxis   float64
}
type SurveyContact struct {
	//a generated uuid used for associates
	SurveyContactID   string `gorm:"type:varchar(36);primary_key"`
	EncryptionVersion EncryptionVersion
	OSPhoneNumber     string `gorm:"type:varchar(12)"`
	//the encrypted PiInfo
	PII string `gorm:"type:text"`
}

type SurveyResult struct {
	//a generated uuid used for associates
	SurveyResultID  string `gorm:"type:varchar(36);primary_key"`
	SurveyContactID string
	Contact         SurveyContact `gorm:"foreignkey:SurveyContactID"`
	//Region RegionInfo `gorm:"foreignkey:RegionID"`
	//Needs []SupportConcern `gorm:"foreignkey:ConcernID;association_foreignkey:SurveyResultID"`
}

//Burning Man region
//These are in the survey.json file
type RegionInfo struct {
	RegionID string `json:"regionID" gorm:"type:varchar(36);primary_key"`
	Name     string `json:"name"`
}

//different support issues such as phyiscal, financial, health, mental health
//These are in the survey.json file
type SupportConcern struct {
	ConcernID string `json:"concernID" gorm:"type:varchar(36);primary_key"`
	Concern   string `json:"concern"`
}

type FormData struct {
	BmRegions []RegionInfo     `json:"bmRegions"`
	Concerns  []SupportConcern `json:"concerns"`
}
