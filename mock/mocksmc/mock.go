package mocksmc

//go:generate mockgen -package=mocksmc -destination=mock_gen.go github.com/quarksgroup/sms-client/sms SendService,AuthService
