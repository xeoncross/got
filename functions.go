package got

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
)

// Functions I've found to be required in most every web-site template engine
// Many borrowed from https://github.com/Masterminds/sprig

// DefaultFunctions for templates
var DefaultFunctions = template.FuncMap{
	"title":     strings.Title,
	"upper":     strings.ToUpper,
	"lower":     strings.ToLower,
	"trim":      strings.TrimSpace,
	"urlencode": url.QueryEscape,
	// Often used for tables of rows
	"yesno": func(yes string, no string, value bool) string {
		if value {
			return yes
		}
		return no
	},
	// Display singluar or plural based on count
	"plural": func(one, many string, count int) string {
		if count == 1 {
			return one
		}
		return many
	},
	// Current Date (Local server time)
	"date": func() string {
		return time.Now().Format("2006-01-02")
	},
	// Current Unix timestamp
	"unixtimestamp": func() string {
		return fmt.Sprintf("%d", time.Now().Unix())
	},
	// json encodes an item into a JSON string
	"json": func(v interface{}) string {
		output, _ := json.Marshal(v)
		return string(output)
	},
	// Allow unsafe injection into HTML
	"noescape": func(a ...interface{}) template.HTML {
		return template.HTML(fmt.Sprint(a...))
	},
	// Allow unsafe URL injections into HTML
	"noescapeurl": func(u string) template.URL {
		return template.URL(u)
	},
	// Modern Hash
	"sha256": func(input string) string {
		hash := sha256.Sum256([]byte(input))
		return hex.EncodeToString(hash[:])
	},
	// Legacy
	"sha1": func(input string) string {
		hash := sha1.Sum([]byte(input))
		return hex.EncodeToString(hash[:])
	},
	// Gravatar
	"md5": func(input string) string {
		hash := md5.Sum([]byte(input))
		return hex.EncodeToString(hash[:])
	},
	// Popular encodings
	"base64encode": func(v string) string {
		return base64.StdEncoding.EncodeToString([]byte(v))
	},
	"base64decode": func(v string) string {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err.Error()
		}
		return string(data)
	},
	"base32encode": func(v string) string {
		return base32.StdEncoding.EncodeToString([]byte(v))
	},
	"base32decode": func(v string) string {
		data, err := base32.StdEncoding.DecodeString(v)
		if err != nil {
			return err.Error()
		}
		return string(data)
	},
}
