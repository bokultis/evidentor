package student

import (
	"database/sql"
	"encoding/json"
	"evidentor/api/apiutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func StudentsIndexHandler(w http.ResponseWriter, r *http.Request) {

	var students []*StudentDO

	students, err := GetAllStudents()
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	studentsWO := make([]interface{}, len(students))
	for i, student := range students {

		studentsWO[i] = makeStudentWO(student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentsWO)
}

func StudentsShowHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentId, _ := strconv.Atoi(params["studentId"])
	student, err := GetStudent(studentId)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeStudentWO(student))
}

func StudentsGroupsListHandler(w http.ResponseWriter, r *http.Request) {
	var groups []*GroupDO
	params := mux.Vars(r)
	studentId, _ := strconv.Atoi(params["studentId"])
	groups, err := GetStudentGroups(studentId)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	groupsWO := make([]interface{}, len(groups))
	for i, group := range groups {

		groupsWO[i] = makeGroupWO(group)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupsWO)
}

func StudentsNotesListHandler(w http.ResponseWriter, r *http.Request) {
	var notes []*StudentNoteDO
	params := mux.Vars(r)
	studentId, _ := strconv.Atoi(params["studentId"])
	notes, err := GetStudentNotes(studentId)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	notesWO := make([]interface{}, len(notes))
	for i, note := range notes {

		notesWO[i] = makeStudentNoteWO(note)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notesWO)
}

// func StudentsAttendanceListHandler(w http.ResponseWriter, r *http.Request) {
// 	var attendances []*AttendanceDO
// 	params := mux.Vars(r)
// 	studentId, _ := strconv.Atoi(params["studentId"])
// 	attendances, err := GetStudentAttendance(studentId, filter)
// 	if err != nil {
// 		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
// 		return
// 	}

// 	attendancesWO := make([]interface{}, len(attendances))
// 	for i, attendance := range attendances {

// 		attendancesWO[i] = makeAttendanceWO(attendance)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(attendancesWO)
// }

func StudentsCreateHandler(w http.ResponseWriter, r *http.Request) {
	var student StudentInputWO
	var firstName, lastName sql.NullString
	var birthday mysql.NullTime

	err := apiutil.GetRequestBody(r, &student)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
	}

	validationErr := student.validate()
	if validationErr != nil {
		apiutil.NewErrorResponse(w, validationErr)
		return
	}

	if student.FirstName == nil {
		firstName = sql.NullString{Valid: false, String: ""}
	} else {
		firstName = sql.NullString{Valid: true, String: *student.FirstName}
	}

	if student.LastName == nil {
		lastName = sql.NullString{Valid: false, String: ""}
	} else {
		lastName = sql.NullString{Valid: true, String: *student.LastName}
	}

	if student.Birthday == nil {
		birthday = mysql.NullTime{Valid: false, Time: time.Now()}
	} else {
		bt, _ := time.Parse("2006-01-02", *student.Birthday)
		birthday = mysql.NullTime{Valid: true, Time: bt}
	}
	newStudent := &StudentDO{
		FirstName: firstName,
		LastName:  lastName,
		Gender:    *student.Gender,
		Birthday:  birthday,
	}

	studentID, err := CreateStudent(newStudent)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	studentCreated, err := GetStudent(studentID)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeStudentWO(studentCreated))
}

func StudentsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentId, _ := strconv.Atoi(params["studentId"])

	err := DeleteStudent(studentId)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Student deleted")
}

func StudentsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student StudentUpdateWO
	studentId, _ := strconv.Atoi(params["studentId"])

	err := apiutil.GetRequestBody(r, &student)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
	}

	update, validationErr := student.validate()
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(validationErr))
		return
	}

	err = UpdateStudent(studentId, *update)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	updatedStudent, err := GetStudent(studentId)
	if err != nil {
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeStudentWO(updatedStudent))
}
