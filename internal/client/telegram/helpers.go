package telegram

import "fmt"

func getCommandsString(commands map[string]CommandHandler) string {
	res := ""
	i := 0
	for k, _ := range commands {
		i++
		if i < len(commands) {
			res += fmt.Sprintf("/%s, ", k)
			continue
		}
		res += fmt.Sprintf("/%s", k)
	}
	return res
}
