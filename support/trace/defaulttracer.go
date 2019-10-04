package trace

type nooptracer struct{}

func (nt *nooptracer) Name() string {
	return "noop-tracer"
}

func (nt *nooptracer) Inject(tCtx TracingContext, format CarrierFormat, carrier interface{}) error {
	return nil
}
func (nt *nooptracer) Extract(format CarrierFormat, data interface{}) (TracingContext, error) {
	return TracingContext{}, nil
}
func (nt *nooptracer) SetTag(tCtx TracingContext, TagKey string, TagValue interface{}) bool {
	return false
}
func (nt *nooptracer) LogKV(tCtx TracingContext, alternatingKeyValues ...interface{}) bool {
	return false
}
func (nt *nooptracer) StartSpan(config Config, parent TracingContext) (TracingContext, error) {
	return TracingContext{}, nil
}
func (nt *nooptracer) FinishSpan(tContext TracingContext, err error) error {
	return nil
}
func (nt *nooptracer) Start() error {
	//fmt.Println(engine.GetAppName(),"-", engine.GetAppVersion())
	return nil
}
func (nt *nooptracer) Stop() error {
	return nil
}