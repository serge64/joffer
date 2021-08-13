package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/valyala/fastjson"
	"guthub.com/serge64/joffer/internal/config"
)

type Profile struct {
	ID           int    `json:"-" db:"id"`
	UserID       int    `json:"-" db:"user_id"`
	PlatformID   int    `json:"-" db:"platform_id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	AccessToken  string `json:"-" db:"access_token"`
	RefreshToken string `json:"-" db:"refresh_token"`
	Expiry       int    `json:"-" db:"expiry"`
}

type response struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (p *Profile) Resumes() ([]Resume, error) {
	req, err := http.NewRequest("GET", "https://api.hh.ru/resumes/mine", nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+p.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	json := func() []byte {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(res.Body); err != nil {
			log.Print(err)
		}
		return buf.Bytes()
	}()

	parse := fastjson.Parser{}
	value, err := parse.ParseBytes(json)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	resumes := []Resume{}

	for _, o := range value.Get("items").GetArray() {
		resume := Resume{
			ProfileID: p.ID,
			Name:      string(o.GetStringBytes("title")),
			UID:       string(o.GetStringBytes("id")),
		}
		resumes = append(resumes, resume)
	}

	return resumes, nil
}

func (p *Profile) Me() (*Profile, error) {
	req, err := http.NewRequest("GET", "https://api.hh.ru/me", nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+p.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	json := func() []byte {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(res.Body); err != nil {
			log.Print(err)
		}
		return buf.Bytes()
	}()

	parse := fastjson.Parser{}
	value, err := parse.ParseBytes(json)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	name := []string{}
	name = append(name, value.Get("first_name").String())
	name = append(name, value.Get("last_name").String())
	name = append(name, value.Get("middle_name").String())

	fullName := strings.ReplaceAll(strings.Join(name, " "), "\"", "")
	fullName = strings.ReplaceAll(fullName, "null", "")
	email := strings.ReplaceAll(value.Get("email").String(), "\"", "")

	p.Name = fullName
	p.Email = email

	return p, nil
}

func (p *Profile) Authorization(code string, c *config.Config) error {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"code":          {code},
	}

	err := p.postRequest(data.Encode())
	if err != nil {
		return err
	}

	return nil
}

func (p *Profile) UpdateToken() error {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {p.RefreshToken},
	}

	err := p.postRequest(data.Encode())
	if err != nil {
		return err
	}

	return nil
}

func (p *Profile) postRequest(data string) error {
	r, err := http.Post(
		"https://hh.ru/oauth/token",
		"application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(data)),
	)
	if err != nil {
		return err
	}

	fmt.Println(r.StatusCode)

	defer r.Body.Close()
	resp := &response{}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}

	p.AccessToken = resp.AccessToken
	p.RefreshToken = resp.RefreshToken
	p.Expiry = resp.ExpiresIn

	return nil
}
