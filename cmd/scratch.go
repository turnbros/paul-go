package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	testString := "kube - argocd"
	fmt.Println(testString)
	testString = strings.Replace(testString, " ", "", -1)
	fmt.Println(testString)
	re := regexp.MustCompile("^[a-zA-Z0-9-]{1,63}")
	match := re.FindStringSubmatch(testString)
	fmt.Println(match)
}
