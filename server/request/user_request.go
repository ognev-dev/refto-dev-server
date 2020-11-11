package request

import "errors"

type LoginWithGithub struct {
	Code string `json:"code"`
}

func (r *LoginWithGithub) Validate() (err error) {
	if r.Code == "" {
		err = errors.New("login code missing")
		return
	}

	return
}
