package jfzt

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
	TradingDay     int       `json:"TradingDay"`     // 交易日期
	Time           int       `json:"Time"`           // 时间
	HighestPrice   float64   `json:"HighestPrice"`   // 最高价格
	LowestPrice    float64   `json:"LowestPrice"`    // 最低价格
	LastPrice      float64   `json:"LastPrice"`      // 最新价格
	Volume         int       `json:"Volume"`         // 成交量
	Amount         float64   `json:"Amount"`         // 成交额
	TickCount      int       `json:"TickCount"`      // 逐笔成交数
	BuyPrice       []float64 `json:"BuyPrice"`       // 买入价格
	BuyVolume      []int     `json:"BuyVolume"`      // 买入量
	SellPrice      []float64 `json:"SellPrice"`      // 卖出价格
	SellVolume     []int     `json:"SellVolume"`     // 卖出量
	AveragePrice   float64   `json:"AveragePrice"`   // 均价
	Wk52High       float64   `json:"Wk52High"`       // 52周最高价
	Wk52Low        float64   `json:"Wk52Low"`        // 52周最低价
	PERatio        float64   `json:"PERatio"`        // 市盈率
	OrderDirection int       `json:"OrderDirection"` // 委托方向
	BidPrice       float64   `json:"BidPrice"`       // 买入竞价
	AskPrice       float64   `json:"AskPrice"`       // 卖出竞价
	TurnoverRate   float64   `json:"TurnoverRate"`   // 换手率
	SA             float64   `json:"SA"`             // SA
	LimitUp        float64   `json:"LimitUp"`        // 涨停价
	LimitDown      float64   `json:"LimitDown"`      // 跌停价
	CirStock       float64   `json:"CirStock"`       // 流通股本
	TotStock       float64   `json:"TotStock"`       // 总股本
	CirVal         float64   `json:"CirVal"`         // 流通市值
	TotVal         float64   `json:"TotVal"`         // 总市值
	NAV            float64   `json:"NAV"`            // 市净率
	Ratio          float64   `json:"Ratio"`          // 量比
	Committee      float64   `json:"Committee"`      // 委比
	PES            float64   `json:"PES"`            // 委差
	WP             int       `json:"WP"`             // 外盘
	NP             int       `json:"NP"`             // 内盘
	LastTradeVol   int       `json:"LastTradeVol"`   // 最后一笔成交量
	YearUpDown     float64   `json:"YearUpDown"`     // 年涨跌幅
	KindsUpdown    struct {
		FiveMinsUpdown  float64 `json:"FiveMinsUpdown"`  // 五分钟涨跌幅
		ThreeMinsUpdown float64 `json:"ThreeMinsUpdown"` // 三分钟涨跌幅
		OneMinsUpdown   float64 `json:"OneMinsUpdown"`   // 一分钟涨跌幅
		MinUpdown2      float64 `json:"MinUpdown2"`      // 两分钟涨跌幅
		MinUpdown4      float64 `json:"MinUpdown4"`      // 四分钟涨跌幅
	} `json:"KindsUpdown"`
	Updown               float64 `json:"Updown"`               // 涨跌幅
	NextDayPreClosePrice float64 `json:"NextDayPreClosePrice"` // 下个交易日的前收盘价
	ExchangeID           string  `json:"ExchangeID"`           // 交易所代码
	InstrumentID         string  `json:"InstrumentID"`         // 合约代码
	TTM                  float64 `json:"TTM"`                  // 滚动市盈率
}

type TickType struct {
	TradingDay int     `json:"TradingDay"` // 交易日期
	ID         int     `json:"ID"`         // ID
	Time       int64   `json:"Time"`       // 时间
	Price      float64 `json:"Price"`      // 价格
	Volume     int     `json:"Volume"`     // 成交量
	Property   int     `json:"Property"`   // 属性
	Virtual    int     `json:"virtual"`    // 虚拟
}

type InstStatusType struct {
	StatusType   int    `json:"StatusType"`   // 状态类型
	ExchangeID   string `json:"ExchangeID"`   // 交易所代码
	InstrumentID string `json:"InstrumentID"` // 合约代码
}

type KlineType struct {
	TradingDay int     `json:"TradingDay"` // 交易日期
	Time       int     `json:"Time"`       // 时间
	High       float64 `json:"High"`       // 最高价
	Open       float64 `json:"Open"`       // 开盘价
	Low        float64 `json:"Low"`        // 最低价
	Close      float64 `json:"Close"`      // 收盘价
	Volume     int     `json:"Volume"`     // 成交量
	Amount     float64 `json:"Amount"`     // 交易金额
	TickCount  int     `json:"TickCount"`  // 逐笔成交数
}

type StatisticType struct {
	TradingDay      int     `json:"TradingDay"`      // 交易日期
	PreClosePrice   float64 `json:"PreClosePrice"`   // 昨日收盘价
	OpenPrice       float64 `json:"OpenPrice"`       // 开盘价
	UpperLimitPrice float64 `json:"UpperLimitPrice"` // 涨停价
	LowerLimitPrice float64 `json:"LowerLimitPrice"` // 跌停价
	ExchangeID      string  `json:"ExchangeID"`      // 交易所代码
	InstrumentID    string  `json:"InstrumentID"`    // 合约代码
}

