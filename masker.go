// Package masker Provide mask format of Taiwan usually used(Name, Address, Email, ID ...etc.),
package masker

import (
	"math"
	"reflect"
	"strings"
)

const tagName = "mask"

type mtype string

// Maske Types of format string
const (
	MPassword   mtype = "password"
	MName             = "name"
	MAddress          = "addr"
	MEmail            = "email"
	MMobile           = "mobile"
	MTelephone        = "tel"
	MId               = "id"
	MCreditCard       = "credit"
	MStruct           = "struct"
)

// Masker is a instance to marshal masked string
type Masker struct{}

// Struct must input a interface{}, add tag mask on struct fields, after Struct(), return a pointer interface{} of input type and it will be masked with the tag format type
//
// Example:
//
//   type Foo struct {
//       Name      string `mask:"name"`
//       Email     string `mask:"email"`
//       Password  string `mask:"password"`
//       ID        string `mask:"id"`
//       Address   string `mask:"addr"`
//       Mobile    string `mask:"mobile"`
//       Telephone string `mask:"tel"`
//       Credit    string `mask:"credit"`
//       Foo       *Foo   `mask:"struct"`
//   }
//
//   func main() {
//       s := &Foo{
//           Name: ...,
//           Email: ...,
//           Password: ...,
//           Foo: &{
//               Name: ...,
//               Email: ...,
//               Password: ...,
//           }
//       }
//
//       m := masker.New()
//
//       t, err := m.Struct(s)
//
//       fmt.Println(t.(*Foo))
//   }
func (m *Masker) Struct(s interface{}) (interface{}, error) {
	var selem, tptr reflect.Value

	st := reflect.TypeOf(s)

	if st.Kind() == reflect.Ptr {
		tptr = reflect.New(st.Elem())
		selem = reflect.ValueOf(s).Elem()
	} else {
		tptr = reflect.New(st)
		selem = reflect.ValueOf(s)
	}

	for i := 0; i < selem.NumField(); i++ {
		if mtag, ok := selem.Type().Field(i).Tag.Lookup(tagName); ok {
			switch mtype(mtag) {
			case MPassword:
				tptr.Elem().Field(i).SetString(m.Password(selem.Field(i).String()))
			case MName:
				tptr.Elem().Field(i).SetString(m.Name(selem.Field(i).String()))
			case MAddress:
				tptr.Elem().Field(i).SetString(m.Address(selem.Field(i).String()))
			case MEmail:
				tptr.Elem().Field(i).SetString(m.Email(selem.Field(i).String()))
			case MMobile:
				tptr.Elem().Field(i).SetString(m.Mobile(selem.Field(i).String()))
			case MId:
				tptr.Elem().Field(i).SetString(m.ID(selem.Field(i).String()))
			case MTelephone:
				tptr.Elem().Field(i).SetString(m.Telephone(selem.Field(i).String()))
			case MCreditCard:
				tptr.Elem().Field(i).SetString(m.CreditCard(selem.Field(i).String()))
			case MStruct:
				if !selem.Field(i).IsNil() {
					_t, err := m.Struct(selem.Field(i).Interface())
					if err != nil {
						return nil, err
					}
					tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
				}
			default:
				tptr.Elem().Field(i).Set(selem.Field(i))
			}
		} else {
			tptr.Elem().Field(i).Set(selem.Field(i))
		}
	}

	return tptr.Interface(), nil
}

// Name mask the second world and the third world
//
// Example:
//   input: ABCD
//   output: A**D
func (*Masker) Name(i string) string {
	l := len(i)

	if l == 2 || l == 3 {
		return overlay(i, "**", 1, 2)
	}

	if l > 3 {
		return overlay(i, "**", 1, 3)
	}

	return "**"
}

// ID mask last 4 worlds of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func (*Masker) ID(i string) string {
	return overlay(i, "****", 6, 10)
}

