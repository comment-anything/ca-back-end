package util

import (
	"testing"
)

func path_expect(domain string, path string, actual PathExtractResult) bool {
	/* fmt.Println()
	fmt.Println(actual) */
	if domain != actual.Domain || path != actual.Path {
		return false
	}
	return true
}

// It would be good if this is well tested. If we have to change it in the future, it could wreak havoc on our entire database.

func TestExtractPathParts(t *testing.T) {
	res := path_expect("google.com", "queries/search.html", ExtractPathParts("http://www.google.com/queries/search.html?q=blablah"))
	if res == false {
		t.Errorf("Failed to extract google.")
	}
	res = path_expect("google.com", "queries/search.html", ExtractPathParts("http://google.com/queries/search.html?q=blablah"))
	if res == false {
		t.Errorf("Failed to extract google.")
	}
	res = path_expect("google.com", "", ExtractPathParts("http://google.com/"))
	if res == false {
		t.Errorf("Failed to extract google.")
	}
	res = path_expect("google.com", "", ExtractPathParts("https://google.com"))
	if res == false {
		t.Errorf("Failed to extract google.")
	}
	res = path_expect("google.com", "search.html", ExtractPathParts("http://google.com/search.html"))
	if res == false {
		t.Errorf("Failed to extract google.")
	}
	res = path_expect("google.com.biz", "search.html", ExtractPathParts("http://google.com.biz/search.html"))
	if res == false {
		t.Errorf("Failed to extract google.com.biz")
	}
	res = path_expect("google.com.biz", "very/nested/path/search.html", ExtractPathParts("http://google.com.biz/very/nested/path/search.html"))
	if res == false {
		t.Errorf("Failed to extract google.")
	}
	res = path_expect("google.com.biz", "very/nested/path/search.php", ExtractPathParts("http://google.com.biz/very/nested/path/search.php?q=srch"))
	if res == false {
		t.Errorf("Failed to extract google.com.biz")
	}
	res = path_expect("abc123.com.biz.ru.ca.uk", "very/nest.ped/patho2/search.php", ExtractPathParts("http://abc123.com.biz.ru.ca.uk/very/nest.ped/patho2/search.php?????wutq=srch"))
	if res == false {
		t.Errorf("Failed to extract abc123.com.biz.ru.ca.uk")
	}

	res = path_expect("abc123.com.biz.ru.ca.uk:25", "very/nest.ped/patho2/search.php", ExtractPathParts("http://abc123.com.biz.ru.ca.uk:25/very/nest.ped/patho2/search.php?????wutq=srch"))
	if res == false {
		t.Errorf("Failed on port")
	}
}
