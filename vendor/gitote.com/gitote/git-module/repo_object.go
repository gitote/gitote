package git

type ObjectType string

const (
	OBJECT_COMMIT ObjectType = "commit"
	OBJECT_TREE   ObjectType = "tree"
	OBJECT_BLOB   ObjectType = "blob"
	OBJECT_TAG    ObjectType = "tag"
)