// Address keep first 6 worlds, mask the overs
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func (*Masker) Address(i string) string {
	l := len(i)
	if l <= 6 {
		return "******"
	}
	return overlay(i, "******", 6, math.MaxInt64)
}

// CreditCard mask middle 6 worlds from 7'th world
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345` (American Express)(len = 15)
//   output2: 123456******345`
func (*Masker) CreditCard(i string) string {
	return overlay(i, "******", 6, 12)
}

// Email keep domain and first 3 worlds
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func (*Masker) Email(i string) string {
	tmp := strings.Split(i, "@")
	addr := tmp[0]
	domain := tmp[1]

	addr = overlay(addr, "****", 3, 7)

	return addr + "@" + domain
}

// Mobile mask mobile 3 worlds from 4'th world
//
// Example:
//   input: 0987654321
//   output: 0987***321
func (*Masker) Mobile(i string) string {
	return overlay(i, "***", 4, 7)
}

// Telephone remove `(`, `)`, ` `, `-` chart, and mask last 4 worlds of telephone number, format to `(??)????-????`
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func (*Masker) Telephone(i string) string {
	i = strings.Replace(i, " ", "", -1)
	i = strings.Replace(i, "(", "", -1)
	i = strings.Replace(i, ")", "", -1)
	i = strings.Replace(i, "-", "", -1)

	l := len(i)

	if l != 10 && l != 8 {
		return i
	}

	ans := ""

	if l == 10 {
		ans += "("
		ans += i[:2]
		ans += ")"
		i = i[2:]
	}

	ans += i[:4]
	ans += "-"
	ans += "****"

	return ans
}

// Password always return `************`
func (*Masker) Password(i string) string {
	return "************"
}

// New create Masker
func New() *Masker {
	return &Masker{}
}

var instance *Masker

func init() {
	instance = New()
}

// Struct must input a interface{}, add tag mask on struct fields, after Struct(), return a pointer interface{} of input type and it will be masked with the tag format type
//
// Example:
//
//   type Foo struct {
//       Name      string `mask:"name"`
//       Email     string `mask:"email"`
//       Password  string `mask:"password"`
//       ID        string `mask:"id"`
//       Address   string `mask:"addr"`
//       Mobile    string `mask:"mobile"`
//       Telephone string `mask:"tel"`
//       Credit    string `mask:"credit"`
//       Foo       *Foo   `mask:"struct"`
//   }
//
//   func main() {
//       s := &Foo{
//           Name: ...,
//           Email: ...,
//           Password: ...,
//           Foo: &{
//               Name: ...,
//               Email: ...,
//               Password: ...,
//           }
//       }
//
//       t, err := masker.Struct(s)
//
//       fmt.Println(t.(*Foo))
//   }
func Struct(s interface{}) (interface{}, error) {
	return instance.Struct(s)
}

// Name mask the second world and the third world
//
// Example:
//   input: ABCD
//   output: A**D
func Name(i string) string {
	return instance.Name(i)
}

// ID mask last 4 worlds of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func ID(i string) string {
	return instance.ID(i)
}

// Address keep first 6 worlds, mask the overs
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func Address(i string) string {
	return instance.Address(i)
}

// CreditCard mask middle 6 worlds from 7'th world
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345 (American Express)(len = 15)
//   output2: 123456******345
func CreditCard(i string) string {
	return instance.CreditCard(i)
}

// Email keep domain and first 3 worlds
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func Email(i string) string {
	return instance.Email(i)
}

// Mobile mask mobile 3 worlds from 4'th world
//
// Example:
//   input: 0987654321
//   output: 0987***321
func Mobile(i string) string {
	return instance.Mobile(i)
}

// Telephone remove `(`, `)`, ` `, `-` chart, and mask last 4 worlds of telephone number, format to `(??)????-????`
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func Telephone(i string) string {
	return instance.Telephone(i)
}

// Password always return `************`
func Password(i string) string {
	return instance.Password(i)
}
