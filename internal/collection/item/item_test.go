package item

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

func init() {
	seed := time.Now().UnixNano()
	println(seed)
	rand.Seed(seed)
}

func testRandItemHead(t *testing.T, idx int, item *Item) {
	t.Helper()
	if idx == 0 {
		if item.isPoint() {
			t.Fatalf("expected false")
		}
		if item.fieldsLen() != 0 {
			t.Fatalf("expected '%v', got '%v'", 0, item.fieldsLen())
		}
		if item.idLen() != 0 {
			t.Fatalf("expected '%v', got '%v'", 0, item.idLen())
		}
	}
	isPoint := rand.Int()%2 == 0
	fieldsLen := int(rand.Uint32() << 2 >> 2)
	idLen := int(rand.Uint32())

	item.setIsPoint(isPoint)
	item.setFieldsLen(fieldsLen)
	item.setIDLen(idLen)

	if item.isPoint() != isPoint {
		t.Fatalf("isPoint: expected '%v', got '%v'", isPoint, item.isPoint())
	}
	if item.fieldsLen() != fieldsLen {
		t.Fatalf("fieldsLen: expected '%v', got '%v'", fieldsLen, item.fieldsLen())
	}
	if item.idLen() != idLen {
		t.Fatalf("idLen: expected '%v', got '%v'", idLen, item.idLen())
	}
}

func TestItemHead(t *testing.T) {
	start := time.Now()
	for time.Since(start) < time.Second/2 {
		var item Item
		for i := 0; i < 10; i++ {
			testRandItemHead(t, i, &item)
		}
	}
}

func testRandItem(t *testing.T) {
	keyb := make([]byte, rand.Int()%16)
	rand.Read(keyb)
	key := string(keyb)

	packed := rand.Int()%2 == 0
	values := make([]float64, rand.Int()%64)
	for i := range values {
		if packed {
			switch rand.Int() % 8 {
			case 0:
				values[i] = 0
			case 1:
				values[i] = 1
			case 2:
				values[i] = -1
			case 3:
				values[i] = float64(rand.Int() % 120)
			case 4:
				values[i] = float64(-(rand.Int() % 120))
			case 5:
				values[i] = float64(rand.Float32())
			case 6:
				values[i] = float64(rand.Float64())
			case 7:
				values[i] = float64(rand.Int()%240) / 10
			}
		} else {
			values[i] = rand.Float64()
		}
	}

	var item *Item
	if rand.Int()%2 == 0 {
		item = New(key, geojson.NewSimplePoint(geometry.Point{X: 1, Y: 2}), packed)
	} else {
		item = New(key, geojson.NewPoint(geometry.Point{X: 1, Y: 2}), packed)
	}
	if item.ID() != key {
		t.Fatalf("expected '%v', got '%v'", key, item.ID())
	}

	var setValues []int
	for _, i := range rand.Perm(len(values)) {
		if !item.SetField(i, values[i]) && values[i] != 0 {
			t.Fatal("expected true")
		}
		setValues = append(setValues, i)
		if item.ID() != key {
			t.Fatalf("expected '%v', got '%v'", key, item.ID())
		}
		for _, i := range setValues {
			if item.GetField(i) != values[i] {
				t.Fatalf("expected '%v', got '%v'", values[i], item.GetField(i))
			}
		}
		fields := item.fields()
		for i := 0; i < len(fields); i++ {
			for _, j := range setValues {
				if i == j {
					if fields[i] != values[i] {
						t.Fatalf("expected '%v', got '%v'", values[i], fields[i])
					}
					break
				}
			}
		}
		weight, points := item.WeightAndPoints()
		if weight != len(fields)*8+len(key)+points*16 {
			t.Fatalf("expected '%v', got '%v'", len(fields)*8+len(key)+points*16, weight)
		}
		if points != 1 {
			t.Fatalf("expected '%v', got '%v'", 1, points)
		}
	}
	if item.GetField(len(values)) != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, item.GetField(len(values)))
	}
	for _, i := range rand.Perm(len(values)) {
		if item.SetField(i, values[i]) {
			t.Fatal("expected false")
		}
	}
	var fvalues []float64
	item.ForEachField(len(values), func(value float64) bool {
		fvalues = append(fvalues, value)
		return true
	})
	if !floatsEquals(values, fvalues) {
		t.Fatalf("expected '%v', got  '%v'", values, fvalues)
	}

	fvalues = nil
	item.ForEachField(len(values), func(value float64) bool {
		if len(fvalues) == 1 {
			return false
		}
		fvalues = append(fvalues, value)
		return true
	})
	if len(values) > 0 && len(fvalues) != 1 {
		t.Fatalf("expected '%v', got '%v'", 1, len(fvalues))
	}

	fvalues = nil
	item.ForEachField(-1, func(value float64) bool {
		fvalues = append(fvalues, value)
		return true
	})
	if !floatsEquals(values, fvalues) {
		t.Fatalf("expected '%v', got '%v'", 1, len(fvalues))
	}

	// should not fail, must allow nil receiver
	(*Item)(nil).ForEachField(1, nil)
	if (*Item)(nil).GetField(1) != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, (*Item)(nil).GetField(1))
	}

	if item.ID() != key {
		t.Fatalf("expected '%v', got '%v'", key, item.ID())
	}
	if item.Obj().NumPoints() != 1 {
		t.Fatalf("expected '%v', got '%v'", 1, item.Obj().NumPoints())
	}
	item.CopyOverFields(values)
	weight, points := item.WeightAndPoints()
	if weight != len(values)*8+len(key)+points*16 {
		t.Fatalf("expected '%v', got '%v'", len(values)*8+len(key)+points*16, weight)
	}
	if points != 1 {
		t.Fatalf("expected '%v', got '%v'", 1, points)
	}
	if !floatsEquals(item.fields(), values) {
		t.Fatalf("expected '%v', got '%v'", values, item.fields())
	}
	item.CopyOverFields(item)
	weight, points = item.WeightAndPoints()
	if weight != len(values)*8+len(key)+points*16 {
		t.Fatalf("expected '%v', got '%v'", len(values)*8+len(key)+points*16, weight)
	}
	if points != 1 {
		t.Fatalf("expected '%v', got '%v'", 1, points)
	}
	if !floatsEquals(item.fields(), values) {
		t.Fatalf("expected '%v', got '%v'", values, item.fields())
	}
	if len(values) > 0 && !item.HasFields() {
		t.Fatal("expected true")
	}

	item.CopyOverFields(nil)
	weight, points = item.WeightAndPoints()
	if weight != len(key)+points*16 {
		t.Fatalf("expected '%v', got '%v'", len(key)+points*16, weight)
	}
	if points != 1 {
		t.Fatalf("expected '%v', got '%v'", 1, points)
	}
	if len(item.fields()) != 0 {
		t.Fatalf("expected '%#v', got '%#v'", 0, len(item.fields()))
	}
	if item.ID() != key {
		t.Fatalf("expected '%v', got '%v'", key, item.ID())
	}
	if item.HasFields() {
		t.Fatal("expected false")
	}

}

