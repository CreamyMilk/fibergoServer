package utils

import (
	"os"
)

func GetROOTDomain() string {

	if root := os.Getenv("ROOT_DOMAIN"); root == "" {
		return "google.coms"
	} else {
		return root
	}

}
