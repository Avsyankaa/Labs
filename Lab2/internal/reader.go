package internal

import (
	"Lab2/internal/entities"
	"encoding/json"
	"io/ioutil"
)

type Reader struct {
}

func (r *Reader) ReadStaffs(filename string) ([]entities.Staff, error) {
	staff, _ := ioutil.ReadFile(filename)
	var data []entities.Staff
	err := json.Unmarshal(staff, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (r *Reader) ReadDivisions(filename string) ([]entities.Division, error) {
	division, _ := ioutil.ReadFile(filename)
	var data []entities.Division
	err := json.Unmarshal(division, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (r *Reader) ReadChangeDivision(filename string) (*entities.ChangeDivision, error) {
	division, _ := ioutil.ReadFile(filename)
	var data entities.ChangeDivision
	err := json.Unmarshal(division, &data)
	if err != nil {
		return nil, err
	}
	return &data, err
}
