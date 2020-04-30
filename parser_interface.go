package i18n

type IParser interface {
    SetOptions(opts *Options)
    Parse() error
    Load(key string, defaultVal ...string) interface{}
}
