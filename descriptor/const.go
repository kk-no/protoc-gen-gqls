package descriptor

const (
	Indent = "    "
	Type   = "type "
	Input  = "input "
	Enum   = "enum "
	Open   = " {"
	Close  = "}"
)

const (
	QueryType    = Type + "Query" + Open
	MutationType = Type + "Mutation" + Open
)
