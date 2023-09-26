package validator

import "regexp"

var (
	emailRegex = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)
	phoneRegex = regexp.MustCompile(`^((\+86)|(86))?1([356789][0-9]|4[579]|6[67]|7[0135678]|9[189])[0-9]{8}$`)
	ipRegex    = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
)

const (
	validTag = "valid"
)
