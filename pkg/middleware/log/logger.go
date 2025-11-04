package log

//
//// NewLogger 创建一个满足 kratos/log.Logger 接口的 slog logger
//func NewLogger(serviceName string) log.Logger {
//	// 使用 slog.NewJSONHandler 创建 JSON 格式的 handler
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
//		AddSource: true, // 记录文件和行号
//		Level:     slog.LevelDebug,
//	}))
//
//	slogLogger := logger.With("service", serviceName)
//
//	return log.NewHelper(log.NewSlogLogger(slogLogger))
//}
//
//// NewSlogLogger 是一个适配器
//type slogLogger struct {
//	l *slog.Logger
//}
//
//func (l *slogLogger) Log(level log.Level, keyvals ...interface{}) error {
//	// ... slog 和 kratos log 的适配逻辑 (具体实现略)
//	// 官方已有或即将有实现，这里仅为示意
//	return nil
//}
