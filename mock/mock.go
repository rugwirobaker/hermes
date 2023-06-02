package mock

//go:generate mockgen -package=mock -destination=mock_gen.go github.com/rugwirobaker/hermes SendService,Pubsub,AppStore,MessageStore
