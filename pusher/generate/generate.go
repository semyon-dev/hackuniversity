package generate

import (
	"encoding/json"
	"fmt"
	"github.com/semyon-dev/hackuniversity/pusher/model"
	"math/rand"
)

// генерируем рандомные параметры (типа OPC Server)
func NewData() []byte {
	min := 0.0
	max := 100.0
	data := model.Data{
		PRESSURE: rand.Float64() * (max - min),
		HUMIDITY: rand.Float64() * (max - min),
		TEMP_HOM: rand.Float64() * (max - min),
		TEMP_WOR: rand.Float64() * (max - min),
		LEVEL:    rand.Float64() * (max - min),
		MASS:     rand.Float64() * (max - min),
		WATER:    rand.Float64() * (max - min),
		LEVELCO2: rand.Float64() * (max - min),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Print(err)
	}
	return jsonData
}
