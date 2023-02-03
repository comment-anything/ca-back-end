package communication

import "encoding/json"

// GetMessage is a convenience method for getting a message wrapped as a server response for pushing to controller responses
func GetMessage(success bool, text string) ServerResponse {
	return Wrap("Message", text)
}

// Wrap is a convenience method for wrapping some data as a server response
func Wrap(name string, data interface{}) ServerResponse {
	rsp := ServerResponse{}
	rsp.Name = name
	rsp.Data = data
	return rsp
}

func GetErrMsg(success bool, text string) []byte {
	msg := Message{}
	msg.Success = success
	msg.Text = text
	b, _ := json.Marshal(msg)
	return b
}
