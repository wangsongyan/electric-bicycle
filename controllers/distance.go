package controllers

import (
	"electric-bicycle/models"
	"github.com/kataras/iris/v12"
)

//decoder := schema.NewDecoder()
//decoder.SetAliasTag("form")
//decoder.RegisterConverter(time.Now(), func(value string) reflect.Value {
//	if v, err := time.Parse("2006-01-02", value); err == nil {
//		return reflect.ValueOf(v)
//	}
//	return reflect.Value{}
//})
//var distance models.Distance
//values := ctx.FormValues()
//if len(values) != 0 {
//	err = decoder.Decode(&distance, values)
//}

func DistanceNew(ctx iris.Context) {
	var (
		distance models.Distance
		ret      = iris.Map{}
		err      error
	)
	defer writeResponse(ctx, ret)
	err = ctx.ReadForm(&distance)
	if err != nil && !iris.IsErrPath(err) {
		ret["msg"] = err.Error()
		return
	}
	lastDistance, err := models.GetLastDistance()
	if err == nil {
		distance.Mileage = int(distance.Number) - int(lastDistance.Number)
		if distance.Mileage <= 0 {
			ret["msg"] = "error input"
			return
		}
	}
	err = distance.Insert()
	if err != nil {
		ret["msg"] = err.Error()
		return
	}
	ret["expect"] = (distance.Number) + 70
	ret["succeed"] = true
}

func DistanceDelete(ctx iris.Context) {

}

func DistanceUpdate(ctx iris.Context) {
	var (
		ret = iris.Map{}
		err error
	)
	defer writeResponse(ctx, ret)
	err = models.UpdateMileage()
	if err != nil {
		ret["msg"] = err.Error()
		return
	}
	ret["succeed"] = true
}

func DistanceList(ctx iris.Context) {
	var (
		distances []models.Distance
		err       error
		ret       = iris.Map{}
	)
	defer writeResponse(ctx, ret)
	distances, err = models.DistanceList()
	if err != nil {
		ret["msg"] = err.Error()
		return
	}
	ret["data"] = distances
	ret["succeed"] = true
}

func writeResponse(ctx iris.Context, ret iris.Map) {
	if _, ok := ret["succeed"]; !ok {
		ret["succeed"] = false
	}
	ctx.JSON(ret)
}
