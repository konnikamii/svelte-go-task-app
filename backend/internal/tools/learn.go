package tools

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var m = sync.RWMutex{}
var wg = sync.WaitGroup{}
var dbData = []string{"id1", "id2", "id3", "id4"}
var results = []string{}

func main() {

	println("Started main func")

	errorHandling()
	dataTypes()
	goRoutines()
	channels()
	atomicValues()
}

func channels() {
	var c = make(chan int, 8)
	go process(c)
	for i := range c {
		println(i)
		time.Sleep(time.Second * 1)
	}
}

func process(c chan int) {
	defer close(c)
	for i := range 5 {
		c <- i
	}
	println("exiting")
}

func errorHandling() {
	// handle errors
	println("\n\n-------------------- || --------------------")
	println("Handle Errors")
	a, b, c, d := 1, 4, 8, 0
	var res, rem, err = division(a, b)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Res: %v Rem: %v \n", res, rem)
	}
	res2, rem2, err2 := division(c, d)
	switch {
	case err2 != nil:
		fmt.Println(err2.Error())
	case rem2 == 0:
		fmt.Printf("Res: %v \n", res2)
	default:
		fmt.Printf("Res: %v Rem: %v \n", res2, rem2)
	}
}

func dataTypes() {
	// assigning types
	println("\n\n-------------------- || --------------------")
	println("Assigning Types")
	const typedWithValue int = 1
	const untypedWithValue = 1
	var typedNoValue int // NOTE must be var
	untypedVarWithValue := 1.0
	println(typedNoValue, untypedVarWithValue)

	// types
	println("\n--------------------")
	println("Numeric Types")
	var integer int     // like int32 or int64
	var integer8 int8   // from -128 to 127
	var integer16 int16 // from -32,768 to 32,767
	var integer32 int32 // from -2,147,483,648 to 2,147,483,647
	var integer64 int64 // from -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807

	var uinteger uint     // like uint32 or uint64
	var uinteger8 uint8   // from 0 to 255
	var uinteger16 uint16 // from 0 to 65,535
	var uinteger32 uint32 // from 0 to 4,294,967,295
	var uinteger64 uint64 // from 0 to 18,446,744,073,709,551,615

	fmt.Printf(`
int %v --- like int32 or int64
int8 %v --- Range -128 to 127
int16 %v --- Range -32,768 to 32,767
int32 %v --- Range -2,147,483,648 to 2,147,483,647
int64 %v --- Range -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807

uint %v --- like uint32 or uint64
uint8 %v --- Range 0 to 255
uint16 %v --- Range 0 to 65,535
uint32 %v --- Range 0 to 4,294,967,295
uint64 %v --- Range 0 to 18,446,744,073,709,551,615
	`,
		integer, integer8, integer16, integer32, integer64,
		uinteger, uinteger8, uinteger16, uinteger32, uinteger64)

	var runeChar rune // Synonym for int32 (Unicode code points)
	var byteChar byte // Synonym for uint8

	fmt.Printf(`
rune %v --- Synonym for int32 (Unicode code points)
byte %v --- Synonym for uint8
	`, runeChar, byteChar)

	var decimal32 float32 = 123456789.9876 // from x to x
	var decimal64 float64                  // from x to x

	fmt.Printf(`
float32 %v --- Range -x to x
float64 %v --- Range -x to x
	`, decimal32, decimal64)

	var complexNum64 complex64 = complex(9, 2)
	var complexNum128 complex128 = complex(6, 2)

	fmt.Printf(`
complex64 %v --- Type %T
complex128 %v --- Type %T
	`, complexNum64, complexNum64, complexNum128, complexNum128)

	println("\n--------------------")
	println("String Types (immutable)")
	var str string = "'I am a sequence of runes'"

	fmt.Printf(`
string %v --- Type %T
	`, str, str)

	println("\n--------------------")
	println("Arrays/Slices Types (mutable)")
	var array1 [3]int
	array2 := [3]string{"abc", "def", "ghi"}
	array3 := [3][2]string{{"abc", "def"}, {"abc2", "def2"}}
	array4 := [...]int{1, 2, 3, 4, 5, 6, 7}
	slice1 := []int8{3, 6, 9}
	slice2 := append(slice1, 9)
	slice3 := append(slice1, slice2...)
	slice4 := make([]int32, 3, 10)

	fmt.Printf(`
array1 %v --- Type %T
array2 %v --- Type %T
array3 %v --- Type %T
array4 %v --- Type %T
slice1 %v --- Type %T --- Length %v Capacity %v
slice2 %v --- Type %T --- Length %v Capacity %v
slice3 %v --- Type %T --- Length %v Capacity %v
slice4 %v --- Type %T --- Length %v Capacity %v
	`, array1, array1, array2, array2, array3, array3, array4, array4,
		slice1, slice1, len(slice1), cap(slice1), slice2, slice2, len(slice2), cap(slice2),
		slice3, slice3, len(slice3), cap(slice3), slice4, slice4, len(slice4), cap(slice4))

	println("\n--------------------")
	println("Maps Types (mutable)")
	map1 := make(map[string]uint8)
	map2 := map[string]uint8{"Adam": 2, "Thomash": 4}

	fmt.Printf(`
array1 %v --- Type %T
array2 %v --- Type %T
	`, map1, map1, map2, map2)

	println("\n--------------------")
	println("Structs Types")
	type structType1 struct {
		param1 uint
		param2 string
	}
	type structType3 struct {
		p1 uint
	}
	type structType2 struct {
		p1 uint
		p2 structType1
		p3 float32
		structType3
	}
	struct1 := structType1{param1: 5, param2: "abc"}
	var struct2 structType1
	struct3 := structType2{1, structType1{1, "test"}, 7, structType3{9}}

	fmt.Printf(`
struct1 %v --- Type %T
struct2 %v --- Type %T
struct3 %v --- Type %T
	`, struct1, struct1, struct2, struct2, struct3, struct3)

	println("\n--------------------")
	println("Pointers Types")
	someInt32 := int32(5)
	var pointer1 *int32 = new(int32(4))
	var pointer2 *int32 = new(int32)
	pointer2 = &someInt32
	*pointer2 = 53

	fmt.Printf(`
someInt32 %v
pointer1 %v --- Type %T --- Value %v
pointer2 %v --- Type %T --- Value %v
	`, someInt32, pointer1, pointer1, *pointer1, pointer2, pointer2, *pointer2)

}

