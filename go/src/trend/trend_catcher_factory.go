package trend

type TrendCatcherFactory struct{}

func (f *TrendCatcherFactory) CreateTrendCatcher() TrendCatcher {
	return &GoogleTrendCatcher{}
}