type MinType struct {
	TradingDay       int     `json:"TradingDay"`                 // 交易日期
	Time             int     `json:"Time"`                       // 时间
	Price            float64 `json:"Price"`                      // 价格
	Volume           int     `json:"Volume"`                     // 成交量
	UnmismatchVolume int     `json:"UnmismatchVolume,omitempty"` // 不匹配的成交量
	UnmismatchFlag   int     `json:"UnmismatchFlag,omitempty"`   // 不匹配的标志
}

type StaticType struct {
	ExchangeID     string `json:"ExchangeID"`     // 交易所代码
	ExchangeName   string `json:"ExchangeName"`   // 交易所名称
	InstrumentID   string `json:"InstrumentID"`   // 合约代码
	InstrumentName string `json:"InstrumentName"` // 合约名称
	Tradetime      struct {
		Timezone int   `json:"Timezone"` // 时区
		Pre      int   `json:"Pre"`      // 交易前时间
		Duration []int `json:"Duration"` // 交易时段
		Tail     int   `json:"Tail"`     // 交易后时间
	} `json:"Tradetime"`
	PriceMoneyType             int     `json:"PriceMoneyType"`             // 价格货币类型
	WeightType                 int     `json:"WeightType"`                 // 权重类型
	TradingUnit                int     `json:"TradingUnit"`                // 交易单位
	MinPriceChange             int     `json:"MinPriceChange"`             // 最小价格变动
	MaxPriceLimitPerc          int     `json:"MaxPriceLimitPerc"`          // 最大价格限制百分比
	MinNumSingleForm           int     `json:"MinNumSingleForm"`           // 单个报单最小数量
	MaxNumSingleForm           int     `json:"MaxNumSingleForm"`           // 单个报单最大数量
	MinTradingUnit             int     `json:"MinTradingUnit"`             // 最小交易单位
	MaxTradingUnit             int     `json:"MaxTradingUnit"`             // 最大交易单位
	TradingMargin              int     `json:"TradingMargin"`              // 交易保证金
	DeliveryDate               int     `json:"DeliveryDate"`               // 交割日期
	TradingFeeExtremelyRatio   int     `json:"TradingFeeExtremelyRatio"`   // 交易手续费极限比例
	DeliveryFee                int     `json:"DeliveryFee"`                // 交割费用
	MarginPerc                 int     `json:"MarginPerc"`                 // 保证金百分比
	FetishPaymentPerc          int     `json:"FetishPaymentPerc"`          // 迷信付款百分比
	DelayPaymentExtremelyRatio int     `json:"DelayPaymentExtremelyRatio"` // 延期付款极限比例
	DelayPaymentCalcStyle      int     `json:"DelayPaymentCalcStyle"`      // 延期付款计算方式
	MinDeliveryApplyNum        int     `json:"MinDeliveryApplyNum"`        // 最小交割申报数量
	PriceDecimalBitNum         int     `json:"PriceDecimalBitNum"`         // 价格小数位数
	MinDeliveryUnit            int     `json:"MinDeliveryUnit"`            // 最小交割单位
	CodeType                   int     `json:"CodeType"`                   // 代码类型
	OpenAuctionDuration        int     `json:"OpenAuctionDuration"`        // 开盘集合竞价时长
	DayTradingFlag             int     `json:"DayTradingFlag"`             // 当日交易标志
	StrikePriceNum             int     `json:"StrikePriceNum"`             // 行权价格数量
	PriceRatio                 int     `json:"PriceRatio"`                 // 价格比率
	InstrumentState            int     `json:"InstrumentState"`            // 合约状态
	MaxOrderNum                int     `json:"MaxOrderNum"`                // 最大报单数量
	MinSplitUnit               int     `json:"MinSplitUnit"`               // 最小分割单位
	CouponRate                 int     `json:"CouponRate"`                 // 优惠比例
	MaxTradingVolume           int     `json:"MaxTradingVolume"`           // 最大交易数量
	DeliveryWay                int     `json:"DeliveryWay"`                // 交割方式
	TakerMaxOrderNum           int     `json:"TakerMaxOrderNum"`           // 做市商最大报单数量
	OptionContractType         int     `json:"OptionContractType"`         // 期权合约类型
	Category                   int     `json:"Category"`                   // 合约类别
	UnlistedFlag               int     `json:"UnlistedFlag"`               // 非上市标志
	SeparatePrice              float64 `json:"SeparatePrice"`              // 分离价格
	InsCapitalType             int     `json:"InsCapitalType"`             // 合约资金类型
	CreateTime                 int64   `json:"CreateTime"`                 // 创建时间
	UpdateTime                 int64   `json:"UpdateTime"`                 // 更新时间
	DisplayCode                string  `json:"DisplayCode"`                // 显示代码
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
		//分线 不解析了
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

type VRaInnner struct {
	t   int64
	val float64
}
type VRa struct {
	q   []VRaInnner
	sum float64
}

// 处理消息的程序是 单协程
func (this *VRa) Push(r VRaInnner) {
	this.q = append(this.q, r)
	this.sum += r.val
	i := 0
	n := len(this.q)
	for ; i < n; i++ {
		if r.t-this.q[i].t > LB {
			this.sum -= this.q[i].val
		} else {
			break
		}
	}
	this.q = this.q[i:n]
}

func (this *VRa) GetAvg() float64 {
	if len(this.q) == 0 {
		return 100000000
	}
	return this.sum / float64(len(this.q))
}