func goRoutines() {
	t0 := time.Now()
	for i := 0; i < len(dbData); i++ {
		wg.Add(1)
		go callDb(i)
	}
	wg.Wait()
	fmt.Printf("\nTotal execution time: %v", time.Since(t0))
	fmt.Printf("\n Result: %v", results)
}

// #region atomicValues
func atomicValues() {
	player := NewPlayer()
	go startUILoop(player)
	startGameLoop(player)
}

type Player struct {
	health int32
}

func (p *Player) getHealthAtomic() int {
	return int(atomic.LoadInt32(&p.health))
}

func (p *Player) takeDmgAtomic(value int) {
	health := p.getHealthAtomic()
	atomic.StoreInt32(&p.health, int32(health-value))
}

func NewPlayer() *Player {
	return &Player{health: 100}
}

func startUILoop(p *Player) {
	ticker := time.NewTicker(time.Millisecond * 150)
	for {
		fmt.Printf("Player health: %d\n", p.getHealthAtomic())
		<-ticker.C
	}
}

func startGameLoop(p *Player) {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		p.takeDmgAtomic(rand.Intn(12))
		if p.getHealthAtomic() <= 0 {
			fmt.Println("\nGame Over")
			break
		}
		<-ticker.C
	}
}

// #endregion atomicValues

// ------------ Helpers ------------
func division(a int, b int) (int, int, error) {
	var err error

	if b == 0 {
		err = errors.New("B cannot be 0")
		return 0, 0, err
	}

	var result = a / b
	var remainder = a % b
	return result, remainder, err
}

func callDb(i int) {
	println("Fetching from db for", i)
	delay := rand.Float32() * 2000
	sleepFor := time.Duration(delay) * time.Millisecond
	time.Sleep(sleepFor)
	fmt.Printf("\n Finished: %v s", sleepFor.Seconds())
	m.Lock()
	results = append(results, dbData[i])
	m.Unlock()
	m.RLock()
	fmt.Printf("\n Current results are: %v", results)
	m.RUnlock()
	wg.Done()
}
