package main

func GetDomains() []string {
	domains, err := readLines("domains.txt")
	if err != nil {
		return nil
	}
	return domains
}
