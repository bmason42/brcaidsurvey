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

type SurveyContact struct {
	//a generated uuid used for associates
	ID                string
	EncryptionVersion EncryptionVersion
	//the encrypted PiInfo
	PII string
}

type SurveyResult struct {
	//a generated uuid used for associates
	ID string
}

type RegionInfo struct {
	RegionID string `json:"regionID"`
	Name     string `json:"name"`
}
type SupportConcern struct {
	ConcernID string `json:"concernID"`
	Concern   string `json:"concern"`
}

type FormData struct {
	BmRegions []RegionInfo     `json:"bmRegions"`
	Concerns  []SupportConcern `json:"concerns"`
}
