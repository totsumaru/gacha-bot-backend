package debug

import (
	"encoding/json"
	"log"
)

// 検証時にしか使用しません
func PrintJSON(v interface{}) {
	// JSONに一度変換
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("初期JSON変換に失敗しました: %v", err)
	}

	// map[string]interface{}に変換
	var dat map[string]interface{}
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatalf("mapへの変換に失敗しました: %v", err)
	}

	// 整形されたJSONに再変換
	b, err = json.MarshalIndent(dat, "", "  ")
	if err != nil {
		log.Fatalf("整形されたJSONへの変換に失敗しました: %v", err)
	}

	// コンソールに出力
	log.Println(string(b))
}
