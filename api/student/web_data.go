package student

type StudentWO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type StudentInputWO struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Gender    *string `json:"gender"`
	Birthday  *string `json:"birthday"`
}

type StudentUpdateWO struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Gender    *string `json:"gender"`
	Birthday  *string `json:"birthday"`
}

type GroupWO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type StudentNoteWO struct {
	ID        int    `json:"id"`
	StudentID int    `json:"student_id"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AttendanceWO struct {
	ID        int    `json:"id"`
	StudentID int    `json:"student_id"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Remark    string `json:"remark"`
	Note      string `json:"note"`
}

func makeStudentWO(student *StudentDO) *StudentWO {

	firstName := ""
	if student.FirstName.Valid {
		firstName = student.FirstName.String
	}

	lastName := ""
	if student.LastName.Valid {
		lastName = student.LastName.String
	}

	birthday := ""
	if student.Birthday.Valid {
		birthday = student.Birthday.Time.Format("2006-01-02")
	}

	return &StudentWO{
		ID:        student.ID,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    student.Gender,
		Birthday:  birthday,
		CreatedAt: student.CreatedAt.String(),
		UpdatedAt: student.UpdatedAt.String(),
	}

}

func makeGroupWO(group *GroupDO) *GroupWO {

	name := ""
	if group.Name.Valid {
		name = group.Name.String
	}

	return &GroupWO{
		ID:        group.ID,
		Name:      name,
		CreatedAt: group.CreatedAt.String(),
		UpdatedAt: group.UpdatedAt.String(),
	}

}

func makeStudentNoteWO(note *StudentNoteDO) *StudentNoteWO {

	nt := ""
	if note.Note.Valid {
		nt = note.Note.String
	}

	return &StudentNoteWO{
		ID:        note.ID,
		StudentID: note.StudentID,
		Note:      nt,
		CreatedAt: note.CreatedAt.String(),
		UpdatedAt: note.UpdatedAt.String(),
	}

}

func makeAttendanceWO(attendance *AttendanceDO) *AttendanceWO {

	note := ""
	if attendance.Note.Valid {
		note = attendance.Note.String
	}

	date := ""
	if attendance.Date.Valid {
		date = attendance.Date.Time.Format("2006-01-02")
	}

	return &AttendanceWO{
		ID:        attendance.ID,
		StudentID: attendance.StudentID,
		Date:      date,
		Note:      note,
		Remark:    attendance.Remark,
		Status:    attendance.Status,
	}

}
