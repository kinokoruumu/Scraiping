package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sclevine/agouti"
)

type Hotel struct {
	Name  string
	Price string
	Image string
}

type Hotels []Hotel

var DB *sql.DB

var NowTime = time.Now()

var baseUrl = "https://www.booking.com/searchresults.ja.html?" +
	"aid=337486&" +
	"sid=25536d504568a4962430499695e3977a&" +
	"class_interval=1&" +
	"dest_id=224&" +
	"dest_type=country&" +
	"dtdisc=0&" +
	"from_sf=1&" +
	"group_adults=2&" +
	"group_children=0&" +
	"inac=0&" +
	"index_postcard=0&" +
	"label_click=undef&" +
	"no_rooms=1&" +
	"postcard=0&" +
	"raw_dest_type=country&" +
	"room1=A%2CA&" +
	"sb_price_type=total&" +
	"src_elem=sb&" +
	"ss_all=0&" +
	"ssb=empty&" +
	"sshis=0&" +
	"ssne=%E3%82%A2%E3%83%A1%E3%83%AA%E3%82%AB&" +
	"ssne_untouched=%E3%82%A2%E3%83%A1%E3%83%AA%E3%82%AB&" +
	"order=review_score_and_price&" +
	"rows=10&"

type Country struct{
	ID int
	Currency string
	KeyWords []string
}

type Keyword struct {
	Keyword string
}

var Countries = []Country{
	{
		2,
		"TWD",
		[]string{
			"台北",
			"高雄",
			"台南",
			"九份",
		},
	},
	{
		3,
		"USD",
		[]string{
			"ニューヨーク",
			"マイアミ",
			"ロスエンジェルス",
			"サンフランシスコ",
			"ラスベガス",
			"ハワイ",
		},
	},
	{
		4,
		"KRW",
		[]string{ // 韓国
			"ソウル",
			"釜山",
			"仁川",
		},
	},
	{
		5,
		"CNY",
		[]string{ // 中国
			"上海",
			"広州",
			"北京",
			"成都",
			"南京",
		},
	},
	{
		6,
		"SGD",
		[]string{ // シンガポール
			"シンガポール",
			"セントーサ",
		},
	},
	{
		7,
		"MXN",
		[]string{ // シンガポール
			"カンクン",
			"グアナファト",
			"グアダラハラ",
			"ケレタロ",
		},
	},
	{
		8,
		"THB",
		[]string{ // シンガポール
			"チェンマイ",
			"バンコク",
			"アユタヤ",
		},
	},
	{
		9,
		"HKD",
		[]string{ // シンガポール
			"ホンコン",
		},
	},
	{
		10,
		"EUR",
		[]string{ // シンガポール
			"ベルリン",
			"ミュンヘン",
			"フランクフルト",
			"パリ",
			"バルセロナ",
			"マドリッド",
			"ローマ",
			"フィレンツェ",
			"ミラノ",
			"ベネチア",
		},
	},
}

func main() {
	db, err := sql.Open("mysql", "root:rootpass@tcp(valuable-production.cci2nw7ztpag.ap-northeast-1.rds.amazonaws.com:3306)/valuable_trip")
	if err != nil {
		fmt.Println("db error")
		panic(err)
	}
	DB = db
	defer DB.Close()

	//driver初期化
	var driver *agouti.WebDriver
	driver = agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	for _, country := range Countries {
		fmt.Println("currency: " + country.Currency)
		fmt.Println("status: scraping...")
		var hotels Hotels
		var recommendHotels Hotels
		for _, keyword := range country.KeyWords {
			hotels = append(hotels, Scraping(keyword, country.Currency, driver)...)
			recommendHotels = GetRecommend(hotels)
		}
		fmt.Println(hotels)
		fmt.Println(recommendHotels)
		Insert(country.ID, hotels)
		InsertByRecommendHotel(country.ID, recommendHotels)
		fmt.Println("-----------------------------------------------------------")
	}
}

