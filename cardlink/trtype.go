package cardlink

// TrType values for redirect field trType (payment vs pre-authorization).
// Tokenization uses trType 8 per docs — not implemented in initial SDK.
const (
	TrTypeSale    = "1" // payment (default)
	TrTypePreauth = "2" // pre-authorization (card payments)
	// TrTypeTokenizer = "8" // tokenizer — deferred
)
