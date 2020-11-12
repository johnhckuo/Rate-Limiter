package environment

const (
	//PersistStorage is used to decide which storage type to use
	PersistStorage = "PERSIST_STORAGE"
	//RateLimiter is used to decide which rate limiter type to use
	RateLimiter = "RATE_LIMITER"
	//GoLimiter is the default value for choosing rate limiter provided by Golang community
	GoLimiter = "GO_LIMITER"
	//RedisLimiter is the default value for choosing my homemade redis rate limiter :)
	RedisLimiter = "REDIS_LIMITER"
	//RedisStorage is the default value for choosing redis as persistent storage
	RedisStorage = "REDIS_STORAGE"
	//DbConnectionString is the DB connection string
	DbConnectionString = "DB_CONNECTION_STRING"

	//RateLimit restricts the number of request each "RateLimitExpirationSecond" can pass
	RateLimit = "RATE_LIMIT"
	//RateLimitExpirationSecond indicates the time interval of rate limiter
	RateLimitExpirationSecond = "RATE_LIMIT_EXPIRATION_SECOND"
	//BurstLimit restricts the number of burst request
	BurstLimit = "BURST_LIMIT"
)
