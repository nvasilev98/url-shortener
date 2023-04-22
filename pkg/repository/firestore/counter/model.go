package counter

type Shard struct {
	Count int64 `firestore:"count"`
}
