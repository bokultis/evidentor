package student

import (
	"database/sql"
	"evidentor/api/db"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type StudentDO struct {
	ID        int            `gorm:"column:id;primary_key"`
	FirstName sql.NullString `gorm:"column:first_name"`
	LastName  sql.NullString `gorm:"column:last_name"`
	Gender    string         `gorm:"column:gender" sql:"default:unknown"`
	Birthday  mysql.NullTime `gorm:"column:birthday" sql:"default:NULL"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

type GroupDO struct {
	ID        int            `gorm:"column:id;primary_key"`
	Name      sql.NullString `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

type StudentNoteDO struct {
	ID        int            `gorm:"column:id;primary_key"`
	StudentID int            `gorm:"column:student_id"`
	Note      sql.NullString `gorm:"column:note"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

type AttendanceDO struct {
	ID        int            `gorm:"column:id;primary_key"`
	StudentID int            `gorm:"column:student_id"`
	Date      mysql.NullTime `gorm:"column:date"`
	Status    string         `gorm:"column:status"`
	Remark    string         `gorm:"column:remark"`
	Note      sql.NullString `gorm:"column:note"`
}

func CreateStudent(student *StudentDO) (int, error) {

	database, err := db.GetDbName()
	if err != nil {
		return 0, err
	}

	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	tx := db.DB.Begin()

	e := tx.Table(database + ".students").Create(student).Error
	if e != nil {
		tx.Rollback()
		return 0, e
	}

	tx.Commit()
	return student.ID, nil
}

func GetStudent(id int) (*StudentDO, error) {
	var student StudentDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}
	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	if err := db.DB.Table(database+".students ").
		Select("students.*").
		Where("students.id = ?", id).
		First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &student, nil
}

func GetAllStudents() ([]*StudentDO, error) {
	var students []*StudentDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}

	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	if err := db.DB.Table(database + ".students ").
		Select("students.*").
		Find(&students).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return students, nil
}

func GetStudentGroups(studentId int) ([]*GroupDO, error) {
	var groups []*GroupDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}

	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	// SELECT * FROM groups as g
	// 	LEFT JOIN student_groups as sg on sg.group_id = g.id
	// 	WHERE sg.student_id = 10

	if err := db.DB.Table(database+".groups ").
		Select("groups.*").
		Joins("left join student_groups on student_groups.group_id = groups.id").
		Where("student_groups.student_id = ?", studentId).
		Find(&groups).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return groups, nil
}

func GetStudentNotes(studentId int) ([]*StudentNoteDO, error) {
	var notes []*StudentNoteDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}

	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	if err := db.DB.Table(database+".student_notes ").
		Select("student_notes.*").
		Where("student_notes.student_id = ?", studentId).
		Find(&notes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return notes, nil
}

func UpdateStudent(studentID int, update map[string]interface{}) error {

	database, err := db.GetDbName()
	if err != nil {
		return err
	}

	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	tx := db.DB.Begin()

	var studentOld StudentDO

	e1 := tx.Table(database+".students").First(&studentOld, studentID).Error
	if e1 != nil {
		if e1 == gorm.ErrRecordNotFound {
			tx.Rollback()
			return e1
		}
		tx.Rollback()
		return e1
	}

	e := tx.Table(database + ".students").Model(&studentOld).Updates(update).Error
	if e != nil {
		tx.Rollback()
		return e
	}

	tx.Commit()
	return nil
}

// Delete
func DeleteStudent(id int) error {
	var student StudentDO

	database, err := db.GetDbName()
	if err != nil {
		return err
	}

	db.DB.LogMode(true)
	defer db.DB.LogMode(false)

	err = db.DB.Table(database+".students").Where("id=?", id).Delete(&student).Error
	if err != nil {
		return err
	}
	return nil
}
