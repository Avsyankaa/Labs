package main

import (
	"Lab2/internal"
	"fmt"
)

func main() {
	storage := internal.NewStorage()
	err := storage.DeleteStaff("234-134-23-14")
	if err != nil {
		fmt.Println(err)
		return
	}
	/*var reader internal.Reader
	divisions, err := reader.ReadDivisions("division.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	storage := internal.NewStorage()
	err = storage.AddNewDivision(divisions)
	if err != nil {
		fmt.Println(err)
		return
	}

	staff, err := reader.ReadStaffs("staff.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = storage.AddNewStaff(staff)
	if err != nil {
		fmt.Println(err)
		return
	}

		changeDivisions, err := reader.ReadChangeDivision("change_division.txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = storage.ChangeDivision(changeDivisions)
		if err != nil {
			fmt.Println(err)
			return
		}

		staff, err := storage.GetAllDivisionNames()
		if err != nil {
			fmt.Println(err)
			return
		}

		/*var reader internal.Reader
		divisions, err := reader.ReadDivisions("division.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		storage := internal.NewStorage()
		err = storage.AddNewDivision(divisions)
		if err != nil {
			fmt.Println(err)
			return
		}

			staff, err := storage.GetAllStaff()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(staff)

			storage := internal.NewStorage()
			//storage.DeleteStaff("234-134-23-14")

			/*staffs, err := storage.GetStaffBySnils("234-134-23-14")
			if err != nil {
				fmt.Println(err)
				return
			}
			res, err := json.Marshal(staffs)
			fmt.Println(string(res))


			staff, err := reader.ReadStaffs("staff.txt")
			if err != nil {
				fmt.Println(err)
				return
			}
			err = storage.AddNewStaff(staff)
			if err != nil {
				fmt.Println(err)
				return
			}*/

	/*names, err := storage.GetAllDivisionNames()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(names)

	changeDivisions, err := reader.ReadChangeDivision("change_division.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = storage.ChangeDivision(changeDivisions)
	if err != nil {
		fmt.Println(err)
		return
	}

	err := storage.DeleteDivision("Бухгалтерия")
	if err != nil {
		fmt.Println(err)
		return
	}*/

}
