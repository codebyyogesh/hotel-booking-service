package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)
const(
    bcryptCost      = 12
    minFirstNameLen = 2
    minLastNameLen  = 2
    minPasswordLen  = 7
)
type UpdateUserParams struct{
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
}

type CreateUserParams struct{
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}

type User struct{
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    FirstName string `bson:"firstName" json:"firstName"`
    LastName  string `bson:"lastName" json:"lastName"`
    Email     string `bson:"email" json:"email"`
    EncryptedPassword string `bson:"EncryptedPassword" json:"-"`
}

// TODO : Add errors for this function
func (params UpdateUserParams) ValidateUpdateUserParams() bool{

    if len(params.FirstName) < minFirstNameLen {
        return false
    }
    if len(params.LastName) < minLastNameLen {
        return false
    }
    return true
}

// TODO: handle errors in this function
func (params UpdateUserParams) ToBsonD() bson.D {
    d := bson.D{}
    if params.ValidateUpdateUserParams() {
        d = append(d, bson.E{Key: "firstName", Value: params.FirstName})
    }
    if params.ValidateUpdateUserParams() {
        d = append(d, bson.E{Key: "lastName", Value: params.LastName})
    }
    return d
}

func (params CreateUserParams)ValidateUserParams() map[string]string{
    errors := map[string]string{}
    if len(params.FirstName) < minFirstNameLen{
        errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
    }
    if len(params.LastName) < minLastNameLen {
        errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
    }
    if len(params.Password) < minPasswordLen {
        errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
    }
    if !IsEmailValid(params.Email) {
        errors["email"] =  "invalid email"
    }
    return errors
}

func IsPasswordValid(encpw, userpw string) bool{
    return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(userpw)) == nil
}

func IsEmailValid(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}

func NewUserFromParams(params CreateUserParams) (*User, error){
    encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

    if err != nil{
        return nil, err
    }
    return &User{
        FirstName:         params.FirstName,
        LastName:          params.LastName,
        Email:             params.Email,
        EncryptedPassword: string(encpw),
    }, nil
}
