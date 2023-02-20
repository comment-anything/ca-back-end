package util

import (
	"regexp"
	"strings"
)

/**
path.go is for splitting up url strings that comment anywhere receives so it can associate paths, domains, and so forth, and generate unique locations for comments in the database.
*/

// test at: https://regex101.com/r/JYSEte/1
var url_re regexp.Regexp = *regexp.MustCompile(`.*://(www\.)?(?P<domain>[\w\d]*\.[\w\d]*(\.[\w\d]*)*(:\d*)?)/?(?P<path>(([\w\d]([\w\d]\.)*)*/)*([\w\d]*(\.[\w\d]*)*)|$)`)

type PathExtractResult struct {
	Domain  string
	Path    string
	Success bool
}

// ExtractPathParts returns a PathExtractResult from a given url string. It contains the domain name and the path as the second. Success describes whether the URL was valid. It discards query strings ,'www's, and hosts. For example, 'http://www.google.com/queries/search.html?q=blablah' would return google.com, queries/search.html, true
func ExtractPathParts(s string) PathExtractResult {
	result := make(map[string]string)
	match := url_re.FindStringSubmatch(s)
	for i, name := range url_re.SubexpNames() {
		if i != 0 && name != "" {
			if !(i > len(match)) {
				result[name] = match[i]
			}
		}
	}
	var returnResult PathExtractResult
	returnResult.Success = true
	var ok bool
	var dom string
	var path string
	dom, ok = result["domain"]
	if ok {
		path, ok = result["path"]
	}
	if ok {
		returnResult.Domain = strings.ToLower(dom)
		returnResult.Path = strings.ToLower(path)
	} else {
		returnResult.Success = false
	}
	return returnResult
}
