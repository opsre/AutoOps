package model

type LicensePayload struct {
	CustomerName string  `json:"customerName"`
	MachineCode  string  `json:"machineCode"`
	LicenseType  string  `json:"licenseType"`
	IssuedAt     string  `json:"issuedAt"`
	ExpireAt     *string `json:"expireAt"`
}

type LicenseStored struct {
	Activated    bool    `json:"activated"`
	LicenseKey   string  `json:"licenseKey"`
	CustomerName string  `json:"customerName"`
	MachineCode  string  `json:"machineCode"`
	LicenseType  string  `json:"licenseType"`
	IssuedAt     string  `json:"issuedAt"`
	ActivatedAt  string  `json:"activatedAt"`
	ExpireAt     *string `json:"expireAt"`
}

type ActivateLicenseRequest struct {
	LicenseKey string `json:"licenseKey" binding:"required"`
}
