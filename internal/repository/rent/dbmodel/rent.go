package dbmodel

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
)

type DBGeoPoint struct {
	Lat  float64
	Long float64
}

func (p *DBGeoPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Long, p.Lat), nil
}

func (p *DBGeoPoint) Scan(val any) error {
	var source []byte
	switch v := val.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return fmt.Errorf("conflict type for DBGeoPoint: %T", val)
	}

	// 1. Декодируем HEX-строку от PostGIS в байты
	decoded := make([]byte, hex.DecodedLen(len(source)))
	_, err := hex.Decode(decoded, source)
	if err != nil {
		decoded = source
	}

	// Минимальный размер EWKB для Point (4326) — 25 байт
	if len(decoded) < 25 {
		return fmt.Errorf("invalid EWKB length: %d", len(decoded))
	}

	// 2. Определяем порядок байт (Byte Order)
	// 1 = LittleEndian (обычно в Postgres), 0 = BigEndian
	var order binary.ByteOrder = binary.BigEndian
	if decoded[0] == 1 {
		order = binary.LittleEndian
	}

	// 3. Считываем тип геометрии (байты 1-5)
	geomType := order.Uint32(decoded[1:5])

	var offset int
	// Если в типе геометрии есть флаг SRID (0x20000000), пропускаем 4 байта самого SRID
	if (geomType & 0x20000000) != 0 {
		offset = 9 // 1 (byte order) + 4 (geom type) + 4 (SRID)
	} else {
		offset = 5 // 1 (byte order) + 4 (geom type)
	}

	// Проверяем, что в массиве осталось достаточно байт для двух float64 (8 + 8 = 16 байт)
	if len(decoded) < offset+16 {
		return fmt.Errorf("insufficient bytes for coordinates")
	}

	// 4. Читаем координаты напрямую из байт-кода
	longBits := order.Uint64(decoded[offset : offset+8])
	latBits := order.Uint64(decoded[offset+8 : offset+16])

	// Превращаем биты в числа с плавающей точкой float64
	p.Long = math.Float64frombits(longBits)
	p.Lat = math.Float64frombits(latBits)

	return nil // Успех!
}

type Rent struct {
	database.Base
	Location DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
	Address  string     `gorm:"column:address;type:varchar;size:255"`
	Info     *string     `gorm:"column:info;type:varchar;size:255;default:null"`
}
