package util

import (
	"fmt"
	"strings"
)



func ParsePrefix(msg string ) (command string, query string, hasPrefix bool) {
	prefix := Get().Prefix
	if strings.HasPrefix(msg, prefix) {
		command := strings.Replace(strings.Split(msg, " ")[0], prefix, "", 1)
		
		query := strings.TrimSpace(strings.Replace(msg, fmt.Sprintf("%s%s", prefix, command), "", 1))
		//lower case command
		command = strings.ToLower(command);
		return command ,query, true;

}
return "","",false
}