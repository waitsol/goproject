package main

type Data_json struct {
	SubType     string `json:"SubType"`
	ReqID       int    `json:"ReqID"`
	Inst        string `json:"Inst"`
	Market      string `json:"Market"`
	ServiceType string `json:"ServiceType"`
}

type Daya_json struct {
	OrgCode string `json:"OrgCode"`
	Token   string `json:"Token"`
	AppName string `json:"AppName"`
	AppVer  string `json:"AppVer"`
	AppType string `json:"AppType"`
	Tag     string `json:"Tag"`
}
type DynaType struct {
	TradingDay     int       `json:"TradingDay"`
	Time           int       `json:"Time"`
	HighestPrice   float64   `json:"HighestPrice"` //最高价格
	LowestPrice    float64   `json:"LowestPrice"`  //最低价格
	LastPrice      float64   `json:"LastPrice"`
	Volume         int       `json:"Volume"`
	Amount         float64   `json:"Amount"`
	TickCount      int       `json:"TickCount"`
	BuyPrice       []float64 `json:"BuyPrice"`
	BuyVolume      []int     `json:"BuyVolume"`
	SellPrice      []float64 `json:"SellPrice"`
	SellVolume     []int     `json:"SellVolume"`
	AveragePrice   float64   `json:"AveragePrice"`
	Wk52High       float64   `json:"Wk52High"`
	Wk52Low        float64   `json:"Wk52Low"`
	PERatio        float64   `json:"PERatio"` //市盈率
	OrderDirection int       `json:"OrderDirection"`
	BidPrice       float64   `json:"BidPrice"`
	AskPrice       float64   `json:"AskPrice"`
	TurnoverRate   float64   `json:"TurnoverRate"`
	SA             float64   `json:"SA"`
	LimitUp        float64   `json:"LimitUp"`
	LimitDown      float64   `json:"LimitDown"`
	CirStock       float64   `json:"CirStock"` //流通股
	TotStock       float64   `json:"TotStock"` //总股本
	CirVal         float64   `json:"CirVal"`   //流通市值
	TotVal         float64   `json:"TotVal"`   //总市值
	NAV            float64   `json:"NAV"`      //市净利率
	Ratio          float64   `json:"Ratio"`    //量比
	Committee      float64   `json:"Committee"`
	PES            float64   `json:"PES"`
	WP             int       `json:"WP"`
	NP             int       `json:"NP"`
	LastTradeVol   int       `json:"LastTradeVol"`
	YearUpDown     float64   `json:"YearUpDown"`
	KindsUpdown    struct {
		FiveMinsUpdown  float64 `json:"FiveMinsUpdown"`
		ThreeMinsUpdown float64 `json:"ThreeMinsUpdown"`
		OneMinsUpdown   float64 `json:"OneMinsUpdown"`
		MinUpdown2      int     `json:"MinUpdown2"`
		MinUpdown4      int     `json:"MinUpdown4"`
	} `json:"KindsUpdown"`
	Updown               float64 `json:"Updown"`
	NextDayPreClosePrice float64 `json:"NextDayPreClosePrice"`
	ExchangeID           string  `json:"ExchangeID"`
	InstrumentID         string  `json:"InstrumentID"`
	TTM                  float64 `json:"TTM"`
}
type TickType struct {
	TradingDay int     `json:"TradingDay"`
	ID         int     `json:"ID"`
	Time       int     `json:"Time"`
	Price      float64 `json:"Price"`
	Volume     int     `json:"Volume"`
	Property   int     `json:"Property"`
	Virtual    int     `json:"virtual"`
}
type InstStatusType struct {
	StatusType   int    `json:"StatusType"`
	ExchangeID   string `json:"ExchangeID"`
	InstrumentID string `json:"InstrumentID"`
}
type KlineType struct {
	TradingDay int     `json:"TradingDay"`
	Time       int     `json:"Time"`
	High       float64 `json:"High"`
	Open       float64 `json:"Open"`
	Low        float64 `json:"Low"`
	Close      float64 `json:"Close"`
	Volume     int     `json:"Volume"`
	Amount     float64 `json:"Amount"` //交易金额
	TickCount  int     `json:"TickCount"`
}
type StatisticType struct {
	TradingDay int `json:"TradingDay"`
	//昨天价格
	PreClosePrice float64 `json:"PreClosePrice"`
	//开盘价
	OpenPrice       float64 `json:"OpenPrice"`
	UpperLimitPrice float64 `json:"UpperLimitPrice"`
	LowerLimitPrice float64 `json:"LowerLimitPrice"`
	ExchangeID      string  `json:"ExchangeID"`
	InstrumentID    string  `json:"InstrumentID"`
}

