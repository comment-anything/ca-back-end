package communication

import (
	"encoding/json"
)

// GetMessage is a convenience method for getting a communication.Message wrapped as a communication.ServerResponse for pushing to controller responses.
func GetMessage(success bool, text string) ServerResponse {
	msg := Message{}
	msg.Success = success
	msg.Text = text
	return Wrap("Message", msg)
}

// Wrap is a convenience method for wrapping some data as a communicaiton.ServerResponse.
func Wrap(name string, data interface{}) ServerResponse {
	rsp := ServerResponse{}
	rsp.Name = name
	rsp.Data = data
	return rsp
}

// GetErrrMsg returns a byte array for sending directly to the front end, which the front end can decode, as JSON, to an array of exactly 1 ServerResponse with 1 Message inside, as provided. This should **not** normally be used; a general error message is sent by calling GetMessage and adding that to the Controller.nextResponse value, and allowing the Controller to handle sending all the responses at the end of the Request pipeline. GetErrMsg is only for situations where *there is no controller and you have to directly Write the HTTP.Response*.
func GetErrMsg(success bool, text string) []byte {
	sr := ServerResponse{}
	msg := Message{}
	msg.Success = success
	msg.Text = text
	sr.Name = "Message"
	sr.Data = msg
	var b []ServerResponse
	b = append(b, sr)
	ret_val, _ := json.Marshal(b)
	return ret_val
}
