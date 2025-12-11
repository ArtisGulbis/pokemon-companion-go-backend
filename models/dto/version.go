package dto

type Version struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	VersionGroupID int    `json:"versionGroupId"`
}
