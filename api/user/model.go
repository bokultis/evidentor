package user

import (
	"database/sql"
	"evidentor/api/db"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserDO struct {
	ID        int            `gorm:"column:id;primary_key"`
	FirstName sql.NullString `gorm:"column:first_name"`
	LastName  sql.NullString `gorm:"column:last_name"`
	Gender    string         `gorm:"column:gender" sql:"default:unknown"`
	Birthday  mysql.NullTime `gorm:"column:birthday" sql:"default:NULL"`
	Email     string         `gorm:"column:email"`
	Address   sql.NullString `gorm:"column:address"`
	Optin     string         `gorm:"column:optin" sql:"default:unknown"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

type RoleDO struct {
	UserID int    `gorm:"column:user_id"`
	RoleID int    `gorm:"column:role_id"`
	Name   string `gorm:"column:name"`
}

type JWTToken struct {
	Token string `json:"token"`
}

func CreateUser(user *UserDO) (int, error) {

	database, err := db.GetDbName()
	if err != nil {
		return 0, err
	}

	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	tx := db.DB.Begin()

	e := tx.Table(database + ".users").Create(user).Error
	if e != nil {
		tx.Rollback()
		return 0, e
	}

	tx.Commit()
	return user.ID, nil
}

func GetUser(id int) (*UserDO, error) {
	var user UserDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}
	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	if err := db.DB.Table(database+".users").
		Select("users.*").
		Where("users.id=? ", id).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*UserDO, error) {
	var user UserDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}
	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	if err := db.DB.Table(database+".users").
		Select("users.*").
		Where("users.email=? ", email).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func GetUserRole(id int) (*RoleDO, error) {
	var role RoleDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}
	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	if err := db.DB.Table(database+".user_roles").
		Select("user_roles.user_id, user_roles.role_id, role.name").
		Joins("LEFT JOIN "+database+".role  on role.id=user_roles.role_id").
		Where("user_roles.user_id=? ", id).
		Find(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &role, nil
}

func GetAllUsers() ([]*UserDO, error) {
	var users []*UserDO

	database, err := db.GetDbName()
	if err != nil {
		return nil, err
	}

	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	if err := db.DB.Table(database + ".users ").
		Select("users.*").
		Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return users, nil
}

func UpdateUser(userID int, update map[string]interface{}) error {

	database, err := db.GetDbName()
	if err != nil {
		return err
	}

	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	tx := db.DB.Begin()

	var userOld UserDO

	e1 := tx.Table(database+".users").First(&userOld, userID).Error
	if e1 != nil {
		if e1 == gorm.ErrRecordNotFound {
			tx.Rollback()
			return e1
		}
		tx.Rollback()
		return e1
	}

	e := tx.Table(database + ".users").Model(&userOld).Updates(update).Error
	if e != nil {
		tx.Rollback()
		return e
	}

	tx.Commit()
	return nil
}

// Delete
func DeleteUser(id int) error {
	var user UserDO

	database, err := db.GetDbName()
	if err != nil {
		return err
	}

	// db.DB.LogMode(true)
	// defer db.DB.LogMode(false)

	err = db.DB.Table(database+".users").Where("id=?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func makeUserWO(usr *UserDO) *UserWO {

	birthday := ""
	if usr.Birthday.Valid {
		birthday = usr.Birthday.Time.Format("2006-01-02")
	}
	return &UserWO{
		ID:        usr.ID,
		FirstName: usr.FirstName.String,
		LastName:  usr.LastName.String,
		Gender:    usr.Gender,
		Birthday:  birthday,
		Email:     usr.Email,
		Address:   usr.Address.String,
		Optin:     usr.Optin,
	}

}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes)
}

func (u UserDO) checkPassword(password string) bool {
	//err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	//return err == nil
	return true
}

func (u UserWO) generateJWT() (JWTToken, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(time.Hour * 24 * 1).Unix(),
		"user_id":    int(u.ID),
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
	})
	tokenString, err := token.SignedString(signingKey)
	return JWTToken{tokenString}, err
}
