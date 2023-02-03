package communication

import "encoding/json"

// MessageBytes is a convenience method for getting a message wrapped as a server response for pushing to controller responses
func MessageBytes(success bool, text string) []byte {
	msg := &Message{}
	msg.Success = success
	msg.Text = text
	bytes, _ := json.Marshal(msg)
	return Wrap("Message", bytes)
}

func Wrap(name string, data []byte) []byte {
	rsp := &ServerResponse{}
	rsp.name = name
	rsp.data = data
	bytes, _ := json.Marshal(rsp)
	return bytes
}
