package entities

type Staff struct {
	Id            int64           `db:"id"`
	Name          string          `json:"name" db:"name"`
	Birthday      string          `json:"birthday" db:"birthday"`
	Snils         string          `json:"snils" db:"snils"`
	AdmissionDate string          `json:"admission_date" db:"admission_date"`
	StaffSchedule []StaffSchedule `json:"staff_schedule"`
}

type BriefStaff struct {
	Name     string `json:"name" db:"name"`
	Birthday string `json:"birthday" db:"birthday"`
	Snils    string `json:"snils" db:"snils"`
}

type Division struct {
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Schedules   []Schedule `json:"schedules"`
}

type StaffSchedule struct {
	PositionName string `json:"position_name" db:"position_name"`
	DivisionName string `json:"division_name" db:"division_name"`
}

type Schedule struct {
	PositionName string  `json:"position_name" db:"position_name"`
	PositionRate float32 `json:"position_rate" db:"position_rate"`
}

type ChangeSchedule struct {
	PositionName    string  `json:"position_name" db:"position_name"`
	NewPositionName string  `json:"new_position_name" db:"new_position_name"`
	NewPositionRate float32 `json:"new_position_rate" db:"new_position_rate"`
}

type ChangeDivision struct {
	Name              string           `json:"name" db:"name"`
	NewName           string           `json:"new_name"`
	NewDescription    string           `json:"new_description"`
	SchedulesToAdd    []Schedule       `json:"add_schedules"`
	SchedulesToDelete []Schedule       `json:"delete_schedules"`
	SchedulesToChange []ChangeSchedule `json:"change_schedules"`
}
