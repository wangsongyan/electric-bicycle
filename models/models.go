package models

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var location, _ = time.LoadLocation("Asia/Shanghai")
var engine *xorm.Engine

type Date time.Time

func (d *Date) UnmarshalText(text []byte) (err error) {
	t, err := time.ParseInLocation("2006-01-02", string(text), location)
	//log.Println(t)
	*d = Date(t)
	return
}

func (d *Date) MarshalJSON() ([]byte, error) {
	format := fmt.Sprintf(`"%s"`, time.Time(*d).Format("2006-01-02"))
	return []byte(format), nil
}

type Number int

func (n *Number) UnmarshalText(text []byte) (err error) {
	number, _ := decimal.NewFromString(string(text))
	*n = Number(int(number.Mul(decimal.NewFromFloat(100)).IntPart()))
	return
}

type Distance struct {
	Id        uint      `xorm:"pk autoincr" json:"id" form:"id"`
	Number    Number    `json:"number" form:"number,required"`
	Mileage   int       `json:"mileage"`
	StartAt   Date      `json:"startAt" form:"startAt,required"`
	CreatedAt time.Time `xorm:"created" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated" json:"updatedAt"`
}

func InitDB() (e *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("mysql", "root:mysql@tcp(127.0.0.1:3307)/test?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		return
	}
	err = engine.Ping()
	if err != nil {
		return
	}
	engine.ShowSQL(true)
	engine.TZLocation = location
	e = engine
	err = engine.Sync2(new(Distance))
	return
}

func (d *Distance) Insert() (err error) {
	_, err = engine.Insert(d)
	return
}

func (d *Distance) Delete() (err error) {
	_, err = engine.Delete(d)
	return
}

func DistanceList() (distances []Distance, err error) {
	err = engine.Find(&distances)
	return
}

func UpdateMileage() (err error) {
	var distances []Distance
	err = engine.Asc("start_at").Find(&distances)
	if err != nil {
		return
	}
	for i, distance := range distances {
		if i == 0 {
			continue
		}
		distance.Mileage = int(distance.Number) - int(distances[i-1].Number)
		engine.Id(distance.Id).Update(distance)
	}
	return
}

func GetLastDistance() (distance Distance, err error) {
	has, err := engine.Table(new(Distance)).Desc("start_at").Get(&distance)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("cannot find record")
	}
	return
}
