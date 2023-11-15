// map[string]interface{}型をパースして値を取り出す関数を提供します。
package seeker

import (
	"time"
)

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合やstring型ではない場合はゼロ値を返します。
func Str(d map[string]interface{}, rp []string) string {
	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return ""
	}

	if len(rp) == 0 {
		v, ok := i.(string)
		if !ok {
			return ""
		}

		return v
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return ""
		}

		return Str(m, rp)
	}

	return ""
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合やint型ではない場合はゼロ値を返します。
func Int(d map[string]interface{}, rp []string) int {
	cp := rp[0]
	rp = rp[1:]
	i, ok := d[cp]
	if !ok {
		return 0
	}

	if len(rp) == 0 {
		v, ok := i.(float64)
		if !ok {
			return 0
		}

		return int(v)
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return 0
		}

		return Int(m, rp)
	}

	return 0
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合やuint型ではない場合はゼロ値を返します。
func Uint(d map[string]interface{}, rp []string) uint {
	cp := rp[0]
	rp = rp[1:]
	i, ok := d[cp]
	if !ok {
		return 0
	}

	if len(rp) == 0 {
		v, ok := i.(float64)
		if !ok {
			return 0
		}

		return uint(v)
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return 0
		}

		return Uint(m, rp)
	}
	return 0
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合やfloat64型ではない場合はゼロ値を返します。
func Float64(d map[string]interface{}, rp []string) float64 {
	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return 0
	}

	if len(rp) == 0 {
		v, ok := i.(float64)
		if !ok {
			return 0
		}

		return v
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return 0
		}

		return Float64(m, rp)
	}

	return 0
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合やbool型ではない場合はゼロ値を返します。
func Bool(d map[string]interface{}, rp []string) bool {
	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return false
	}

	if len(rp) == 0 {
		v, ok := i.(bool)
		if !ok {
			return false
		}

		return v
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return false
		}

		return Bool(m, rp)
	}

	return false
}

// 第一引数dのmap[string]interfaceデータ構造から、第二引数rpで指定したキーを再帰的に検索してその値をtime型に変換して取得します
// 値が見つからなかった場合は、ゼロ値を返します
func Time(d map[string]interface{}, rp []string) time.Time {
	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return time.Time{}
	}

	if len(rp) == 0 {
		v, ok := i.(string)
		if !ok {
			return time.Time{}
		}
		vv, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}
		}

		return vv
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return time.Time{}
		}

		return Time(m, rp)
	}

	return time.Time{}
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合や[]interface[]型ではない場合は、空のmap[string]interface{}型のスライスを返します。
func Slice(d map[string]interface{}, rp []string) []map[string]interface{} {
	var empty []map[string]interface{}

	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return empty
	}

	if len(rp) == 0 {
		v, ok := i.([]interface{})
		if !ok {
			return empty
		}

		var mis []map[string]interface{}
		for _, m := range v {
			mis = append(mis, m.(map[string]interface{}))
		}

		return mis
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return empty
		}

		return Slice(m, rp)
	}

	return empty
}

// SliceString は map[string]interface{} 型から string のスライスを取得します。
// 指定されたキーが存在しない場合、または []interface{} 型でない場合は空のスライスを返します。
func SliceString(d map[string]interface{}, rp []string) []string {
	var empty []string

	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return empty
	}

	if len(rp) == 0 {
		v, ok := i.([]interface{})
		if !ok {
			return empty
		}

		var mis []string
		for _, m := range v {
			mis = append(mis, m.(string))
		}

		return mis
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return empty
		}

		return SliceString(m, rp)
	}

	return empty
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合やmap[stringはinterface{}型ではない場合は、空のmap[string]interface{}型を返します。
func Raw(d map[string]interface{}, rp []string) map[string]interface{} {
	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return map[string]interface{}{}
	}

	if len(rp) == 0 {
		v, ok := i.(map[string]interface{})
		if !ok {
			return map[string]interface{}{}
		}

		return v
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return map[string]interface{}{}
		}

		return Raw(m, rp)
	}

	return map[string]interface{}{}
}

// dからrpで指定したキーを再帰的に検索してその値を取得します
//
// 値が見つからなかった場合は、空のinterface{}型を返します。
func Interface(d map[string]interface{}, rp []string) interface{} {
	var empty interface{}

	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return empty
	}

	if len(rp) == 0 {
		return i
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return empty
		}

		return Interface(m, rp)
	}

	return empty
}

// JSONの指定したパスに値があるか確認します
func IsExist(d map[string]interface{}, rp []string) bool {
	cp := rp[0]
	rp = rp[1:]

	i, ok := d[cp]
	if !ok {
		return false
	}

	if len(rp) == 0 {
		_, ok := i.(interface{})
		if !ok {
			return false
		}

		return true
	}

	if len(rp) > 0 {
		m, ok := i.(map[string]interface{})
		if !ok {
			return false
		}

		return IsExist(m, rp)
	}

	return false
}