type MinType struct {
	TradingDay       int     `json:"TradingDay"`
	Time             int     `json:"Time"`
	Price            float64 `json:"Price"`
	Volume           int     `json:"Volume"`
	UnmismatchVolume int     `json:"UnmismatchVolume,omitempty"`
	UnmismatchFlag   int     `json:"UnmismatchFlag,omitempty"`
}
type StaticType struct {
	ExchangeID     string `json:"ExchangeID"`
	ExchangeName   string `json:"ExchangeName"`
	InstrumentID   string `json:"InstrumentID"`
	InstrumentName string `json:"InstrumentName"`
	Tradetime      struct {
		Timezone int   `json:"Timezone"`
		Pre      int   `json:"Pre"`
		Duration []int `json:"Duration"`
		Tail     int   `json:"Tail"`
	} `json:"Tradetime"`
	PriceMoneyType             int    `json:"PriceMoneyType"`
	WeightType                 int    `json:"WeightType"`
	TradingUnit                int    `json:"TradingUnit"`
	MinPriceChange             int    `json:"MinPriceChange"`
	MaxPriceLimitPerc          int    `json:"MaxPriceLimitPerc"`
	MinNumSingleForm           int    `json:"MinNumSingleForm"`
	MaxNumSingleForm           int    `json:"MaxNumSingleForm"`
	MinTradingUnit             int    `json:"MinTradingUnit"`
	MaxTradingUnit             int    `json:"MaxTradingUnit"`
	TradingMargin              int    `json:"TradingMargin"`
	DeliveryDate               int    `json:"DeliveryDate"`
	TradingFeeExtremelyRatio   int    `json:"TradingFeeExtremelyRatio"`
	DeliveryFee                int    `json:"DeliveryFee"`
	MarginPerc                 int    `json:"MarginPerc"`
	FetishPaymentPerc          int    `json:"FetishPaymentPerc"`
	DelayPaymentExtremelyRatio int    `json:"DelayPaymentExtremelyRatio"`
	DelayPaymentCalcStyle      int    `json:"DelayPaymentCalcStyle"`
	MinDeliveryApplyNum        int    `json:"MinDeliveryApplyNum"`
	PriceDecimalBitNum         int    `json:"PriceDecimalBitNum"`
	MinDeliveryUnit            int    `json:"MinDeliveryUnit"`
	CodeType                   int    `json:"CodeType"`
	CodeSecondType             int    `json:"CodeSecondType"`
	IsCrdBuyUnderlying         int    `json:"IsCrdBuyUnderlying"`
	IsCrdSellUnderlying        int    `json:"IsCrdSellUnderlying"`
	SecurityType               string `json:"SecurityType"`
	TransactionsMultiplier     int    `json:"TransactionsMultiplier"`
}
type dataRes struct {
	Market      string `json:"Market"`
	Inst        string `json:"Inst"`
	ServiceType string `json:"ServiceType"`
	SubType     string `json:"SubType"`
	ReqID       int    `json:"ReqID"`
	QuoteData   struct {
		//动态数据
		DynaData []DynaType `json:"DynaData"`
		//最后交易买单  实时
		TickData []TickType `json:"TickData"`
		//不知道干嘛的
		InstStatusData []InstStatusType `json:"InstStatusData"`
		//k线 不解析了
		//KlineData []KlineType `json:"KlineData"`
		//不动的数据 开盘价  昨天收盘价等等
		StatisticsData []StatisticType `json:"StatisticsData"`
		//早盘 数据
		MinData []MinType `json:"MinData"`
		//基本信息
		StaticData []StaticType `json:"StaticData"`
	} `json:"QuoteData"`
}
type PingType struct {
	ServiceType string `json:"ServiceType"`
}
type Pong struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}
type MsgType struct {
	Id  string
	Msg string
}

type VRaInnner struct {
	t   int
	val float64
}
type VRa struct {
	q   []VRaInnner
	n   int
	sum float64
}

// 处理消息的程序是 单协程
func (this *VRa) Push(r VRaInnner) {
	this.q = append(this.q, r)
	this.n += 1
	this.sum += r.val
	i := 0
	tn := this.n
	for ; i < tn; i++ {
		if r.t-this.q[i].t > LB {
			this.n--
			this.sum -= this.q[i].val
		} else {
			break
		}
	}
	this.q = this.q[i:tn]
}

func (this *VRa) GetAvg() float64 {
	if this.n == 0 {
		return 100000000
	}
	return this.sum / float64(this.n)
}
