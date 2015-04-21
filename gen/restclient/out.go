package restclient

import (
	"bytes"
	"net/http"
)

type IUserServiceClient struct {
	service          core.IUserService
	RequestDecorator func(req *http.Request) (*http.Request, error)
}

func (svc *IUserServiceClient) SendRemoveAllUsersRequest(arg0 *core.Request, arg1 time.Time) (*http.Response, error) {
	var body *bytes.Buffer = bytes.NewBuffer(nil)
	body.Write([]byte("["))
	Write_Request(body, arg0)
	UnknownWriter(body, arg1)
	body.Write([]byte("]"))

	httpreq, err := http.NewRequest("GET", "http://hello.world/", body)
	if err != nil {
		return nil, err
	}
	httpreq.Header.Add("Content-Type", "application/json")
	if svc.RequestDecorator != nil {
		httpreq, err = svc.RequestDecorator(httpreq)
		if err != nil {
			return nil, err
		}
	}
	c := http.Client{}
	return c.Do(httpreq)
}

func (svc *IUserServiceClient) SendGetUserRequest(arg0 *core.User) (*http.Response, error) {
	var body *bytes.Buffer = bytes.NewBuffer(nil)
	Write_User(body, arg0)

	httpreq, err := http.NewRequest("GET", "http://hello.world/", body)
	if err != nil {
		return nil, err
	}
	httpreq.Header.Add("Content-Type", "application/json")
	if svc.RequestDecorator != nil {
		httpreq, err = svc.RequestDecorator(httpreq)
		if err != nil {
			return nil, err
		}
	}
	c := http.Client{}
	return c.Do(httpreq)
}

func (svc *IUserServiceClient) SendSaveUserRequest(arg0 *core.SaveUserRequest) (*http.Response, error) {
	var body *bytes.Buffer = bytes.NewBuffer(nil)
	Write_SaveUserRequest(body, arg0)

	httpreq, err := http.NewRequest("GET", "http://hello.world/", body)
	if err != nil {
		return nil, err
	}
	httpreq.Header.Add("Content-Type", "application/json")
	if svc.RequestDecorator != nil {
		httpreq, err = svc.RequestDecorator(httpreq)
		if err != nil {
			return nil, err
		}
	}
	c := http.Client{}
	return c.Do(httpreq)
}