func TestItem(t *testing.T) {
	start := time.Now()
	for time.Since(start) < time.Second/2 {
		testRandItem(t)
	}
}

func TestItemFieldOutOfRange(t *testing.T) {
	var item Item
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	item.GetField(-1)
}

func TestItemLess(t *testing.T) {
	item0 := New("0", testString("0"), false)
	item1 := New("1", testString("1"), false)
	item2 := New("1", testString("2"), false)
	item3 := New("3", testString("2"), false)
	if !item0.Less(item1, nil) {
		t.Fatal("expected true")
	}
	if item1.Less(item0, nil) {
		t.Fatal("expected false")
	}
	if !item1.Less(item2, nil) {
		t.Fatal("expected true")
	}
	if item2.Less(item1, nil) {
		t.Fatal("expected false")
	}
	if !item2.Less(item3, nil) {
		t.Fatal("expected true")
	}
	if item3.Less(item2, nil) {
		t.Fatal("expected false")
	}
	weight, points := item0.WeightAndPoints()
	if weight != 2 {
		t.Fatalf("expected '%v', got '%v'", 2, weight)
	}
	if points != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, points)
	}
}

type testString string

func (s testString) Spatial() geojson.Spatial {
	return geojson.EmptySpatial{}
}
func (s testString) ForEach(iter func(geom geojson.Object) bool) bool {
	return iter(s)
}
func (s testString) Empty() bool {
	return true
}
func (s testString) Valid() bool {
	return false
}
func (s testString) Rect() geometry.Rect {
	return geometry.Rect{}
}
func (s testString) Center() geometry.Point {
	return geometry.Point{}
}
func (s testString) AppendJSON(dst []byte) []byte {
	data, _ := json.Marshal(string(s))
	return append(dst, data...)
}
func (s testString) String() string {
	return string(s)
}
func (s testString) JSON() string {
	return string(s.AppendJSON(nil))
}
func (s testString) MarshalJSON() ([]byte, error) {
	return s.AppendJSON(nil), nil
}
func (s testString) Within(obj geojson.Object) bool {
	return false
}
func (s testString) Contains(obj geojson.Object) bool {
	return false
}
func (s testString) Intersects(obj geojson.Object) bool {
	return false
}
func (s testString) NumPoints() int {
	return 0
}
func (s testString) Distance(obj geojson.Object) float64 {
	return 0
}

func floatsEquals(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
