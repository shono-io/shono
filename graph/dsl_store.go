package graph

func ListFromStore(storeKey Key, filters map[string]Expression) StoreLogic {
	return StoreLogic{
		StoreKey:  storeKey,
		Operation: StoreOperationList,
		Filters:   filters,
	}
}

func GetFromStore(storeKey Key, key *Expression) StoreLogic {
	return StoreLogic{
		StoreKey:  storeKey,
		Operation: StoreOperationGet,
		Key:       key,
	}
}

func AddToStore(storeKey Key, key Expression, value ...Mapping) StoreLogic {
	return StoreLogic{
		StoreKey:  storeKey,
		Operation: StoreOperationAdd,
		Key:       &key,
		Value:     value,
	}
}

func SetInStore(storeKey Key, key Expression, value ...Mapping) StoreLogic {
	return StoreLogic{
		StoreKey:  storeKey,
		Operation: StoreOperationSet,
		Key:       &key,
		Value:     value,
	}
}

func RemoveFromStore(storeKey Key, key Expression) StoreLogic {
	return StoreLogic{
		StoreKey:  storeKey,
		Operation: StoreOperationDelete,
		Key:       &key,
	}
}

type StoreOperation string

var (
	StoreOperationList StoreOperation = "list"
	StoreOperationGet  StoreOperation = "get"

	StoreOperationAdd    StoreOperation = "add"
	StoreOperationSet    StoreOperation = "set"
	StoreOperationDelete StoreOperation = "delete"
)

type StoreLogic struct {
	StoreKey  Key
	Operation StoreOperation
	Key       *Expression
	Value     []Mapping
	Filters   map[string]Expression
}

func (s StoreLogic) Kind() string {
	return "store"
}
