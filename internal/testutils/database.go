package testutils

import (
	"fmt"
	"math"
	"math/big"
	"time"

	crand "crypto/rand"
	"math/rand"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID            uint    `gorm:"column:ID,primaryKey"`
	CAGId         int64   `gorm:"column:CAGId"`
	Domain        uint    `gorm:"column:Domain"`
	Username      string  `gorm:"column:Username"`
	Flags         uint    `gorm:"column:Flags"`
	MailboxLimit  *uint   `gorm:"column:MailboxLimit"`
	Forward       *string `gorm:"column:Forward"`
	RealName      *string `gorm:"column:RealName"`
	Comment       *string `gorm:"column:Comment"`
	Storage       string  `gorm:"column:Storage"`
	Target        string  `gorm:"column:Target"`
	PwdCode       uint    `gorm:"column:PwdCode"`
	AttachStorage string  `gorm:"column:AttachStorage"`
	Flags2        uint    `gorm:"column:Flags2"`
}

type SecurityImage struct {
	ID   int        `gorm:"column:id,primaryKey"`
	Word *string    `gorm:"column:word"`
	Time *time.Time `gorm:"column:time"`
	IP   *uint      `gorm:"column:ip"`
}

func (User) TableName() string {
	return "user"
}

func (SecurityImage) TableName() string {
	return "security_image"
}

func ptr[T any](v T) *T {
	return &v
}

func Shuffle[T any](src []T) []T {
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(src))

	dest := make([]T, len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

func randStr() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	res := make([]byte, 32)
	crand.Read(res)

	for i, b := range res {
		res[i] = letterBytes[int(b)%len(letterBytes)]
	}
	return string(res)
}

func randInt64(max int64) int64 {
	v, err := crand.Int(crand.Reader, big.NewInt(max))
	if err != nil {
		panic(fmt.Sprintf("rang returned error: %s", err))
	}

	return v.Int64()
}

func randInt() uint {
	return uint(randInt64(math.MaxInt32))
}

func randTime() time.Time {
	return time.Unix(randInt64(time.Now().Unix()), 0)
}

func generate[T any](n int, gen func() T) []T {
	res := make([]T, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, gen())
	}
	return res
}

func NewUser() *User {
	return &User{
		// ID should be automatically set after insert
		CAGId:         randInt64(5000000000),
		Domain:        randInt(),
		Username:      randStr(),
		Flags:         randInt(),
		MailboxLimit:  ptr(randInt()),
		Forward:       ptr(randStr()),
		RealName:      ptr(randStr()),
		Comment:       ptr(randStr()),
		Storage:       randStr(),
		Target:        randStr(),
		PwdCode:       randInt(),
		AttachStorage: randStr(),
		Flags2:        randInt(),
	}
}

func NewUsers(n int) []*User {
	return generate(n, NewUser)
}

func NewSecurityImage() *SecurityImage {
	return &SecurityImage{
		// ID should be automatically set after insert
		Word: ptr(randStr()),
		Time: ptr(randTime()),
		IP:   ptr(randInt()),
	}
}

func NewSecurityImages(n int) []*SecurityImage {
	return generate(n, NewSecurityImage)
}

func Gorm() (*gorm.DB, func()) {
	g, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("unable to initialize gorm: %s", err))
	}

	return g, func() {
		db, err := g.DB()
		if err != nil {
			panic(fmt.Sprintf("unablt to get Gorm DB: %s", err))
		}

		db.Close()
	}
}