func Scraping(keyword string, currency string, driver *agouti.WebDriver) Hotels {
	var hotels Hotels
	fmt.Println(keyword)

	// ページ作成
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	// 総表示ページ数
	pages := 0

	//現在の時刻取得
	nowYear, nowMonth, nowDay := TimeToString(NowTime)
	//翌日の日付を取得
	nextTime := NowTime.AddDate(0, 0, 1)
	nextYear, nextMonth, nextDay := TimeToString(nextTime)


	for {
		fmt.Println(len(hotels))
		// リクエストURL指定
		reqUrl := baseUrl + "&offset=" + strconv.Itoa(pages*10) + "&selected_currency=" + currency + "&ss=" + url.QueryEscape(keyword) +
			"&checkin_month=" + nowMonth +
			"&checkin_monthday=" + nowDay +
			"&checkin_year=" + nowYear +
			"&checkout_month=" + nextMonth +
			"&checkout_monthday=" + nextDay +
			"&checkout_year=" + nextYear

		if err := page.Navigate(reqUrl); err != nil {
			log.Fatalf("%v", err)
		}
		h1, _ := page.FirstByClass("sr_header").Text()
		s1 := strings.Split(h1, "：")
		s2 := strings.Split(s1[1], "おすすめ")
		s2 = strings.Split(s2[0], "（")


		re := regexp.MustCompile("[^0-9]")
		num, _ := strconv.Atoi(re.ReplaceAllString(s2[0], ""))

		// ホテル情報取得
		items := page.AllByClass("sr_item")

		// 取得したホテル数取得
		itemsCount, _ := items.Count()
		for i := 0; i < itemsCount; i++ {
			// hotel作成
			hotel := Hotel{}

			// name取得
			name, nameErr := items.At(i).FirstByClass("sr-hotel__name").Text()
			if nameErr != nil {
				fmt.Println("not found name: continue...")
				continue
			}
			hotel.Name = name

			// price取得
			price, priceErr := items.At(i).FirstByClass("price").Text()
			if priceErr != nil {
				fmt.Println("not found price: continue...")
				continue
			}
			hotel.Price = price

			//画像url取得
			image, imageErr := items.At(i).FirstByClass("hotel_image").Attribute("src")
			if imageErr != nil {
				fmt.Println("not found image: continue...")
				continue
			}
			hotel.Image = image

			// name, price両方取れたらslice格納
			hotels = append(hotels, hotel)

		}
		if len(hotels) < 100 && len(hotels) < num {
			pages++
			continue
		}
		break
	}
	page.CloseWindow()
	return hotels
}

func GetRecommend(hotels Hotels) Hotels {
	var res Hotels
	for i, hotel := range hotels {
		if i < 3 {
			res = append(res, hotel)
		}
	}
	return res
}

func FindValue(hotelData Hotels) (int, int) {
	average := 0            //平均値
	median := 0             //中央値
	var hotelPrice []string //ホテルデータの金額だけを格納するスライス
	var prices []int        //料金計算をするためにINT型へキャストするためスライス
	var parse int

	//金額の部分だけを取り出し、格納する
	for i := 0; i < len(hotelData); i++ {
		hotelPrice = append(hotelPrice, hotelData[i].Price)
	}

	//正規表現で数値以外を除去する
	re := regexp.MustCompile("[^0-9]")

	for i := 0; i < len(hotelPrice); i++ {
		hotelPrice[i] = re.ReplaceAllString(hotelPrice[i], "")
		parse, _ = strconv.Atoi(hotelPrice[i])
		prices = append(prices, parse)

		average += prices[i]
	}

	//昇順にソート
	sort.Sort(sort.IntSlice(prices))
	fmt.Println(prices)

	//平均値の算出
	average = average / len(hotelData)

	//中央値の算出
	if len(hotelData)%2 == 0 {
		//偶数の場合
		median = (prices[len(hotelData)/2-1] + prices[len(hotelData)/2]) / 2
	} else {
		//奇数の場合
		median = prices[len(hotelData)/2]
	}

	return average, median
}

func TimeToString(t time.Time) (string, string, string) {
	timeString := t.Format("2006-1-2")
	timeSplit := strings.Split(timeString, "-")
	return timeSplit[0], timeSplit[1], timeSplit[2]
}

func Insert(countryID int, hotelData Hotels) {
	fmt.Println("insert")
	average, median := FindValue(hotelData)
	fmt.Println("中央値", median)
	fmt.Println("平均値", average)

	//インサートの実行　国ＩＤ，中央値、平均値、取得日付を格納する
	_, err := DB.Exec("INSERT INTO hotels (country_id ,median_price ,average_price ,acquisition_date) VALUES (?,?,?,?)", countryID, median, average, NowTime.Format("2006-01-02"))
	if err != nil {
		fmt.Println("insert error")
		panic(err)
	}
}

//国別に3件ずつホテルの名前,料金,画像のURL,取得日　をDBに格納する
func InsertByRecommendHotel(countryID int, hotelDatas Hotels) {

	//正規表現で料金の数値以外を除去する
	re := regexp.MustCompile("[^0-9]")

	//インサートの実行　国ＩＤ，中央値、平均値、取得日付を格納する。
	for i, hotel := range hotelDatas {
		if _, err := DB.Exec("INSERT INTO recommend_hotels (country_id, num, name, price, image, hotel_acquisition_date) VALUES (?,?,?,?,?,?)", countryID, i + 1, hotel.Name, re.ReplaceAllString(hotel.Price, ""), hotel.Image, NowTime.Format("2006-01-02")); err != nil {
			panic(err)
		}
	}
}
