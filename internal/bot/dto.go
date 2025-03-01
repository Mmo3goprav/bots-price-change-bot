package bot

const (
	initMsg             = "Hello, I am a bot that shows price changes in financial markets.\nall you need to do is select <b>add chart</b> and I will notify you of price changes"
	addChartMsg         = "specify the symbol you want to track in the format: <b>BTCUSD</b>"
	chartAdded          = "chart <b>%s</b> successfully added"
	chartAlreadyAdded   = "chart <b>%s</b> already added"
	removeChartMsg      = "specify the symbol you don't want to track in the format: <b>BTCUSD</b>"
	chartRemoved        = "chart <b>%s</b> successfully removed"
	chartAlreadyRemoved = "chart <b>%s</b> already removed"
	helpMsg             = "unknown command, please type <b>/start</b> to start use bot"
	chartMsg            = "Symbol: <b>%s</b> \ncurrent price: %f \nprice change: %f \n price change percentage %f\n\n"
)

const (
	commandStart = "start"
	//
	callbackAddChart    = "Add Chart"
	callbackRemoveChart = "Remove Chart"
)

const (
	parseModHTML = "HTML"
)
