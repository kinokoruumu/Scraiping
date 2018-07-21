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

var hotels Hotels

var getTime string
var baseUrl string = "https://www.booking.com/searchresults.ja.html?" +
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

func main() {

	var countryId int        //国別のＩＤ
	var countryHotels Hotels //国別に3件のホテルデータ格納用の構造体

	// 台湾
	countryId = 2
	var twdHotels Hotels
	twdHotels = append(twdHotels, Scraping("台北", 100, "TWD")...)
	twdHotels = append(twdHotels, Scraping("高雄", 100, "TWD")...)
	twdHotels = append(twdHotels, Scraping("台南", 100, "TWD")...)
	twdHotels = append(twdHotels, Scraping("九份", 20, "TWD")...)
	fmt.Println(twdHotels)
	fmt.Println(len(twdHotels))
	Insert(twdHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("台湾", 3, "TWD")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	// アメリカ
	countryId = 3
	var usdHotels Hotels
	usdHotels = append(usdHotels, Scraping("ニューヨーク", 100, "USD")...)
	usdHotels = append(usdHotels, Scraping("マイアミ", 100, "USD")...)
	usdHotels = append(usdHotels, Scraping("ロスエンジェルス", 100, "USD")...)
	usdHotels = append(usdHotels, Scraping("サンフランシスコ", 100, "USD")...)
	usdHotels = append(usdHotels, Scraping("ラスベガス", 100, "USD")...)
	usdHotels = append(usdHotels, Scraping("ハワイ", 100, "USD")...)
	fmt.Println(usdHotels)
	fmt.Println(len(usdHotels))
	Insert(usdHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("アメリカ", 3, "USD")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	// 韓国
	countryId = 4
	var krwHotels Hotels
	krwHotels = append(krwHotels, Scraping("ソウル", 100, "KRW")...)
	krwHotels = append(krwHotels, Scraping("釜山", 100, "KRW")...)
	krwHotels = append(krwHotels, Scraping("仁川", 100, "KRW")...)
	fmt.Println(krwHotels)
	fmt.Println(len(krwHotels))
	Insert(krwHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("韓国", 3, "KRW")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	//中国
	countryId = 5
	var cnyHotels Hotels
	cnyHotels = append(cnyHotels, Scraping("上海", 100, "CNY")...)
	cnyHotels = append(cnyHotels, Scraping("広州	", 100, "CNY")...)
	cnyHotels = append(cnyHotels, Scraping("北京", 100, "CNY")...)
	cnyHotels = append(cnyHotels, Scraping("成都", 100, "CNY")...)
	cnyHotels = append(cnyHotels, Scraping("南京", 100, "CNY")...)
	fmt.Println(countryId)
	fmt.Println(cnyHotels)
	fmt.Println(len(cnyHotels))
	Insert(cnyHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("中国", 3, "CNY")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	//　シンガポール
	countryId = 6
	var sgdHotels Hotels
	sgdHotels = append(sgdHotels, Scraping("シンガポール", 100, "SGD")...)
	sgdHotels = append(sgdHotels, Scraping("セントーサ", 100, "SGD")...)
	fmt.Println(sgdHotels)
	fmt.Println(len(sgdHotels))
	Insert(sgdHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("シンガポール", 3, "SGD")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	//　メキシコ
	countryId = 7
	var mxnHotels Hotels
	mxnHotels = append(mxnHotels, Scraping("カンクン", 100, "MXN")...)
	mxnHotels = append(mxnHotels, Scraping("グアナファト", 30, "MXN")...)
	mxnHotels = append(mxnHotels, Scraping("グアダラハラ", 100, "MXN")...)
	mxnHotels = append(mxnHotels, Scraping("ケレタロ", 100, "MXN")...)
	fmt.Println(mxnHotels)
	fmt.Println(len(mxnHotels))
	Insert(mxnHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("メキシコ", 3, "MXN")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	//	タイ
	countryId = 8
	var thbHotels Hotels
	thbHotels = append(thbHotels, Scraping("チェンマイ", 100, "THB")...)
	thbHotels = append(thbHotels, Scraping("バンコク", 100, "THB")...)
	thbHotels = append(thbHotels, Scraping("アユタヤ", 50, "THB")...)
	fmt.Println(thbHotels)
	fmt.Println(len(thbHotels))
	Insert(thbHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("タイ", 3, "THB")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	//	香港
	countryId = 9
	var hkdHotels Hotels
	hkdHotels = append(hkdHotels, Scraping("ホンコン", 100, "HKD")...)
	fmt.Println(hkdHotels)
	fmt.Println(len(hkdHotels))
	Insert(hkdHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("香港", 3, "HKD")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	//	//	ユーロ圏
	countryId = 10
	var eurHotels Hotels
	eurHotels = append(eurHotels, Scraping("ベルリン", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("ミュンヘン", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("フランクフルト", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("パリ", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("バルセロナ", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("マドリッド", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("ローマ", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("フィレンツェ", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("ミラノ", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("ベネチア", 100, "EUR")...)
	eurHotels = append(eurHotels, Scraping("ウィーン", 100, "EUR")...)
	fmt.Println(eurHotels)
	fmt.Println(len(eurHotels))
	Insert(eurHotels, countryId)
	countryHotels = append(countryHotels, CountryScraping("ドイツ", 3, "EUR")...)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	fmt.Println(countryHotels)
	fmt.Println(len(countryHotels))

	InsertByRecommendHotel(countryHotels)

}

func Scraping(keyword string, number int, currency string) Hotels {
	var hotels Hotels

	//driver初期化
	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}

	defer driver.Stop()

	// ページ作成
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	// 総表示ページ数
	pages := 0

	//現在の時刻取得
	nowTime := time.Now()
	getTime = nowTime.Format("2006-01-02")
	nowYear, nowMonth, nowDay := TimetoString(nowTime)
	//翌日の日付を取得
	nextTime := nowTime.AddDate(0, 0, 1)
	nextYear, nextMonth, nextDay := TimetoString(nextTime)

	for len(hotels) < number {
		//TODO 日付動的化
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

		// ホテル情報取得
		items := page.AllByClass("sr_item")

		// 取得したホテル数取得
		itemsCount, _ := items.Count()
		for i := 0; i < itemsCount; i++ {
			// hotel作成
			hotel := Hotel{}

			// name取得
			if name, nameErr := items.At(i).FirstByClass("sr-hotel__name").Text(); nameErr == nil {
				hotel.Name = name
			} else {
				fmt.Println("not found name: continue...")
				continue
			}

			// price取得
			if price, priceErr := items.At(i).FirstByClass("price").Text(); priceErr == nil {
				hotel.Price = price
			} else {
				fmt.Println("not found price: continue...")
				continue
			}

			// name, price両方取れたらslice格納
			hotels = append(hotels, hotel)
		}
		pages++
	}
	return hotels
}

//ホテルを国別に3件取得する関数（名前、料金、画像のURL）
func CountryScraping(keyword string, number int, currency string) Hotels {
	var hotels Hotels

	//driver初期化
	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}

	defer driver.Stop()

	// ページ作成
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	// 総表示ページ数
	pages := 0

	//現在の時刻取得
	nowTime := time.Now()
	getTime = nowTime.Format("2006-01-02")
	nowYear, nowMonth, nowDay := TimetoString(nowTime)
	//翌日の日付を取得
	nextTime := nowTime.AddDate(0, 0, 1)
	nextYear, nextMonth, nextDay := TimetoString(nextTime)

	//TODO 日付動的化
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

	// ホテル情報取得
	items := page.AllByClass("sr_item")

	for i := 0; i < number+1; i++ {
		// hotel作成
		hotel := Hotel{}

		// name取得
		if name, nameErr := items.At(i).FirstByClass("sr-hotel__name").Text(); nameErr == nil {
			hotel.Name = name
		} else {
			fmt.Println("not found name: continue...")
			continue
		}

		// price取得
		if price, priceErr := items.At(i).FirstByClass("price").Text(); priceErr == nil {
			hotel.Price = price
		} else {
			fmt.Println("not found price: continue...")
			continue
		}

		//画像url取得
		if image, imageErr := items.At(i).FirstByClass("hotel_image").Attribute("src"); imageErr == nil {
			hotel.Image = image
		} else {
			fmt.Println("not found image: continue...")
			continue
		}

		// name, price, image　が取れたらslice格納
		hotels = append(hotels, hotel)
	}

	return hotels
}

func Findvalue(hoteldatas Hotels) (int, int) {
	average := 0            //平均値
	median := 0             //中央値
	var hotelPrice []string //ホテルデータの金額だけを格納するスライス
	var prices []int        //料金計算をするためにINT型へキャストするためスライス
	var parse int

	//金額の部分だけを取り出し、格納する
	for i := 0; i < len(hoteldatas); i++ {
		hotelPrice = append(hotelPrice, hoteldatas[i].Price)
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
	average = average / len(hoteldatas)

	//中央値の算出
	if len(hoteldatas)%2 == 0 {
		//偶数の場合
		median = (prices[len(hoteldatas)/2-1] + prices[len(hoteldatas)/2]) / 2
	} else {
		//奇数の場合
		median = prices[len(hoteldatas)/2]
	}

	return average, median
}

func TimetoString(t time.Time) (string, string, string) {
	timeString := t.Format("2006-1-2")
	timeSplit := strings.Split(timeString, "-")
	return timeSplit[0], timeSplit[1], timeSplit[2]
}

func Insert(hoteldatas Hotels, countryId int) {

	average, median := Findvalue(hoteldatas)
	fmt.Println("中央値", median)
	fmt.Println("平均値", average)

	//データベース接続
	con, err := sql.Open("mysql", "root:rootpass@tcp(valuable-production.cci2nw7ztpag.ap-northeast-1.rds.amazonaws.com:3306)/valuable_trip")
	if err != nil {
		panic(err)
	}
	defer con.Close()

	//インサートの実行　国ＩＤ，中央値、平均値、取得日付を格納する
	if _, err := con.Exec("INSERT INTO hotels (country_id ,median_price ,average_price ,acquisition_date) VALUES (?,?,?,?)", countryId, median, average, getTime); err != nil {
		panic(err)
	}
}

//国別に3件ずつホテルの名前,料金,画像のURL,取得日　をDBに格納する
func InsertByRecommendHotel(hoteldatas Hotels) {

	var hotelName []string
	var hotelPrice []string
	var hotelImage []string

	//正規表現で料金の数値以外を除去する
	re := regexp.MustCompile("[^0-9]")
	for i := 0; i < len(hoteldatas); i++ {
		hotelName = append(hotelName, hoteldatas[i].Name)
		hotelPrice = append(hotelPrice, hoteldatas[i].Price)
		hotelPrice[i] = re.ReplaceAllString(hotelPrice[i], "")
		hotelImage = append(hotelImage, hoteldatas[i].Image)
	}
	fmt.Println("ホテル名：", hotelName)
	fmt.Println("料金名：", hotelPrice)
	fmt.Println("画像：", hotelImage)

	//データベース接続
	con, err := sql.Open("mysql", "root:rootpass@tcp(valuable-production.cci2nw7ztpag.ap-northeast-1.rds.amazonaws.com:3306)/valuable_trip")
	if err != nil {
		panic(err)
	}
	defer con.Close()

	countryId := 2
	//インサートの実行　国ＩＤ，中央値、平均値、取得日付を格納する。
	for i := 0; i < len(hoteldatas); i++ {
		if i%3 == 0 && i != 0 {
			countryId++
		}
		if _, err := con.Exec("INSERT INTO recommend_hotels (country_id ,name ,price ,image,hotel_acquisition_date) VALUES (?,?,?,?,?)", countryId, hotelName[i], hotelPrice[i], hotelImage[i], getTime); err != nil {
			panic(err)
		}
	}

}
