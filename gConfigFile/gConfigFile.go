package gConfigFile

// use json as config file

import (
	"os"
	"encoding/json"
	"fmt"
)

type configf struct {
	Mainurl string
	Header  string
	Keyword string
	RountineNum int
}

func getElement(ele string) (interface{}, error) {
	f, err := os.OpenFile("config.json", os.O_RDONLY, os.ModePerm)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("open file confi.json error: %v", err)
	} else if err != nil {
		f, err = os.Create("config.json")
		if err != nil {
			return nil, fmt.Errorf("create file config.json error: %v", err)
		} else { // write json body
			//err = json.NewEncoder(f).Encode(&configf{mainurl:"", header:"", keyword:"", rountineNum:0})
			conf := []configf{
				{Mainurl:"", Header:"", Keyword:"", RountineNum:0}}
			data, err := json.MarshalIndent(conf, "", "		")
			if err != nil {
				f.Close()
				return nil, fmt.Errorf("error write config.json: %v", err)
			} else {
				f.Write(data)
				f.Close()
				return nil, fmt.Errorf("please config file config.json")
			}
		}
	} else {
		var config []configf
		err = json.NewDecoder(f).Decode(&config)
		configOne := config[0]
		if err != nil {
			return nil, fmt.Errorf("get element error: %v", err)
		}

		defer f.Close()
		// get element
		if ele == "mainurl" {
			return configOne.Mainurl, nil
		}
		if ele == "header" {
			return configOne.Header, nil
		}
		if ele == "keyword" {
			return configOne.Keyword, nil
		}
		if ele == "rountineNum" {
			return configOne.RountineNum, nil
		}
		return nil, nil
	}
}

func GetMainUrl() (string, error) {
	s, err := getElement("mainurl")
	if err != nil {
		return "", err
	}
	return s.(string), err
}

func GetHeader() (string, error) {
	s, err := getElement("header")
	if err != nil {
		return "", nil
	}
	return s.(string), err
}

func GetKeyWord() (string, error) {
	s, err := getElement("keyword")
	if err != nil {
		return "", nil
	}
	return s.(string), err
}

func GetRoutineNum() (int, error) {
	s, err := getElement("rountineNum")
	if err != nil {
		return 0, err
	}
	return s.(int), err
}
