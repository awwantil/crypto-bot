package options

type ApiOptions struct {
	Key           string
	Secret        string
	Passphrase    string
	ClientId      string
	IsDemoTrading string
}

type ApiOption func(options *ApiOptions)

func WithApiKey(key string) ApiOption {
	return func(options *ApiOptions) {
		options.Key = key
	}
}

func WithApiSecretKey(secret string) ApiOption {
	return func(options *ApiOptions) {
		options.Secret = secret
	}
}

func WithPassphrase(passphrase string) ApiOption {
	return func(options *ApiOptions) {
		options.Passphrase = passphrase
	}
}

func WithIsDemoTrade(isDemo string) ApiOption {
	return func(options *ApiOptions) {
		options.IsDemoTrading = isDemo
	}
}

func WithClientId(clientId string) ApiOption {
	return func(options *ApiOptions) {
		options.ClientId = clientId
	}
}
