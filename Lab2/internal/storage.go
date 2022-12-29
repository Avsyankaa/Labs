package internal

import (
	"Lab2/internal/entities"
	"fmt"
	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
)

type Storage struct {
	Session *dbr.Session
}

func NewStorage() *Storage {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"localhost", "5432", "postgres", "Organization", "123", "disable")

	conn, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		return nil
	}
	conn.SetMaxOpenConns(10)
	return &Storage{
		Session: conn.NewSession(nil),
	}
}

func (s *Storage) AddNewDivision(divisions []entities.Division) error {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil
	}
	defer tx.RollbackUnlessCommitted()

	for _, division := range divisions {
		// Add division
		var divisionId int64
		err = tx.InsertInto("division").
			Pair("name", division.Name).
			Pair("description", division.Description).
			Returning("id").
			Load(&divisionId)
		if err != nil {
			return err
		}

		// Add schedules
		for _, schedule := range division.Schedules {
			_, err = tx.InsertInto("schedule").
				Pair("position_name", schedule.PositionName).
				Pair("position_rate", schedule.PositionRate).
				Pair("division_id", divisionId).
				Exec()
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (s *Storage) ChangeDivision(change *entities.ChangeDivision) error {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil
	}
	defer tx.RollbackUnlessCommitted()
	// Get division id

	var divisionId int64
	err = tx.Select("id").
		From("division").
		Where(dbr.Eq("name", change.Name)).
		LoadOne(&divisionId)
	if err != nil {
		return err
	}
	if change.NewName != "" {
		// Change division name
		_, err = tx.Update("division").
			Set("name", change.NewName).
			Where(dbr.Eq("id", divisionId)).
			Exec()
		if err != nil {
			return err
		}
	}
	if change.NewDescription != "" {
		// Change division description
		_, err = tx.Update("division").
			Set("description", change.NewDescription).
			Where(dbr.Eq("id", divisionId)).
			Exec()
		if err != nil {
			return err
		}
	}
	// Change schedules
	for _, schedule := range change.SchedulesToChange {
		if schedule.NewPositionRate != 0 {
			//Change rate
			_, err = tx.Update("schedule").
				Set("position_rate", schedule.NewPositionRate).
				Where(
					dbr.And(
						dbr.Eq("division_id", divisionId),
						dbr.Eq("position_name", schedule.PositionName),
					),
				).Exec()
			if err != nil {
				return err
			}
		}

		if schedule.NewPositionName != "" {
			// Change name
			_, err = tx.Update("schedule").
				Set("position_name", schedule.NewPositionName).
				Where(
					dbr.And(
						dbr.Eq("division_id", divisionId),
						dbr.Eq("position_name", schedule.PositionName),
					),
				).Exec()
			if err != nil {
				return err
			}
		}
	}

	// Add schedules
	for _, schedule := range change.SchedulesToAdd {
		_, err = tx.InsertInto("schedule").
			Pair("position_name", schedule.PositionName).
			Pair("position_rate", schedule.PositionRate).
			Pair("division_id", divisionId).
			Exec()
		if err != nil {
			return err
		}
	}

	//Delete schedules
	for _, schedule := range change.SchedulesToDelete {
		_, err = tx.DeleteFrom("schedule").
			Where(
				dbr.And(
					dbr.Eq("division_id", divisionId),
					dbr.Eq("position_name", schedule.PositionName),
				),
			).Exec()
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (s *Storage) AddNewStaff(staffs []entities.Staff) error {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil
	}
	defer tx.RollbackUnlessCommitted()

	for _, staff := range staffs {
		var staffId int64
		// Add staff
		err = tx.InsertInto("staff").
			Pair("name", staff.Name).
			Pair("birthday", staff.Birthday).
			Pair("snils", staff.Snils).
			Pair("admission_date", staff.AdmissionDate).
			Returning("staff.id").
			Load(&staffId)
		if err != nil {
			return err
		}

		for _, staffSchedule := range staff.StaffSchedule {
			var scheduleId int64
			err := tx.Select("schedule.id").
				From("schedule").
				Join("division", "schedule.division_id = division.id").
				Where(
					dbr.And(
						dbr.Eq("position_name", staffSchedule.PositionName),
						dbr.Eq("division.name", staffSchedule.DivisionName),
					),
				).
				LoadOne(&scheduleId)
			if err != nil {
				return err
			}
			//Add staff schedule
			_, err = tx.InsertInto("staff_schedule").
				Pair("staff_id", staffId).
				Pair("position_id", scheduleId).
				Exec()
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (s *Storage) GetAllStaff() ([]entities.BriefStaff, error) {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var staff []entities.BriefStaff

	_, err = tx.Select("name, birthday, snils").
		From("staff").
		Load(&staff)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return staff, nil
}

func (s *Storage) GetStaffBy(query dbr.Builder) ([]entities.Staff, error) {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var staffs []entities.Staff

	_, err = tx.Select("id, name, birthday, snils, admission_date").
		From("staff").
		Where(query).
		Load(&staffs)
	if err != nil {
		return nil, err
	}

	for index, _ := range staffs {
		var positionsId []int64
		_, err = tx.Select("position_id").
			From("staff_schedule").
			Where(dbr.Eq("staff_id", staffs[index].Id)).
			Load(&positionsId)
		if err != nil {
			return nil, err
		}

		var positionsNames []string
		for _, position := range positionsId {
			_, err = tx.Select("position_name").
				From("schedule").
				Where(dbr.Eq("id", position)).
				Load(&positionsNames)
			if err != nil {
				return nil, err
			}
		}

		for _, positionName := range positionsNames {
			var divisionName string
			err = tx.Select("name").
				From("division").
				Join("schedule", "schedule.division_id = division.id").
				Where(dbr.Eq("schedule.position_name", positionName)).
				LoadOne(&divisionName)
			if err != nil {
				return nil, err
			}
			staffs[index].StaffSchedule = append(staffs[index].StaffSchedule, entities.StaffSchedule{DivisionName: divisionName, PositionName: positionName})
		}
	}

	tx.Commit()
	return staffs, nil
}

func (s *Storage) GetStaffByName(name string) ([]entities.Staff, error) {
	query := dbr.Eq("name", name)
	return s.GetStaffBy(query)
}

func (s *Storage) GetStaffBySnils(snils string) ([]entities.Staff, error) {
	query := dbr.Eq("snils", snils)
	return s.GetStaffBy(query)
}

func (s *Storage) GetAllDivisionNames() ([]string, error) {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var names []string

	_, err = tx.Select("name").
		From("division").
		Load(&names)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return names, nil
}

func (s *Storage) DeleteDivision(name string) error {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.DeleteFrom("division").
		Where(dbr.Eq("name", name)).Exec()
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (s *Storage) DeleteStaff(snils string) error {
	tx, err := s.Session.Begin()
	if err != nil {
		return nil
	}
	defer tx.RollbackUnlessCommitted()

	var staffId int64
	err = tx.Select("id").
		From("staff").
		Where(dbr.Eq("staff.snils", snils)).
		LoadOne(&staffId)
	if err != nil {
		return err
	}
	_, err = tx.DeleteFrom("staff_schedule").
		Where(dbr.Eq("staff_id", staffId)).
		Exec()
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom("staff").
		Where(dbr.Eq("snils", snils)).
		Exec()
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
